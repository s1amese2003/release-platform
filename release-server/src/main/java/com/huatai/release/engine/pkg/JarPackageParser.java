package com.huatai.release.engine.pkg;

import com.huatai.release.model.dto.DependencyCoordinate;
import com.huatai.release.model.dto.PackageParseResult;
import org.apache.commons.codec.digest.DigestUtils;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.web.multipart.MultipartFile;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;
import java.util.jar.Attributes;
import java.util.jar.Manifest;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;

@Component
public class JarPackageParser {

    private static final Pattern UPGRADE_VERSION_PATTERN = Pattern.compile("^BOOT-INF/classes/upgrade/(\\d{8})/.*");
    private static final Pattern CONFIG_PATTERN = Pattern.compile("^BOOT-INF/classes/(bootstrap.*\\.yml|application.*\\.yml)$");
    private static final Pattern DEP_PATTERN = Pattern.compile("^BOOT-INF/lib/(.+)\\.jar$");

    @Value("${release.work-dir:./workdir}")
    private String workDir;

    public PackageParseResult parse(MultipartFile file, String appName, String targetEnv, String submitter) {
        try {
            String requestNo = buildRequestNo();
            Path base = Paths.get(workDir).resolve(requestNo);
            Files.createDirectories(base);

            Path packagePath = base.resolve(Optional.ofNullable(file.getOriginalFilename()).orElse("release-package.jar"));
            file.transferTo(packagePath);
            String md5 = md5(packagePath);

            List<String> versions = new ArrayList<>();
            Map<String, String> sqlByVersion = new HashMap<>();
            Map<String, String> configTxtByVersion = new HashMap<>();
            Map<String, String> configFiles = new LinkedHashMap<>();
            List<DependencyCoordinate> dependencies = new ArrayList<>();
            String appVersion = "unknown";
            String buildJdk = "unknown";
            LocalDateTime buildTime = null;

            try (ZipFile zipFile = openZipWithFallback(packagePath)) {
                Manifest manifest = loadManifest(zipFile);
                if (manifest != null) {
                    Attributes attrs = manifest.getMainAttributes();
                    appVersion = firstNonBlank(
                            attrs.getValue("Implementation-Version"),
                            attrs.getValue("Build-Version"),
                            attrs.getValue("App-Version"),
                            "unknown");
                    buildJdk = firstNonBlank(attrs.getValue("Build-Jdk"), attrs.getValue("Build-JDK"), "unknown");
                    buildTime = parseBuildTime(attrs.getValue("Build-Time"));
                }

                zipFile.stream().forEach(entry -> {
                    String normalizedName = entry.getName().replace('\\', '/');
                    Matcher versionMatcher = UPGRADE_VERSION_PATTERN.matcher(normalizedName);
                    if (versionMatcher.matches()) {
                        String version = versionMatcher.group(1);
                        if (!versions.contains(version)) {
                            versions.add(version);
                        }
                    }

                    if (normalizedName.endsWith("/upgrade.sql") && normalizedName.contains("BOOT-INF/classes/upgrade/")) {
                        String version = extractVersion(normalizedName);
                        if (version != null) {
                            sqlByVersion.put(version, readEntry(zipFile, entry));
                        }
                    }

                    if (normalizedName.endsWith("/config.txt") && normalizedName.contains("BOOT-INF/classes/upgrade/")) {
                        String version = extractVersion(normalizedName);
                        if (version != null) {
                            configTxtByVersion.put(version, readEntry(zipFile, entry));
                        }
                    }

                    if (CONFIG_PATTERN.matcher(normalizedName).matches()) {
                        configFiles.put(normalizedName, readEntry(zipFile, entry));
                    }

                    Matcher depMatcher = DEP_PATTERN.matcher(normalizedName);
                    if (depMatcher.matches()) {
                        dependencies.add(parseDependency(depMatcher.group(1)));
                    }
                });
            }

            versions.sort(Comparator.reverseOrder());
            String latest = versions.isEmpty() ? null : versions.get(0);
            String sqlContent = latest == null ? "" : sqlByVersion.getOrDefault(latest, "");
            String configTxt = latest == null ? "" : configTxtByVersion.getOrDefault(latest, "");
            String sqlFilePath = latest == null ? "" : "BOOT-INF/classes/upgrade/" + latest + "/upgrade.sql";

            return new PackageParseResult(
                    requestNo,
                    appName,
                    appVersion,
                    targetEnv,
                    buildTime,
                    buildJdk,
                    packagePath.toString(),
                    md5,
                    versions,
                    latest,
                    sqlFilePath,
                    sqlContent,
                    configTxt,
                    parseConfigTxt(configTxt),
                    configFiles,
                    dependencies
            );
        } catch (IOException ex) {
            throw new IllegalStateException("Package parse failed: " + ex.getMessage(), ex);
        }
    }

