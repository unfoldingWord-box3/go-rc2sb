package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nichmahn/go-rc2sb/books"
	"github.com/nichmahn/go-rc2sb/rc"
	"github.com/nichmahn/go-rc2sb/sb"
)

// NewTQHandler creates a new TSV Translation Questions handler.
func NewTQHandler() Handler {
	return &tqHandler{}
}

type tqHandler struct{}

func (h *tqHandler) Subject() string {
	return "TSV Translation Questions"
}

func (h *tqHandler) Convert(ctx context.Context, manifest *rc.Manifest, inDir, outDir string, opts Options) (*sb.Metadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m := BuildBaseMetadata(manifest, "uWBurritos", "TQ")

	// Set type - parascriptural/x-bcvquestions
	currentScope := make(map[string][]string)
	m.Type = sb.Type{
		FlavorType: sb.FlavorType{
			Name: "parascriptural",
			Flavor: sb.Flavor{
				Name: "x-bcvquestions",
			},
		},
	}

	m.Copyright = BuildCopyright(manifest, false)

	lang := manifest.DublinCore.Language.Identifier

	// Process each project (TSV file per book)
	for _, project := range manifest.Projects {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		srcPath := filepath.Join(inDir, strings.TrimPrefix(project.Path, "./"))
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			continue
		}
		srcFilename := filepath.Base(srcPath)

		// Strip "tq_" prefix: "tq_GEN.tsv" -> "GEN.tsv"
		destFilename := strings.TrimPrefix(srcFilename, "tq_")
		ingredientKey := "ingredients/" + destFilename

		// Get book code for scope
		bookID := strings.ToLower(project.Identifier)
		bookCode := books.CodeFromProjectID(bookID)

		scope := map[string][]string{bookCode: {}}
		currentScope[bookCode] = []string{}

		// Add localized name: try USFM from USFMPath, then manifest title, then English
		var usfmNames *books.LocalizedBookNames
		if opts.USFMPath != "" {
			if usfmFile := books.FindUSFMFile(opts.USFMPath, bookID); usfmFile != "" {
				usfmNames = books.ParseUSFMBookNames(usfmFile)
			}
		}
		key, localizedName := books.LocalizedNameEntryWithNames(bookID, lang, project.Title, usfmNames)
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
