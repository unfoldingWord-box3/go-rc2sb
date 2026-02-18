package sb

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MIMETypeForExt returns the MIME type for a given file extension.
func MIMETypeForExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".md":
		return "text/markdown"
	case ".usfm":
		return "text/plain"
	case ".tsv":
		return "text/tab-separated-values"
	case ".yaml", ".yml":
		return "text/yaml"
	case ".json":
		return "application/json"
	case ".txt":
		return "text/plain"
	default:
		// Default to text/markdown as the samples show this as fallback
		return "text/markdown"
	}
}

// ComputeIngredient computes the Ingredient (MD5 checksum, size, MIME type) for a file.
func ComputeIngredient(filePath string) (Ingredient, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return Ingredient{}, fmt.Errorf("opening file %s: %w", filePath, err)
	}
	defer f.Close()

	h := md5.New()
	size, err := io.Copy(h, f)
	if err != nil {
		return Ingredient{}, fmt.Errorf("reading file %s: %w", filePath, err)
	}

	ext := filepath.Ext(filePath)

	return Ingredient{
		Checksum: Checksum{
			MD5: fmt.Sprintf("%x", h.Sum(nil)),
		},
		MimeType: MIMETypeForExt(ext),
		Size:     size,
	}, nil
}

// ComputeIngredientWithScope computes the Ingredient and attaches the given scope.
func ComputeIngredientWithScope(filePath string, scope map[string][]string) (Ingredient, error) {
	ing, err := ComputeIngredient(filePath)
	if err != nil {
		return Ingredient{}, err
	}
	ing.Scope = scope
	return ing, nil
}