    private ZipFile openZipWithFallback(Path packagePath) throws IOException {
        List<Charset> candidates = List.of(
                StandardCharsets.UTF_8,
                Charset.forName("GBK"),
                Charset.defaultCharset()
        );
        LinkedHashSet<Charset> unique = new LinkedHashSet<>(candidates);

        IOException last = null;
        for (Charset charset : unique) {
            try {
                return new ZipFile(packagePath.toFile(), charset);
            } catch (IOException ex) {
                last = ex;
            }
        }
        if (last != null) {
            throw last;
        }
        throw new IOException("Unable to open archive with supported charsets");
    }

    private String buildRequestNo() {
        String date = DateTimeFormatter.ofPattern("yyyyMMdd", Locale.CHINA).format(LocalDateTime.now());
        String suffix = UUID.randomUUID().toString().replace("-", "").substring(0, 4).toUpperCase(Locale.ROOT);
        return "REL-" + date + "-" + suffix;
    }

    private String md5(Path file) throws IOException {
        try (InputStream is = new FileInputStream(file.toFile())) {
            return DigestUtils.md5Hex(is);
        }
    }

    private Manifest loadManifest(ZipFile zipFile) {
        ZipEntry entry = zipFile.getEntry("META-INF/MANIFEST.MF");
        if (entry == null) {
            return null;
        }
        try (InputStream is = zipFile.getInputStream(entry)) {
            return new Manifest(is);
        } catch (IOException ex) {
            return null;
        }
    }

    private String readEntry(ZipFile zipFile, ZipEntry entry) {
        try (InputStream is = zipFile.getInputStream(entry)) {
            String content = new String(is.readAllBytes(), StandardCharsets.UTF_8);
            if (content.startsWith("\uFEFF")) {
                return content.substring(1);
            }
            return content;
        } catch (IOException ex) {
            return "";
        }
    }

    private String extractVersion(String path) {
        Matcher matcher = UPGRADE_VERSION_PATTERN.matcher(path);
        if (!matcher.matches()) {
            return null;
        }
        return matcher.group(1);
    }

    private DependencyCoordinate parseDependency(String depName) {
        int idx = depName.lastIndexOf('-');
        if (idx <= 0 || idx == depName.length() - 1) {
            return new DependencyCoordinate("unknown", depName, "unknown", depName + ".jar");
        }
        String artifact = depName.substring(0, idx);
        String version = depName.substring(idx + 1);
        return new DependencyCoordinate("unknown", artifact, version, depName + ".jar");
    }

    private LocalDateTime parseBuildTime(String raw) {
        if (raw == null || raw.isBlank()) {
            return null;
        }
        List<DateTimeFormatter> formatters = List.of(
                DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"),
                DateTimeFormatter.ISO_DATE_TIME
        );
        for (DateTimeFormatter formatter : formatters) {
            try {
                return LocalDateTime.parse(raw, formatter);
            } catch (Exception ignored) {
            }
        }
        return null;
    }

    private String firstNonBlank(String... values) {
        for (String value : values) {
            if (value != null && !value.isBlank()) {
                return value;
            }
        }
        return "";
    }

    private List<String> parseConfigTxt(String configTxt) {
        if (configTxt == null || configTxt.isBlank()) {
            return List.of();
        }
        return configTxt.lines()
                .map(String::trim)
                .filter(line -> !line.isBlank())
                .filter(line -> !line.startsWith("#"))
                .collect(Collectors.toList());
    }
}
