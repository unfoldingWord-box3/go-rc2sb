package rc_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nichmahn/go-rc2sb/rc"
)

func TestLoadManifest_MissingFile(t *testing.T) {
	dir := t.TempDir()
	_, err := rc.LoadManifest(dir)
	if err == nil {
		t.Fatal("expected error for missing manifest.yaml")
	}
	if got := err.Error(); got == "" {
		t.Error("error message should not be empty")
	}
}

func TestLoadManifest_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	// Use truly invalid YAML (tab characters in flow context)
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte("{\t\tinvalid:\n\t[broken"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := rc.LoadManifest(dir)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoadManifest_Valid(t *testing.T) {
	dir := t.TempDir()
	yaml := `dublin_core:
  conformsto: 'rc0.2'
  creator: 'TestCreator'
  description: 'Test description'
  format: 'text/markdown'
  identifier: 'test'
  issued: '2024-01-01'
  language:
    direction: 'ltr'
    identifier: 'en'
    title: 'English'
  modified: '2024-01-01'
  publisher: 'TestPublisher'
  rights: 'CC BY-SA 4.0'
  subject: 'Open Bible Stories'
  title: 'Test Title'
  type: 'book'
  version: '1'
checking:
  checking_entity:
    - 'TestEntity'
  checking_level: '3'
projects:
  - identifier: 'obs'
    path: './content'
    sort: 0
    title: 'Test Project'
`
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := rc.LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.DublinCore.Subject != "Open Bible Stories" {
		t.Errorf("Subject = %q; want %q", m.DublinCore.Subject, "Open Bible Stories")
	}
	if m.DublinCore.Identifier != "test" {
		t.Errorf("Identifier = %q; want %q", m.DublinCore.Identifier, "test")
	}
	if m.DublinCore.Language.Identifier != "en" {
		t.Errorf("Language = %q; want %q", m.DublinCore.Language.Identifier, "en")
	}
	if len(m.Projects) != 1 {
		t.Errorf("Projects count = %d; want 1", len(m.Projects))
	}
}
