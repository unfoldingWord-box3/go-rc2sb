package handler

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nichmahn/go-rc2sb/books"
	"github.com/nichmahn/go-rc2sb/rc"
	"github.com/nichmahn/go-rc2sb/sb"
)

// twLinkRegexp parses RC links like "rc://*/tw/dict/bible/other/creation"
var twLinkRegexp = regexp.MustCompile(`rc://[^/]*/tw/dict/bible/([^/]+)/([^/]+)`)

// NewTWLHandler creates a new TSV Translation Words Links handler.
func NewTWLHandler() Handler {
	return &twlHandler{}
}

type twlHandler struct{}

func (h *twlHandler) Subject() string {
	return "TSV Translation Words Links"
}

func (h *twlHandler) Convert(ctx context.Context, manifest *rc.Manifest, inDir, outDir string, opts Options) (*sb.Metadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m := BuildBaseMetadata(manifest, "uWBurritos", "TW")

	// Set type - parascriptural/x-bcvarticles
	currentScope := make(map[string][]string)
	m.Type = sb.Type{
		FlavorType: sb.FlavorType{
			Name: "parascriptural",
			Flavor: sb.Flavor{
				Name: "x-bcvarticles",
			},
		},
	}

	m.Copyright = BuildCopyright(manifest, false)

	// Collect all TWLink references for payload processing
	twLinks := make(map[string]bool) // "category/article" -> true

	// Process each project (TSV file per book)
	for _, project := range manifest.Projects {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		srcPath := filepath.Join(inDir, strings.TrimPrefix(project.Path, "./"))
		srcFilename := filepath.Base(srcPath)

		// Strip "twl_" prefix: "twl_GEN.tsv" -> "GEN.tsv"
		destFilename := strings.TrimPrefix(srcFilename, "twl_")
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

		// Extract TWLink references from the TSV file
		links, err := extractTWLinks(srcPath)
		if err != nil {
			return nil, fmt.Errorf("extracting TWLinks from %s: %w", srcFilename, err)
		}
		for _, link := range links {
			twLinks[link] = true
		}
	}

	// Set the currentScope
	m.Type.FlavorType.CurrentScope = currentScope

	// Process payload if Translation Words directory is available
	twDir, hasTW := opts.PayloadDirs["Translation Words"]
	if hasTW && twDir != "" {
		if err := processTWPayload(ctx, twDir, outDir, twLinks, m); err != nil {
			return nil, fmt.Errorf("processing TW payload: %w", err)
		}
	}

	// Copy LICENSE.md to ingredients/
	licIng, err := CopyLicenseIngredient(inDir, outDir)
	if err != nil {
		return nil, fmt.Errorf("copying LICENSE.md: %w", err)
	}
	m.Ingredients["ingredients/LICENSE.md"] = licIng

	return m, nil
}

// extractTWLinks reads a TSV file and extracts all unique TWLink RC references.
// Returns a slice of "category/article" strings (e.g., "other/creation").
func extractTWLinks(tsvPath string) ([]string, error) {
	f, err := os.Open(tsvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer for large lines

	// Read header to find TWLink column index
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty TSV file: %s", tsvPath)
	}
	header := scanner.Text()
	cols := strings.Split(header, "\t")

	twLinkCol := -1
	for i, col := range cols {
		if strings.TrimSpace(col) == "TWLink" {
			twLinkCol = i
			break
		}
	}
	if twLinkCol < 0 {
		return nil, nil // No TWLink column, nothing to extract
	}

	seen := make(map[string]bool)
	var links []string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if twLinkCol >= len(fields) {
			continue
		}

		twLink := fields[twLinkCol]
		matches := twLinkRegexp.FindAllStringSubmatch(twLink, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				key := match[1] + "/" + match[2] // "category/article"
				if !seen[key] {
					seen[key] = true
					links = append(links, key)
				}
			}
		}
	}

	return links, scanner.Err()
}

// processTWPayload copies referenced Translation Words articles to the payload directory.
func processTWPayload(ctx context.Context, twDir, outDir string, twLinks map[string]bool, m *sb.Metadata) error {
	// The TW RC has articles in bible/{category}/{article}.md
	bibleDir := filepath.Join(twDir, "bible")

	for link := range twLinks {
		if err := ctx.Err(); err != nil {
			return err
		}

		// link is "category/article" (e.g., "other/creation")
		srcFile := filepath.Join(bibleDir, link+".md")
		if _, err := os.Stat(srcFile); os.IsNotExist(err) {
			// Article not found, skip
			continue
		}

		ingredientKey := "ingredients/payload/" + link + ".md"
		ing, err := CopyFileAndComputeIngredient(srcFile, outDir, ingredientKey)
		if err != nil {
			return fmt.Errorf("copying TW article %s: %w", link, err)
		}
		m.Ingredients[ingredientKey] = ing
	}

	return nil
}
