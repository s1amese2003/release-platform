package service

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractZipSecure_WindowsStylePath(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "windows-style.zip")
	outDir := filepath.Join(tmp, "out")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip file: %v", err)
	}
	zw := zip.NewWriter(f)
	entry, err := zw.Create("1.9.0\\BOOT-INF\\classes\\upgrade\\20251128\\upgrade.sql")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := entry.Write([]byte("SELECT 1;")); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	if err := extractZipSecure(zipPath, outDir); err != nil {
		t.Fatalf("extract zip: %v", err)
	}

	expected := filepath.Join(outDir, "1.9.0", "BOOT-INF", "classes", "upgrade", "20251128", "upgrade.sql")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("expected extracted file not found: %v", err)
	}
}

func TestExtractZipSecure_WindowsStyleDirectoryEntries(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "windows-dir-entries.zip")
	outDir := filepath.Join(tmp, "out")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip file: %v", err)
	}
	zw := zip.NewWriter(f)
	if _, err := zw.Create("1.9.0\\BOOT-INF\\"); err != nil {
		t.Fatalf("create BOOT-INF dir entry: %v", err)
	}
	if _, err := zw.Create("1.9.0\\BOOT-INF\\classes\\upgrade\\"); err != nil {
		t.Fatalf("create upgrade dir entry: %v", err)
	}
	entry, err := zw.Create("1.9.0\\BOOT-INF\\classes\\upgrade\\20251128\\upgrade.sql")
	if err != nil {
		t.Fatalf("create sql entry: %v", err)
	}
	if _, err := entry.Write([]byte("SELECT 1;")); err != nil {
		t.Fatalf("write sql entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	if err := extractZipSecure(zipPath, outDir); err != nil {
		t.Fatalf("extract zip with dir entries: %v", err)
	}

	expected := filepath.Join(outDir, "1.9.0", "BOOT-INF", "classes", "upgrade", "20251128", "upgrade.sql")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("expected extracted file not found: %v", err)
	}
}
