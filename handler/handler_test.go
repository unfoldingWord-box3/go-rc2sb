package handler_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nichmahn/go-rc2sb/handler"
	"github.com/nichmahn/go-rc2sb/rc"
	"github.com/nichmahn/go-rc2sb/sb"

	// Register all handlers so Lookup works.
	_ "github.com/nichmahn/go-rc2sb/handler/subjects"
)

// --- CopyCommonRootFiles tests ---

func TestCopyCommonRootFiles_CopiesREADMEAndGitignore(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// Create README.md and .gitignore in the input directory
	if err := os.WriteFile(filepath.Join(inDir, "README.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, ".gitignore"), []byte("*.tmp\n"), 0644); err != nil {
		t.Fatal(err)
	}

	m := sb.NewMetadata()
	if err := handler.CopyCommonRootFiles(inDir, outDir, m); err != nil {
		t.Fatalf("CopyCommonRootFiles failed: %v", err)
	}

	// Verify README.md was copied and is not in ingredients
	if _, err := os.Stat(filepath.Join(outDir, "README.md")); os.IsNotExist(err) {
		t.Error("README.md was not copied to outDir")
	}
	if _, ok := m.Ingredients["README.md"]; ok {
		t.Error("README.md should not be in metadata ingredients")
	}

	// Verify .gitignore was copied and is not in ingredients
	if _, err := os.Stat(filepath.Join(outDir, ".gitignore")); os.IsNotExist(err) {
		t.Error(".gitignore was not copied to outDir")
	}
	if _, ok := m.Ingredients[".gitignore"]; ok {
		t.Error(".gitignore should not be in metadata ingredients")
	}
}

func TestCopyCommonRootFiles_CopiesGiteaDir(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// Create .gitea/ directory with a file
	giteaDir := filepath.Join(inDir, ".gitea")
	if err := os.MkdirAll(giteaDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(giteaDir, "auto_merge.yaml"), []byte("auto_merge: true\n"), 0644); err != nil {
		t.Fatal(err)
	}

	m := sb.NewMetadata()
	if err := handler.CopyCommonRootFiles(inDir, outDir, m); err != nil {
		t.Fatalf("CopyCommonRootFiles failed: %v", err)
	}

	// Verify the .gitea directory was copied
	if _, err := os.Stat(filepath.Join(outDir, ".gitea", "auto_merge.yaml")); os.IsNotExist(err) {
		t.Error(".gitea/auto_merge.yaml was not copied to outDir")
	}
	if _, ok := m.Ingredients[".gitea/auto_merge.yaml"]; ok {
		t.Error(".gitea/auto_merge.yaml should not be in metadata ingredients")
	}
}

func TestCopyCommonRootFiles_CopiesGithubDir(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// Create .github/ directory with nested structure
	workflowsDir := filepath.Join(inDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowsDir, "ci.yml"), []byte("name: CI\n"), 0644); err != nil {
		t.Fatal(err)
	}

	m := sb.NewMetadata()
	if err := handler.CopyCommonRootFiles(inDir, outDir, m); err != nil {
		t.Fatalf("CopyCommonRootFiles failed: %v", err)
	}

	// Verify the .github directory was copied recursively
	if _, err := os.Stat(filepath.Join(outDir, ".github", "workflows", "ci.yml")); os.IsNotExist(err) {
		t.Error(".github/workflows/ci.yml was not copied to outDir")
	}
	if _, ok := m.Ingredients[".github/workflows/ci.yml"]; ok {
		t.Error(".github/workflows/ci.yml should not be in metadata ingredients")
	}
}

func TestCopyCommonRootFiles_SkipsMissingFiles(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// No files at all — should succeed without copying anything
	m := sb.NewMetadata()
	if err := handler.CopyCommonRootFiles(inDir, outDir, m); err != nil {
		t.Fatalf("CopyCommonRootFiles should not fail when no root files exist: %v", err)
	}

	if len(m.Ingredients) != 0 {
		t.Errorf("Expected 0 ingredients for empty inDir, got %d", len(m.Ingredients))
	}
}

func TestCopyCommonRootFiles_DoesNotCopyGitDir(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// Create .git/ directory (should NOT be copied)
	gitDir := filepath.Join(inDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte("[core]\n"), 0644); err != nil {
		t.Fatal(err)
	}

	m := sb.NewMetadata()
	if err := handler.CopyCommonRootFiles(inDir, outDir, m); err != nil {
		t.Fatalf("CopyCommonRootFiles failed: %v", err)
	}

	// Verify .git was NOT copied
	if _, err := os.Stat(filepath.Join(outDir, ".git")); !os.IsNotExist(err) {
		t.Error(".git directory should NOT be copied to outDir")
	}
	for key := range m.Ingredients {
		if strings.HasPrefix(key, ".git/") {
			t.Errorf("ingredient key %q should not start with .git/", key)
		}
	}
}

