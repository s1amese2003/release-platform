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
import java.nio.file.StandardCopyOption;
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
    private static final int MAX_NESTED_DEPTH = 2;

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

            ParsedArchiveData parsed = parseArchive(packagePath, 0, base);
            parsed.versions.sort(Comparator.reverseOrder());

            String latest = parsed.versions.isEmpty() ? null : parsed.versions.get(0);
            String sqlContent = latest == null ? "" : parsed.sqlByVersion.getOrDefault(latest, "");
            String configTxt = latest == null ? "" : parsed.configTxtByVersion.getOrDefault(latest, "");
            String sqlFilePath = latest == null ? "" : "BOOT-INF/classes/upgrade/" + latest + "/upgrade.sql";

            return new PackageParseResult(
                    requestNo,
                    appName,
                    parsed.appVersion,
                    targetEnv,
                    parsed.buildTime,
                    parsed.buildJdk,
                    packagePath.toString(),
                    md5,
                    parsed.versions,
                    latest,
                    sqlFilePath,
                    sqlContent,
                    configTxt,
                    parseConfigTxt(configTxt),
                    parsed.configFiles,
                    parsed.dependencies
            );
        } catch (IOException ex) {
            throw new IllegalStateException("Package parse failed: " + ex.getMessage(), ex);
        }
    }

    private ParsedArchiveData parseArchive(Path archivePath, int depth, Path baseDir) throws IOException {
        ParsedArchiveData current = parseArchiveDirect(archivePath);

        if (current.hasReleaseMarkers() || depth >= MAX_NESTED_DEPTH) {
            return current;
        }

        ParsedArchiveData best = current;
        int index = 0;

        try (ZipFile zipFile = openZipWithFallback(archivePath)) {
            var entries = zipFile.entries();
            while (entries.hasMoreElements()) {
                ZipEntry entry = entries.nextElement();
                if (entry.isDirectory()) {
                    continue;
                }

                String name = entry.getName().replace('\\', '/');
                String lower = name.toLowerCase(Locale.ROOT);
                if (!lower.endsWith(".jar") && !lower.endsWith(".zip")) {
                    continue;
                }

                Path nestedPath = baseDir.resolve("nested-" + depth + "-" + (index++) + "-" + sanitizeFileName(Paths.get(name).getFileName().toString()));
                try (InputStream is = zipFile.getInputStream(entry)) {
                    Files.copy(is, nestedPath, StandardCopyOption.REPLACE_EXISTING);
                }

                ParsedArchiveData nested = parseArchive(nestedPath, depth + 1, baseDir);
                if (nested.score() > best.score()) {
                    best = nested;
                }
            }
        }

        return best;
    }

    private ParsedArchiveData parseArchiveDirect(Path archivePath) throws IOException {
        ParsedArchiveData data = new ParsedArchiveData();

        try (ZipFile zipFile = openZipWithFallback(archivePath)) {
            Manifest manifest = loadManifest(zipFile);
            if (manifest != null) {
                Attributes attrs = manifest.getMainAttributes();
                data.appVersion = firstNonBlank(
                        attrs.getValue("Implementation-Version"),
                        attrs.getValue("Build-Version"),
                        attrs.getValue("App-Version"),
                        "unknown");
                data.buildJdk = firstNonBlank(attrs.getValue("Build-Jdk"), attrs.getValue("Build-JDK"), "unknown");
                data.buildTime = parseBuildTime(attrs.getValue("Build-Time"));
            }

            zipFile.stream().forEach(entry -> {
                String normalizedName = entry.getName().replace('\\', '/');
                Matcher versionMatcher = UPGRADE_VERSION_PATTERN.matcher(normalizedName);
                if (versionMatcher.matches()) {
                    String version = versionMatcher.group(1);
                    if (!data.versions.contains(version)) {
                        data.versions.add(version);
                    }
                }

                if (normalizedName.endsWith("/upgrade.sql") && normalizedName.contains("BOOT-INF/classes/upgrade/")) {
                    String version = extractVersion(normalizedName);
                    if (version != null) {
                        data.sqlByVersion.put(version, readEntry(zipFile, entry));
                    }
                }

                if (normalizedName.endsWith("/config.txt") && normalizedName.contains("BOOT-INF/classes/upgrade/")) {
                    String version = extractVersion(normalizedName);
                    if (version != null) {
                        data.configTxtByVersion.put(version, readEntry(zipFile, entry));
                    }
                }

                if (CONFIG_PATTERN.matcher(normalizedName).matches()) {
                    data.configFiles.put(normalizedName, readEntry(zipFile, entry));
                }

                Matcher depMatcher = DEP_PATTERN.matcher(normalizedName);
                if (depMatcher.matches()) {
                    data.dependencies.add(parseDependency(depMatcher.group(1)));
                }
            });
        }

        return data;
    }

    private ZipFile openZipWithFallback(Path packagePath) throws IOException {
        List<Charset> candidates = List.of(
                StandardCharsets.UTF_8,
                Charset.forName("GBK"),
                Charset.defaultCharset()
        );

        IOException last = null;
        for (Charset charset : new LinkedHashSet<>(candidates)) {
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

    private String sanitizeFileName(String name) {
        return name.replaceAll("[^a-zA-Z0-9._-]", "_");
    }

    private static class ParsedArchiveData {
        private String appVersion = "unknown";
        private String buildJdk = "unknown";
        private LocalDateTime buildTime;

        private final List<String> versions = new ArrayList<>();
        private final Map<String, String> sqlByVersion = new HashMap<>();
        private final Map<String, String> configTxtByVersion = new HashMap<>();
        private final Map<String, String> configFiles = new LinkedHashMap<>();
        private final List<DependencyCoordinate> dependencies = new ArrayList<>();

        private boolean hasReleaseMarkers() {
            return !versions.isEmpty()
                    || !sqlByVersion.isEmpty()
                    || !configTxtByVersion.isEmpty()
                    || !configFiles.isEmpty()
                    || !dependencies.isEmpty();
        }

        private int score() {
            int score = 0;
            score += versions.size() * 100;
            score += sqlByVersion.size() * 80;
            score += configTxtByVersion.size() * 50;
            score += configFiles.size() * 30;
            score += Math.min(dependencies.size(), 50);
            if (!"unknown".equalsIgnoreCase(appVersion)) {
                score += 10;
            }
            return score;
        }
    }
}
