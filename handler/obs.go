package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nichmahn/go-rc2sb/rc"
	"github.com/nichmahn/go-rc2sb/sb"
)

// NewOBSHandler creates a new Open Bible Stories handler.
func NewOBSHandler() Handler {
	return &obsHandler{}
}

type obsHandler struct{}

func (h *obsHandler) Subject() string {
	return "Open Bible Stories"
}

func (h *obsHandler) Convert(ctx context.Context, manifest *rc.Manifest, inDir, outDir string, opts Options) (*sb.Metadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m := BuildBaseMetadata(manifest, "BurritoTruck", "OBS")

	// Set type - OBS uses gloss/textStories
	m.Type = sb.Type{
		FlavorType: sb.FlavorType{
			Name: "gloss",
			Flavor: sb.Flavor{
				Name: "textStories",
			},
			CurrentScope: map[string][]string{"GEN": {}},
		},
	}

	// OBS uses a different copyright format
	m.Copyright = BuildCopyright(manifest, true)

	// Copy common root files (README.md, .gitignore, .gitea, .github)
	if err := CopyCommonRootFiles(inDir, outDir, m); err != nil {
		return nil, err
	}

	// Copy additional OBS-specific root-level files: LICENSE.md, manifest.yaml, media.yaml
	obsRootFiles := []string{"LICENSE.md", "manifest.yaml", "media.yaml"}
	for _, name := range obsRootFiles {
		src := filepath.Join(inDir, name)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}
		ing, err := CopyFileAndComputeIngredient(src, outDir, name)
		if err != nil {
			return nil, fmt.Errorf("copying root file %s: %w", name, err)
		}
		m.Ingredients[name] = ing
	}

	// Copy content/ directory to ingredients/content/
	contentDir := filepath.Join(inDir, "content")
	if err := copyContentDir(contentDir, outDir, m); err != nil {
		return nil, err
	}

	// Copy LICENSE.md to ingredients/LICENSE.md
	licSrc := filepath.Join(inDir, "LICENSE.md")
	if _, err := os.Stat(licSrc); err == nil {
		ing, err := CopyFileAndComputeIngredient(licSrc, outDir, "ingredients/LICENSE.md")
		if err != nil {
			return nil, fmt.Errorf("copying ingredients/LICENSE.md: %w", err)
		}
		m.Ingredients["ingredients/LICENSE.md"] = ing
	}

	return m, nil
}

// copyContentDir recursively copies content files to ingredients/content/.
func copyContentDir(contentDir, outDir string, m *sb.Metadata) error {
	return filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(contentDir, path)
		if err != nil {
			return err
		}

		ingredientKey := "ingredients/content/" + filepath.ToSlash(relPath)

		ing, err := CopyFileAndComputeIngredient(path, outDir, ingredientKey)
		if err != nil {
			return fmt.Errorf("copying content file %s: %w", relPath, err)
		}
		m.Ingredients[ingredientKey] = ing

		return nil
	})
}
