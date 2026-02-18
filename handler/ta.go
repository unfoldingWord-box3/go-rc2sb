package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nichmahn/go-rc2sb/rc"
	"github.com/nichmahn/go-rc2sb/sb"
)

// NewTAHandler creates a new Translation Academy handler.
func NewTAHandler() Handler {
	return &taHandler{}
}

type taHandler struct{}

func (h *taHandler) Subject() string {
	return "Translation Academy"
}

func (h *taHandler) Convert(ctx context.Context, manifest *rc.Manifest, inDir, outDir string, opts Options) (*sb.Metadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m := BuildBaseMetadata(manifest, "uWBurritos", "TA")

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

	// Copy root-level files: .gitignore, LICENSE.md, README.md, manifest.yaml, media.yaml
	rootFiles := []string{".gitignore", "LICENSE.md", "README.md", "manifest.yaml", "media.yaml"}
	for _, name := range rootFiles {
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

	// Copy each project directory to ingredients/
	// Projects are: intro, process, translate, checking
	for _, project := range manifest.Projects {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		projectDir := filepath.Join(inDir, project.Identifier)
		if _, err := os.Stat(projectDir); os.IsNotExist(err) {
			continue
		}

		destPrefix := "ingredients/" + project.Identifier
		if err := copyTreeToIngredients(projectDir, outDir, destPrefix, m); err != nil {
			return nil, fmt.Errorf("copying project %s: %w", project.Identifier, err)
		}
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
