package rc2sb_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	rc2sb "github.com/nichmahn/go-rc2sb"
)

func TestConvert_MissingManifest(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	_, err := rc2sb.Convert(context.Background(), inDir, outDir, rc2sb.Options{})
	if err == nil {
		t.Fatal("expected error for missing manifest.yaml")
	}
	if !strings.Contains(err.Error(), "manifest.yaml") {
		t.Errorf("error should mention manifest.yaml: %v", err)
	}
}

func TestConvert_UnsupportedSubject(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	yaml := `dublin_core:
  subject: 'Unknown Subject Type'
  identifier: 'test'
  title: 'Test'
  language:
    identifier: 'en'
    title: 'English'
    direction: 'ltr'
projects: []
`
	if err := os.WriteFile(filepath.Join(inDir, "manifest.yaml"), []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := rc2sb.Convert(context.Background(), inDir, outDir, rc2sb.Options{})
	if err == nil {
		t.Fatal("expected error for unsupported subject")
	}
	if !strings.Contains(err.Error(), "unsupported subject") {
		t.Errorf("error should mention unsupported subject: %v", err)
	}
	if !strings.Contains(err.Error(), "Unknown Subject Type") {
		t.Errorf("error should include the subject name: %v", err)
	}
}

func TestConvert_CancelledContext(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// Create a valid manifest
	yaml := `dublin_core:
  subject: 'Open Bible Stories'
  identifier: 'obs'
  title: 'Test'
  language:
    identifier: 'en'
    title: 'English'
    direction: 'ltr'
projects:
  - identifier: 'obs'
    path: './content'
    sort: 0
    title: 'Test'
`
	if err := os.WriteFile(filepath.Join(inDir, "manifest.yaml"), []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}

	// Create content directory (empty)
	if err := os.MkdirAll(filepath.Join(inDir, "content"), 0755); err != nil {
		t.Fatal(err)
	}

	// Cancel the context before calling Convert
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := rc2sb.Convert(ctx, inDir, outDir, rc2sb.Options{})
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestConvert_InvalidYAML(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(inDir, "manifest.yaml"), []byte(":::bad:::"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := rc2sb.Convert(context.Background(), inDir, outDir, rc2sb.Options{})
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