func TestCopyCommonRootFiles_DoesNotAddRootFilesToIngredients(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	content := []byte("Hello, World!")
	if err := os.WriteFile(filepath.Join(inDir, "README.md"), content, 0644); err != nil {
		t.Fatal(err)
	}

	m := sb.NewMetadata()
	if err := handler.CopyCommonRootFiles(inDir, outDir, m); err != nil {
		t.Fatalf("CopyCommonRootFiles failed: %v", err)
	}

	if len(m.Ingredients) != 0 {
		t.Errorf("Expected no root-file ingredient entries, got %d", len(m.Ingredients))
	}
}

// --- Bible subject alias tests ---

func TestBibleSubjectAliases_AllRegistered(t *testing.T) {
	subjects := []string{
		"Aligned Bible",
		"Bible",
		"Hebrew Old Testament",
		"Greek New Testament",
	}

	for _, subject := range subjects {
		t.Run(subject, func(t *testing.T) {
			h, err := handler.Lookup(subject)
			if err != nil {
				t.Fatalf("Lookup(%q) failed: %v", subject, err)
			}
			if h.Subject() != subject {
				t.Errorf("Subject() = %q; want %q", h.Subject(), subject)
			}
		})
	}
}

func TestBibleSubjectAliases_AbbreviationFromIdentifier(t *testing.T) {
	tests := []struct {
		subject    string
		identifier string
		wantAbbr   string
	}{
		{"Aligned Bible", "ult", "ULT"},
		{"Bible", "ust", "UST"},
		{"Hebrew Old Testament", "uhb", "UHB"},
		{"Greek New Testament", "ugnt", "UGNT"},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			h, err := handler.Lookup(tt.subject)
			if err != nil {
				t.Fatalf("Lookup(%q) failed: %v", tt.subject, err)
			}

			// Create a minimal RC structure to test abbreviation derivation
			inDir := t.TempDir()
			outDir := t.TempDir()

			// Write a minimal manifest
			manifest := &rc.Manifest{
				DublinCore: rc.DublinCore{
					Subject:    tt.subject,
					Identifier: tt.identifier,
					Title:      "Test " + tt.subject,
					Issued:     "2024-01-01",
					Publisher:  "test",
					Rights:     "CC BY-SA 4.0",
					Language: rc.Language{
						Identifier: "en",
						Title:      "English",
						Direction:  "ltr",
					},
				},
			}

			// Create a minimal USFM file for the handler to process
			os.MkdirAll(filepath.Join(inDir, "content"), 0755)
			os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644)

			metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, handler.Options{})
			if err != nil {
				t.Fatalf("Convert failed: %v", err)
			}

			gotAbbr := metadata.Identification.Abbreviation["en"]
			if gotAbbr != tt.wantAbbr {
				t.Errorf("Abbreviation = %q; want %q", gotAbbr, tt.wantAbbr)
			}
		})
	}
}

// --- TWL handler tests ---

func writeTWLManifest(t *testing.T, inDir string) *rc.Manifest {
	t.Helper()
	return &rc.Manifest{
		DublinCore: rc.DublinCore{
			Subject:    "TSV Translation Words Links",
			Identifier: "twl",
			Title:      "Test TWL",
			Issued:     "2024-01-01",
			Publisher:  "test",
			Rights:     "CC BY-SA 4.0",
			Language: rc.Language{
				Identifier: "en",
				Title:      "English",
				Direction:  "ltr",
			},
		},
		Projects: []rc.Project{
			{
				Identifier: "gen",
				Path:       "./twl_GEN.tsv",
				Sort:       1,
				Title:      "Genesis",
			},
		},
	}
}

