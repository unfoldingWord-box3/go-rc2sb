package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unfoldingWord/go-rc2sb/rc"
	"github.com/unfoldingWord/go-rc2sb/sb"
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

	// Copy LICENSE.md to root (uses embedded default if RC doesn't have one).
	if err := CopyLicenseToRoot(inDir, outDir); err != nil {
		return nil, fmt.Errorf("copying root LICENSE.md: %w", err)
	}

	// Determine the content directory from the manifest project path.
	// OBS has a single project whose path is typically "./content" but may be "."
	// when the markdown files live in the repository root.
	contentPath := "content"
	if len(manifest.Projects) > 0 {
		p := strings.TrimPrefix(manifest.Projects[0].Path, "./")
		if p != "" {
			contentPath = p
		}
	}

	if contentPath == "." {
		// Content lives in the repo root — copy everything except known
		// non-content files (manifest.yaml, media.yaml, README.md, LICENSE.md,
		// .gitignore, and dot-directories like .git, .gitea, .github).
		if err := copyOBSRootContent(inDir, outDir, m); err != nil {
			return nil, err
		}
	} else {
		// Content lives in a subdirectory — copy everything in it.
		contentDir := filepath.Join(inDir, contentPath)
		if err := copyContentDir(contentDir, outDir, m); err != nil {
			return nil, err
		}
	}

	// Copy LICENSE.md to ingredients/LICENSE.md (uses embedded default if RC doesn't have one).
	licIng, err := CopyLicenseIngredient(inDir, outDir)
	if err != nil {
		return nil, fmt.Errorf("copying ingredients/LICENSE.md: %w", err)
	}
	m.Ingredients["ingredients/LICENSE.md"] = licIng

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

// copyOBSRootContent copies OBS content from the repo root when the manifest
// project path is ".". It copies all files and directories except known
// non-content entries: *.yaml files, README.md, LICENSE.md, .gitignore,
// and dot-directories (.git, .gitea, .github). This handles both flat layouts
// (numbered .md files, front.md, back.md) and layouts with subdirectories
// (front/, back/).
func copyOBSRootContent(inDir, outDir string, m *sb.Metadata) error {
	entries, err := os.ReadDir(inDir)
	if err != nil {
		return fmt.Errorf("reading OBS root directory: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()

		if isOBSExcludedEntry(name, entry.IsDir()) {
			continue
		}

		srcPath := filepath.Join(inDir, name)

		if entry.IsDir() {
			// Recursively copy the subdirectory into ingredients/content/{dir}/
			// We walk the subdirectory and prefix each relative path with the
			// directory name so that e.g. front/intro.md maps to
			// ingredients/content/front/intro.md.
			if err := copyOBSSubdir(srcPath, name, outDir, m); err != nil {
				return fmt.Errorf("copying OBS content directory %s: %w", name, err)
			}
		} else {
			ingredientKey := "ingredients/content/" + name
			ing, err := CopyFileAndComputeIngredient(srcPath, outDir, ingredientKey)
			if err != nil {
				return fmt.Errorf("copying OBS content file %s: %w", name, err)
			}
			m.Ingredients[ingredientKey] = ing
		}
	}

	return nil
}

// copyOBSSubdir recursively copies a subdirectory from the OBS root into
// ingredients/content/{dirName}/. For example, a file front/intro.md is
// copied to ingredients/content/front/intro.md.
func copyOBSSubdir(srcDir, dirName, outDir string, m *sb.Metadata) error {
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

		ingredientKey := "ingredients/content/" + dirName + "/" + filepath.ToSlash(relPath)

		ing, err := CopyFileAndComputeIngredient(path, outDir, ingredientKey)
		if err != nil {
			return fmt.Errorf("copying %s/%s: %w", dirName, relPath, err)
		}
		m.Ingredients[ingredientKey] = ing

		return nil
	})
}

// isOBSExcludedEntry returns true if the given root-level entry should be
// excluded from OBS content copying. Excluded entries are repository metadata
// and infrastructure files that are not part of the OBS content itself.
func isOBSExcludedEntry(name string, isDir bool) bool {
	if isDir {
		// Exclude dot-directories (.git, .gitea, .github, etc.)
		return strings.HasPrefix(name, ".")
	}
	// Exclude YAML metadata files
	if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
		return true
	}
	// Exclude known root-level non-content files
	switch name {
	case "README.md", "LICENSE.md", ".gitignore":
		return true
	}
	return false
}
