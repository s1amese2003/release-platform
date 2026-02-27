package service

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ArtifactService struct {
	storageRoot string
}

func NewArtifactService(storageRoot string) *ArtifactService {
	return &ArtifactService{storageRoot: storageRoot}
}

func (s *ArtifactService) SaveUploadedArtifact(releaseID uint, fileHeader *multipart.FileHeader) (storagePath string, size int64, sha string, err error) {
	if err := os.MkdirAll(filepath.Join(s.storageRoot, "artifacts"), 0o755); err != nil {
		return "", 0, "", fmt.Errorf("create artifact directory: %w", err)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", 0, "", fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	safeName := sanitizeFileName(fileHeader.Filename)
	storedName := fmt.Sprintf("%d_%s_%s", releaseID, uuid.NewString(), safeName)
	storagePath = filepath.Join(s.storageRoot, "artifacts", storedName)

	dst, err := os.Create(storagePath)
	if err != nil {
		return "", 0, "", fmt.Errorf("create artifact file: %w", err)
	}
	defer dst.Close()

	hasher := sha256.New()
	writer := io.MultiWriter(dst, hasher)
	written, err := io.Copy(writer, src)
	if err != nil {
		return "", 0, "", fmt.Errorf("persist artifact: %w", err)
	}

	size = written
	sha = hex.EncodeToString(hasher.Sum(nil))
	return storagePath, size, sha, nil
}

func (s *ArtifactService) ExtractArtifact(releaseID uint, artifactPath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(artifactPath))
	if ext != ".zip" && ext != ".jar" && ext != ".war" {
		return "", fmt.Errorf("unsupported artifact extension %q", ext)
	}

	target := filepath.Join(s.storageRoot, "extracted", fmt.Sprintf("%d_%d", releaseID, time.Now().UnixNano()))
	if err := os.MkdirAll(target, 0o755); err != nil {
		return "", fmt.Errorf("create extraction directory: %w", err)
	}

	if err := extractZipSecure(artifactPath, target); err != nil {
		return "", err
	}

	return target, nil
}

func (s *ArtifactService) LocateLatestUpgradeSQL(extractedDir string) (string, error) {
	pattern := regexp.MustCompile(`(?i)^upgrade/(\d{8})/upgrade\.sql$`)
	roots, err := findClassesRoots(extractedDir)
	if err != nil {
		return "", err
	}

	var candidates []struct {
		date int
		path string
	}
	for _, root := range roots {
		err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			rel = filepath.ToSlash(rel)
			matched := pattern.FindStringSubmatch(rel)
			if len(matched) != 2 {
				return nil
			}
			date, convErr := strconv.Atoi(matched[1])
			if convErr != nil {
				return nil
			}
			candidates = append(candidates, struct {
				date int
				path string
			}{date: date, path: path})
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("walk upgrade sql path: %w", err)
		}
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("latest upgrade.sql not found")
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].date > candidates[j].date })
	return candidates[0].path, nil
}

func (s *ArtifactService) LocateBootstrapDevYAML(extractedDir string) (string, error) {
	roots, err := findClassesRoots(extractedDir)
	if err != nil {
		return "", err
	}

	for _, root := range roots {
		yml := filepath.Join(root, "bootstrap-dev.yml")
		if _, err := os.Stat(yml); err == nil {
			return yml, nil
		}
		yaml := filepath.Join(root, "bootstrap-dev.yaml")
		if _, err := os.Stat(yaml); err == nil {
			return yaml, nil
		}
	}
	return "", fmt.Errorf("bootstrap-dev.yml/bootstrap-dev.yaml not found")
}

func extractZipSecure(zipPath, destination string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip archive: %w", err)
	}
	defer reader.Close()

	destClean := filepath.Clean(destination)
	prefix := destClean + string(os.PathSeparator)

	for _, file := range reader.File {
		entryName := strings.ReplaceAll(file.Name, "\\", "/")
		entryName = strings.TrimLeft(entryName, "/")
		if entryName == "" {
			continue
		}

		targetPath := filepath.Join(destination, filepath.FromSlash(entryName))
		cleanTarget := filepath.Clean(targetPath)
		if cleanTarget != destClean && !strings.HasPrefix(cleanTarget, prefix) {
			return fmt.Errorf("zip slip detected for %s", entryName)
		}

		isDirEntry := file.FileInfo().IsDir() || strings.HasSuffix(entryName, "/")
		if isDirEntry {
			if err := os.MkdirAll(cleanTarget, 0o755); err != nil {
				return fmt.Errorf("create directory %s: %w", cleanTarget, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(cleanTarget), 0o755); err != nil {
			return fmt.Errorf("create parent directory: %w", err)
		}

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("open file from archive: %w", err)
		}

		dst, err := os.OpenFile(cleanTarget, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			src.Close()
			return fmt.Errorf("create extracted file: %w", err)
		}

		_, copyErr := io.Copy(dst, src)
		closeErr := dst.Close()
		srcErr := src.Close()
		if copyErr != nil {
			return fmt.Errorf("extract file: %w", copyErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close extracted file: %w", closeErr)
		}
		if srcErr != nil {
			return fmt.Errorf("close archive file: %w", srcErr)
		}
	}
	return nil
}

func sanitizeFileName(fileName string) string {
	replacer := strings.NewReplacer("..", "", "\\", "_", "/", "_", " ", "_")
	return replacer.Replace(fileName)
}

func findClassesRoots(extractedDir string) ([]string, error) {
	roots := make([]string, 0, 2)
	err := filepath.WalkDir(extractedDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if isBootInfClassesDir(path) {
			roots = append(roots, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk extracted path: %w", err)
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("BOOT-INF/classes not found")
	}
	sort.Slice(roots, func(i, j int) bool {
		if len(roots[i]) == len(roots[j]) {
			return roots[i] < roots[j]
		}
		return len(roots[i]) < len(roots[j])
	})
	return roots, nil
}

func isBootInfClassesDir(path string) bool {
	clean := filepath.ToSlash(filepath.Clean(path))
	parts := strings.Split(clean, "/")
	if len(parts) < 2 {
		return false
	}
	return strings.EqualFold(parts[len(parts)-2], "BOOT-INF") && strings.EqualFold(parts[len(parts)-1], "classes")
}