func TestTWL_AutoDetectsPayload(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	manifest := writeTWLManifest(t, inDir)

	// Create the TWL TSV file with an rc:// link
	tsvContent := "Reference\tID\tTags\tOrigWords\tOccurrence\tTWLink\n" +
		"1:1\tabcd\t\tword\t1\trc://*/tw/dict/bible/names/adam\n"
	os.WriteFile(filepath.Join(inDir, "twl_GEN.tsv"), []byte(tsvContent), 0644)
	os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644)

	// Create the en_tw/bible/ directory (auto-detection target)
	twBibleDir := filepath.Join(inDir, "en_tw", "bible", "names")
	os.MkdirAll(twBibleDir, 0755)
	os.WriteFile(filepath.Join(twBibleDir, "adam.md"), []byte("# Adam\n\nThe first man."), 0644)

	h, err := handler.Lookup("TSV Translation Words Links")
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, handler.Options{})
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify payload was auto-detected and copied
	if _, ok := metadata.Ingredients["ingredients/payload/names/adam.md"]; !ok {
		t.Error("Payload article ingredients/payload/names/adam.md not found; auto-detection failed")
	}

	// Verify TSV was rewritten
	data, err := os.ReadFile(filepath.Join(outDir, "ingredients", "GEN.tsv"))
	if err != nil {
		t.Fatalf("Reading output TSV: %v", err)
	}
	content := string(data)
	if strings.Contains(content, "rc://") {
		t.Error("TSV still contains rc:// links after auto-detection rewrite")
	}
	if !strings.Contains(content, "./payload/names/adam.md") {
		t.Error("TSV does not contain expected ./payload/names/adam.md path")
	}
}

func TestTWL_ExplicitPayloadPath(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()
	payloadDir := t.TempDir() // Separate directory for payload

	manifest := writeTWLManifest(t, inDir)

	// Create the TWL TSV file with an rc:// link
	tsvContent := "Reference\tID\tTags\tOrigWords\tOccurrence\tTWLink\n" +
		"1:1\tabcd\t\tword\t1\trc://*/tw/dict/bible/kt/god\n"
	os.WriteFile(filepath.Join(inDir, "twl_GEN.tsv"), []byte(tsvContent), 0644)
	os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644)

	// Create the TW directory at the explicit payload path
	twBibleDir := filepath.Join(payloadDir, "bible", "kt")
	os.MkdirAll(twBibleDir, 0755)
	os.WriteFile(filepath.Join(twBibleDir, "god.md"), []byte("# God\n\nThe creator."), 0644)

	h, err := handler.Lookup("TSV Translation Words Links")
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	opts := handler.Options{PayloadPath: payloadDir}
	metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, opts)
	if err != nil {
		t.Fatalf("Convert with PayloadPath failed: %v", err)
	}

	// Verify payload from explicit path was copied
	if _, ok := metadata.Ingredients["ingredients/payload/kt/god.md"]; !ok {
		t.Error("Payload article ingredients/payload/kt/god.md not found; explicit PayloadPath failed")
	}

	// Verify TSV was rewritten
	data, err := os.ReadFile(filepath.Join(outDir, "ingredients", "GEN.tsv"))
	if err != nil {
		t.Fatalf("Reading output TSV: %v", err)
	}
	content := string(data)
	if strings.Contains(content, "rc://") {
		t.Error("TSV still contains rc:// links after PayloadPath rewrite")
	}
	if !strings.Contains(content, "./payload/kt/god.md") {
		t.Error("TSV does not contain expected ./payload/kt/god.md path")
	}
}

func TestTWL_NoPayloadCopiesAsIs(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	manifest := writeTWLManifest(t, inDir)

	// Create the TWL TSV file with an rc:// link — but NO en_tw/ directory
	tsvContent := "Reference\tID\tTags\tOrigWords\tOccurrence\tTWLink\n" +
		"1:1\tabcd\t\tword\t1\trc://*/tw/dict/bible/names/adam\n"
	os.WriteFile(filepath.Join(inDir, "twl_GEN.tsv"), []byte(tsvContent), 0644)
	os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644)

	h, err := handler.Lookup("TSV Translation Words Links")
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, handler.Options{})
	if err != nil {
		t.Fatalf("Convert without payload failed: %v", err)
	}

	// Verify no payload ingredients
	for key := range metadata.Ingredients {
		if strings.HasPrefix(key, "ingredients/payload/") {
			t.Errorf("Unexpected payload ingredient %s when no TW directory exists", key)
		}
	}

	// Verify TSV was copied as-is (rc:// links preserved)
	data, err := os.ReadFile(filepath.Join(outDir, "ingredients", "GEN.tsv"))
	if err != nil {
		t.Fatalf("Reading output TSV: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "rc://") {
		t.Error("TSV should preserve rc:// links when no payload exists")
	}
	if strings.Contains(content, "./payload/") {
		t.Error("TSV should NOT contain ./payload/ paths when no payload exists")
	}
}

