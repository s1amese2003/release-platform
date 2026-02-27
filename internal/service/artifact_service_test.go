package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocateLatestUpgradeSQL_WithOuterVersionFolder(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "1.9.0", "BOOT-INF", "classes", "upgrade")
	if err := os.MkdirAll(filepath.Join(base, "20240101"), 0o755); err != nil {
		t.Fatalf("mkdir old sql dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(base, "20251128"), 0o755); err != nil {
		t.Fatalf("mkdir new sql dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(base, "20240101", "upgrade.sql"), []byte("CREATE TABLE t1(id bigint);"), 0o644); err != nil {
		t.Fatalf("write old sql: %v", err)
	}
	expected := filepath.Join(base, "20251128", "upgrade.sql")
	if err := os.WriteFile(expected, []byte("CREATE TABLE t2(id bigint);"), 0o644); err != nil {
		t.Fatalf("write new sql: %v", err)
	}

	svc := NewArtifactService(tmp)
	found, err := svc.LocateLatestUpgradeSQL(tmp)
	if err != nil {
		t.Fatalf("locate latest sql: %v", err)
	}
	if found != expected {
		t.Fatalf("expected %s, got %s", expected, found)
	}
}

func TestLocateBootstrapDevYAML_SupportYamlExtension(t *testing.T) {
	tmp := t.TempDir()
	classesDir := filepath.Join(tmp, "pkg", "BOOT-INF", "classes")
	if err := os.MkdirAll(classesDir, 0o755); err != nil {
		t.Fatalf("mkdir classes dir: %v", err)
	}
	expected := filepath.Join(classesDir, "bootstrap-dev.yaml")
	if err := os.WriteFile(expected, []byte("spring:\n  datasource:\n    url: test"), 0o644); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	svc := NewArtifactService(tmp)
	found, err := svc.LocateBootstrapDevYAML(tmp)
	if err != nil {
		t.Fatalf("locate bootstrap yaml: %v", err)
	}
	if found != expected {
		t.Fatalf("expected %s, got %s", expected, found)
	}
}
