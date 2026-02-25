package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/unfoldingWord/go-rc2sb/rc"
	"github.com/unfoldingWord/go-rc2sb/sb"
)

// NewTWHandler creates a new Translation Words handler.
func NewTWHandler() Handler {
	return &twHandler{}
}

type twHandler struct{}

func (h *twHandler) Subject() string {
	return "Translation Words"
}

func (h *twHandler) Convert(ctx context.Context, manifest *rc.Manifest, inDir, outDir string, opts Options) (*sb.Metadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m := BuildBaseMetadata(manifest, "uWBurritos", "TW")

	// Set type - peripheral/x-peripheralArticles
	m.Type = sb.Type{
		FlavorType: sb.FlavorType{
			Name: "peripheral",
			Flavor: sb.Flavor{
				Name: "x-peripheralArticles",
			},
		},
	}

	m.Copyright = BuildCopyright(manifest, false)
	m.LocalizedNames = map[string]sb.LocalizedName{}

	// Copy common root files (README.md, .gitignore, .gitea, .github)
	if err := CopyCommonRootFiles(inDir, outDir, m); err != nil {
		return nil, err
	}

	// Copy LICENSE.md to root (uses embedded default if RC doesn't have one).
	if err := CopyLicenseToRoot(inDir, outDir); err != nil {
		return nil, fmt.Errorf("copying root LICENSE.md: %w", err)
	}

	// Copy bible/ contents to ingredients/
	// Structure: bible/{kt,other,names}/*.md and bible/config.yaml
	bibleDir := filepath.Join(inDir, "bible")
	if err := copyTreeToIngredients(bibleDir, outDir, "ingredients", m); err != nil {
		return nil, fmt.Errorf("copying bible directory: %w", err)
	}

	// Copy LICENSE.md to ingredients/LICENSE.md (uses embedded default if RC doesn't have one).
	licIng, err := CopyLicenseIngredient(inDir, outDir)
	if err != nil {
		return nil, fmt.Errorf("copying ingredients/LICENSE.md: %w", err)
	}
	m.Ingredients["ingredients/LICENSE.md"] = licIng

	return m, nil
}

// copyTreeToIngredients recursively copies a directory tree into the ingredients directory.
func copyTreeToIngredients(srcDir, outDir, destPrefix string, m *sb.Metadata) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		ingredientKey := destPrefix + "/" + filepath.ToSlash(relPath)

		ing, err := CopyFileAndComputeIngredient(path, outDir, ingredientKey)
		if err != nil {
			return fmt.Errorf("copying %s: %w", relPath, err)
		}
		m.Ingredients[ingredientKey] = ing

		return nil
	})
}