func TestTWL_LinkRewriteMultipleLinks(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	manifest := writeTWLManifest(t, inDir)

	// Create a TSV with multiple rc:// links across several rows
	tsvContent := "Reference\tID\tTags\tOrigWords\tOccurrence\tTWLink\n" +
		"1:1\ta001\t\tword1\t1\trc://*/tw/dict/bible/names/adam\n" +
		"1:2\ta002\t\tword2\t1\trc://*/tw/dict/bible/kt/god\n" +
		"1:3\ta003\t\tword3\t1\trc://en/tw/dict/bible/other/creation\n"
	os.WriteFile(filepath.Join(inDir, "twl_GEN.tsv"), []byte(tsvContent), 0644)
	os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644)

	// Create the en_tw/bible/ directory
	for _, path := range []string{"names/adam.md", "kt/god.md", "other/creation.md"} {
		fullPath := filepath.Join(inDir, "en_tw", "bible", path)
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		os.WriteFile(fullPath, []byte("# Article\n"), 0644)
	}

	h, err := handler.Lookup("TSV Translation Words Links")
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, handler.Options{})
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify all three payload articles were copied
	expectedPayload := []string{
		"ingredients/payload/names/adam.md",
		"ingredients/payload/kt/god.md",
		"ingredients/payload/other/creation.md",
	}
	for _, key := range expectedPayload {
		if _, ok := metadata.Ingredients[key]; !ok {
			t.Errorf("Missing payload ingredient: %s", key)
		}
	}

	// Verify all rc:// links were rewritten
	data, err := os.ReadFile(filepath.Join(outDir, "ingredients", "GEN.tsv"))
	if err != nil {
		t.Fatalf("Reading output TSV: %v", err)
	}
	content := string(data)
	if strings.Contains(content, "rc://") {
		t.Error("TSV still contains rc:// links — not all were rewritten")
	}

	// Verify specific rewrites
	expectedPaths := []string{
		"./payload/names/adam.md",
		"./payload/kt/god.md",
		"./payload/other/creation.md",
	}
	for _, p := range expectedPaths {
		if !strings.Contains(content, p) {
			t.Errorf("TSV missing expected rewritten path: %s", p)
		}
	}
}

func TestTWL_StripsTWLPrefix(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	manifest := writeTWLManifest(t, inDir)

	tsvContent := "Reference\tID\tTags\tOrigWords\tOccurrence\tTWLink\n" +
		"1:1\ta001\t\tword1\t1\trc://*/tw/dict/bible/names/adam\n"
	os.WriteFile(filepath.Join(inDir, "twl_GEN.tsv"), []byte(tsvContent), 0644)
	os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644)

	h, err := handler.Lookup("TSV Translation Words Links")
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, handler.Options{})
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify twl_ prefix was stripped: "twl_GEN.tsv" -> "ingredients/GEN.tsv"
	if _, ok := metadata.Ingredients["ingredients/GEN.tsv"]; !ok {
		t.Error("Expected ingredient key 'ingredients/GEN.tsv' (twl_ prefix should be stripped)")
	}

	// Verify the file exists on disk with the stripped name
	if _, err := os.Stat(filepath.Join(outDir, "ingredients", "GEN.tsv")); os.IsNotExist(err) {
		t.Error("ingredients/GEN.tsv file does not exist on disk")
	}
}

func TestTWL_CopiesRootFilesWithoutIngredientEntries(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	manifest := writeTWLManifest(t, inDir)

	tsvContent := "Reference\tID\tTags\tOrigWords\tOccurrence\tTWLink\n" +
		"1:1\ta001\t\tword1\t1\trc://*/tw/dict/bible/names/adam\n"
	os.WriteFile(filepath.Join(inDir, "twl_GEN.tsv"), []byte(tsvContent), 0644)
	os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644)
	os.WriteFile(filepath.Join(inDir, "README.md"), []byte("# TWL Readme"), 0644)
	os.WriteFile(filepath.Join(inDir, ".gitignore"), []byte("*.tmp\n"), 0644)

	h, err := handler.Lookup("TSV Translation Words Links")
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, handler.Options{})
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify root files are not in ingredients metadata
	if _, ok := metadata.Ingredients["README.md"]; ok {
		t.Error("README.md should not be present in TWL metadata ingredients")
	}
	if _, ok := metadata.Ingredients[".gitignore"]; ok {
		t.Error(".gitignore should not be present in TWL metadata ingredients")
	}

	// Verify files exist on disk
	if _, err := os.Stat(filepath.Join(outDir, "README.md")); os.IsNotExist(err) {
		t.Error("README.md was not copied to TWL output")
	}
	if _, err := os.Stat(filepath.Join(outDir, ".gitignore")); os.IsNotExist(err) {
		t.Error(".gitignore was not copied to TWL output")
	}
}

