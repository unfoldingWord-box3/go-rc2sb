package sb_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/unfoldingWord/go-rc2sb/sb"
)

func TestMIMETypeForExt(t *testing.T) {
	tests := []struct {
		ext  string
		want string
	}{
		{".md", "text/markdown"},
		{".MD", "text/markdown"},
		{".usfm", "text/plain"},
		{".tsv", "text/tab-separated-values"},
		{".yaml", "text/yaml"},
		{".yml", "text/yaml"},
		{".json", "application/json"},
		{".txt", "text/plain"},
		{".unknown", "text/markdown"}, // default fallback
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := sb.MIMETypeForExt(tt.ext)
			if got != tt.want {
				t.Errorf("MIMETypeForExt(%q) = %q; want %q", tt.ext, got, tt.want)
			}
		})
	}
}

func TestComputeIngredient(t *testing.T) {
	dir := t.TempDir()
	content := []byte("Hello, World!")
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	ing, err := sb.ComputeIngredient(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ing.Size != 13 {
		t.Errorf("Size = %d; want 13", ing.Size)
	}
	if ing.Checksum.MD5 != "65a8e27d8879283831b664bd8b7f0ad4" {
		t.Errorf("MD5 = %q; want %q", ing.Checksum.MD5, "65a8e27d8879283831b664bd8b7f0ad4")
	}
	if ing.MimeType != "text/markdown" {
		t.Errorf("MimeType = %q; want %q", ing.MimeType, "text/markdown")
	}
}

func TestComputeIngredient_Missing(t *testing.T) {
	_, err := sb.ComputeIngredient("/nonexistent/file.md")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestComputeIngredientWithScope(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.tsv")
	if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	scope := map[string][]string{"GEN": {}}
	ing, err := sb.ComputeIngredientWithScope(path, scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ing.Scope == nil {
		t.Fatal("Scope should not be nil")
	}
	if _, ok := ing.Scope["GEN"]; !ok {
		t.Error("Scope should contain GEN")
	}
}
