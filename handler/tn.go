package handler

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nichmahn/go-rc2sb/books"
	"github.com/nichmahn/go-rc2sb/rc"
	"github.com/nichmahn/go-rc2sb/sb"
)

// NewTNHandler creates a new TSV Translation Notes handler.
func NewTNHandler() Handler {
	return &tnHandler{}
}

type tnHandler struct{}

func (h *tnHandler) Subject() string {
	return "TSV Translation Notes"
}

func (h *tnHandler) Convert(ctx context.Context, manifest *rc.Manifest, inDir, outDir string, opts Options) (*sb.Metadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m := BuildBaseMetadata(manifest, "uWBurritos", "TN")

	// Set type - parascriptural/x-bcvnotes
	currentScope := make(map[string][]string)
	m.Type = sb.Type{
		FlavorType: sb.FlavorType{
			Name: "parascriptural",
			Flavor: sb.Flavor{
				Name: "x-bcvnotes",
			},
		},
	}

	m.Copyright = BuildCopyright(manifest, false)

	// Process each project (TSV file per book)
	for _, project := range manifest.Projects {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		srcPath := filepath.Join(inDir, strings.TrimPrefix(project.Path, "./"))
		srcFilename := filepath.Base(srcPath)

		// Strip "tn_" prefix: "tn_GEN.tsv" -> "GEN.tsv"
		destFilename := strings.TrimPrefix(srcFilename, "tn_")
		ingredientKey := "ingredients/" + destFilename

		// Get book code for scope
		bookID := strings.ToLower(project.Identifier)
		bookCode := books.CodeFromProjectID(bookID)

		scope := map[string][]string{bookCode: {}}
		currentScope[bookCode] = []string{}

		// Add localized name
		key, localizedName := books.LocalizedNameEntry(bookID)
		if key != "" {
			m.LocalizedNames[key] = localizedName
		}

		// Copy TSV file with scope
		ing, err := CopyFileWithScope(srcPath, outDir, ingredientKey, scope)
		if err != nil {
			return nil, fmt.Errorf("copying %s: %w", srcFilename, err)
		}
		m.Ingredients[ingredientKey] = ing
	}

	// Set the currentScope
	m.Type.FlavorType.CurrentScope = currentScope

	// Copy common root files (README.md, .gitignore, .gitea, .github)
	if err := CopyCommonRootFiles(inDir, outDir, m); err != nil {
		return nil, err
	}

	// Copy LICENSE.md to ingredients/
	licIng, err := CopyLicenseIngredient(inDir, outDir)
	if err != nil {
		return nil, fmt.Errorf("copying LICENSE.md: %w", err)
	}
	m.Ingredients["ingredients/LICENSE.md"] = licIng

	return m, nil
}