func TestTA_DoesNotCopyManifestOrMediaToRoot(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	manifest := &rc.Manifest{
		DublinCore: rc.DublinCore{
			Subject:    "Translation Academy",
			Identifier: "ta",
			Title:      "Test TA",
			Issued:     "2024-01-01",
			Publisher:  "test",
			Rights:     "CC BY-SA 4.0",
			Language: rc.Language{
				Identifier: "en",
				Title:      "English",
				Direction:  "ltr",
			},
		},
		Projects: []rc.Project{
			{Identifier: "intro"},
		},
	}

	if err := os.MkdirAll(filepath.Join(inDir, "intro"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, "intro", "01.md"), []byte("# Intro"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, "manifest.yaml"), []byte("dublin_core: {}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, "media.yaml"), []byte("projects: []"), 0644); err != nil {
		t.Fatal(err)
	}

	h, err := handler.Lookup("Translation Academy")
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, handler.Options{})
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "LICENSE.md")); os.IsNotExist(err) {
		t.Error("LICENSE.md should be copied to TA output root")
	}
	if _, err := os.Stat(filepath.Join(outDir, "manifest.yaml")); !os.IsNotExist(err) {
		t.Error("manifest.yaml should not be copied to TA output root")
	}
	if _, err := os.Stat(filepath.Join(outDir, "media.yaml")); !os.IsNotExist(err) {
		t.Error("media.yaml should not be copied to TA output root")
	}
	if _, ok := metadata.Ingredients["ingredients/LICENSE.md"]; !ok {
		t.Error("ingredients/LICENSE.md should exist in TA metadata ingredients")
	}
}

func TestOBS_DoesNotCopyManifestOrMediaToRoot(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	manifest := &rc.Manifest{
		DublinCore: rc.DublinCore{
			Subject:    "Open Bible Stories",
			Identifier: "obs",
			Title:      "Test OBS",
			Issued:     "2024-01-01",
			Publisher:  "test",
			Rights:     "CC BY-SA 4.0",
			Language: rc.Language{
				Identifier: "en",
				Title:      "English",
				Direction:  "ltr",
			},
		},
	}

	if err := os.MkdirAll(filepath.Join(inDir, "content"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, "content", "01.md"), []byte("# Story"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, "LICENSE.md"), []byte("License"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, "manifest.yaml"), []byte("dublin_core: {}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inDir, "media.yaml"), []byte("projects: []"), 0644); err != nil {
		t.Fatal(err)
	}

	h, err := handler.Lookup("Open Bible Stories")
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	metadata, err := h.Convert(context.Background(), manifest, inDir, outDir, handler.Options{})
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "LICENSE.md")); os.IsNotExist(err) {
		t.Error("LICENSE.md should be copied to OBS output root")
	}
	if _, err := os.Stat(filepath.Join(outDir, "manifest.yaml")); !os.IsNotExist(err) {
		t.Error("manifest.yaml should not be copied to OBS output root")
	}
	if _, err := os.Stat(filepath.Join(outDir, "media.yaml")); !os.IsNotExist(err) {
		t.Error("media.yaml should not be copied to OBS output root")
	}
	if _, ok := metadata.Ingredients["ingredients/LICENSE.md"]; !ok {
		t.Error("ingredients/LICENSE.md should exist in OBS metadata ingredients")
	}
}

// --- Registry tests ---

func TestLookup_AllRegisteredSubjects(t *testing.T) {
	expectedSubjects := []string{
		"Open Bible Stories",
		"Aligned Bible",
		"Bible",
		"Hebrew Old Testament",
		"Greek New Testament",
		"Translation Words",
		"Translation Academy",
		"TSV Translation Notes",
		"TSV Translation Questions",
		"TSV Translation Words Links",
		"TSV OBS Study Notes",
		"TSV OBS Study Questions",
		"TSV OBS Translation Notes",
		"TSV OBS Translation Questions",
	}

	for _, subject := range expectedSubjects {
		t.Run(subject, func(t *testing.T) {
			h, err := handler.Lookup(subject)
			if err != nil {
				t.Fatalf("Lookup(%q) failed: %v", subject, err)
			}
			if h.Subject() != subject {
				t.Errorf("Subject() = %q; want %q", h.Subject(), subject)
			}
		})
	}
}

func TestSupportedSubjects_Count(t *testing.T) {
	subjects := handler.SupportedSubjects()
	if len(subjects) != 14 {
		t.Errorf("SupportedSubjects() returned %d subjects; want 14. Got: %v", len(subjects), subjects)
	}
}

func TestLookup_UnsupportedSubject(t *testing.T) {
	_, err := handler.Lookup("Nonexistent Subject")
	if err == nil {
		t.Fatal("expected error for unsupported subject")
	}
	if !strings.Contains(err.Error(), "unsupported subject") {
		t.Errorf("error should mention 'unsupported subject': %v", err)
	}
}
