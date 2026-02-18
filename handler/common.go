package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nichmahn/go-rc2sb/rc"
	"github.com/nichmahn/go-rc2sb/sb"
)

// CopyFile copies a file from src to dst, creating any necessary directories.
func CopyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", dst, err)
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening source %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("creating destination %s: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copying %s to %s: %w", src, dst, err)
	}

	return out.Close()
}

// CopyFileAndComputeIngredient copies a file and computes its ingredient entry.
// Returns the ingredient key (relative path in SB) and the Ingredient.
func CopyFileAndComputeIngredient(src, outDir, ingredientKey string) (sb.Ingredient, error) {
	dst := filepath.Join(outDir, ingredientKey)
	if err := CopyFile(src, dst); err != nil {
		return sb.Ingredient{}, err
	}
	return sb.ComputeIngredient(dst)
}

// CopyFileWithScope copies a file and computes its ingredient entry with scope.
func CopyFileWithScope(src, outDir, ingredientKey string, scope map[string][]string) (sb.Ingredient, error) {
	dst := filepath.Join(outDir, ingredientKey)
	if err := CopyFile(src, dst); err != nil {
		return sb.Ingredient{}, err
	}
	return sb.ComputeIngredientWithScope(dst, scope)
}

// BuildBaseMetadata creates a base SB Metadata from an RC manifest with common fields set.
func BuildBaseMetadata(manifest *rc.Manifest, idAuthority, abbreviation string) *sb.Metadata {
	m := sb.NewMetadata()

	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	m.Meta.DateCreated = now

	dc := manifest.DublinCore

	// Set ID authority
	if idAuthority == "BurritoTruck" {
		m.IDAuthorities[idAuthority] = sb.IDAuthority{
			ID:   "https://git.door43.org/BurritoTruck",
			Name: map[string]string{"en": "Door43 Burrito Truck"},
		}
	} else {
		m.IDAuthorities[idAuthority] = sb.IDAuthority{
			ID:   "https://git.door43.org/uW",
			Name: map[string]string{"en": "Door43 uW Burritos"},
		}
	}

	// Set identification
	abbr := abbreviation
	if abbr == "" {
		abbr = strings.ToUpper(dc.Identifier)
	}

	m.Identification = sb.Identification{
		Primary: map[string]map[string]sb.PrimaryEntry{
			idAuthority: {
				abbr: {
					Revision:  "1",
					Timestamp: now,
				},
			},
		},
		Name:         map[string]string{"en": dc.Title},
		Description:  map[string]string{"en": dc.Title},
		Abbreviation: map[string]string{"en": abbr},
	}

	// Set language
	m.Languages = []sb.LanguageEntry{
		{
			Tag:             dc.Language.Identifier,
			Name:            map[string]string{"en": dc.Language.Title},
			ScriptDirection: dc.Language.Direction,
		},
	}

	return m
}

// BuildCopyright generates a copyright statement from the RC manifest.
// Uses the format "© {publisher} {year}, {rights}" for most types,
// or "Copyright © {year} by {publisher}" for OBS.
func BuildCopyright(manifest *rc.Manifest, isOBS bool) sb.Copyright {
	dc := manifest.DublinCore
	year := dc.Issued
	if len(year) >= 4 {
		year = year[:4]
	}

	if isOBS {
		return sb.Copyright{
			ShortStatements: []sb.CopyrightStatement{
				{
					Statement: fmt.Sprintf("Copyright \u00a9 %s by %s", year, dc.Publisher),
				},
			},
		}
	}

	return sb.Copyright{
		ShortStatements: []sb.CopyrightStatement{
			{
				Statement: fmt.Sprintf("\u00a9 %s %s, %s", dc.Publisher, year, dc.Rights),
				MimeType:  "text/plain",
				Lang:      "en",
			},
		},
	}
}

// CopyLicenseIngredient copies LICENSE.md from RC to ingredients/LICENSE.md and returns the ingredient.
func CopyLicenseIngredient(inDir, outDir string) (sb.Ingredient, error) {
	src := filepath.Join(inDir, "LICENSE.md")
	return CopyFileAndComputeIngredient(src, outDir, "ingredients/LICENSE.md")
}

// CopyRootFile copies a root-level file from RC to SB root and returns the ingredient.
func CopyRootFile(inDir, outDir, filename string) (sb.Ingredient, error) {
	src := filepath.Join(inDir, filename)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return sb.Ingredient{}, nil // File doesn't exist, skip silently
	}
	return CopyFileAndComputeIngredient(src, outDir, filename)
}

// CopyCommonRootFiles copies common root-level files from the RC repo to the SB output
// if they exist: README.md, .gitea, .github, .gitignore (but NOT .git).
// Files are copied to the SB root and added to the metadata ingredients map.
func CopyCommonRootFiles(inDir, outDir string, m *sb.Metadata) error {
	// Individual files to copy
	files := []string{"README.md", ".gitignore"}
	for _, name := range files {
		src := filepath.Join(inDir, name)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}
		ing, err := CopyFileAndComputeIngredient(src, outDir, name)
		if err != nil {
			return fmt.Errorf("copying root file %s: %w", name, err)
		}
		m.Ingredients[name] = ing
	}

	// Directories to copy recursively
	dirs := []string{".gitea", ".github"}
	for _, dirName := range dirs {
		src := filepath.Join(inDir, dirName)
		info, err := os.Stat(src)
		if os.IsNotExist(err) || !info.IsDir() {
			continue
		}
		if err := copyTreeToIngredients(src, outDir, dirName, m); err != nil {
			return fmt.Errorf("copying root directory %s: %w", dirName, err)
		}
	}

	return nil
}
