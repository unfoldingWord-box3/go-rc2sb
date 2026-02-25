package books_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/unfoldingWord/go-rc2sb/books"
)

func TestByID(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"gen", "GEN"},
		{"GEN", "GEN"},
		{"rev", "REV"},
		{"1co", "1CO"},
		{"psa", "PSA"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			b := books.ByID(tt.id)
			if b == nil {
				t.Fatalf("ByID(%q) returned nil", tt.id)
			}
			if b.Code != tt.want {
				t.Errorf("ByID(%q).Code = %q; want %q", tt.id, b.Code, tt.want)
			}
		})
	}
}

func TestByID_NotFound(t *testing.T) {
	b := books.ByID("xyz")
	if b != nil {
		t.Error("ByID('xyz') should return nil for unknown book")
	}
}

func TestByCode(t *testing.T) {
	b := books.ByCode("GEN")
	if b == nil {
		t.Fatal("ByCode('GEN') returned nil")
	}
	if b.ID != "gen" {
		t.Errorf("ByCode('GEN').ID = %q; want %q", b.ID, "gen")
	}
}

func TestIsBookID(t *testing.T) {
	if !books.IsBookID("gen") {
		t.Error("IsBookID('gen') should be true")
	}
	if books.IsBookID("frt") {
		t.Error("IsBookID('frt') should be false")
	}
	if books.IsBookID("obs") {
		t.Error("IsBookID('obs') should be false")
	}
}

func TestLocalizedNameEntry(t *testing.T) {
	key, ln := books.LocalizedNameEntry("gen")
	if key != "book-gen" {
		t.Errorf("key = %q; want %q", key, "book-gen")
	}
	if ln.Abbr["en"] != "Gen" {
		t.Errorf("Abbr = %q; want %q", ln.Abbr["en"], "Gen")
	}
	if ln.Short["en"] != "Genesis" {
		t.Errorf("Short = %q; want %q", ln.Short["en"], "Genesis")
	}
}

func TestLocalizedNameEntry_NotFound(t *testing.T) {
	key, _ := books.LocalizedNameEntry("xyz")
	if key != "" {
		t.Errorf("key should be empty for unknown book; got %q", key)
	}
}

func TestCodeFromProjectID(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"gen", "GEN"},
		{"1co", "1CO"},
		{"rev", "REV"},
		{"frt", "FRT"}, // not a book, should uppercase
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := books.CodeFromProjectID(tt.id)
			if got != tt.want {
				t.Errorf("CodeFromProjectID(%q) = %q; want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestAllBooks_Count(t *testing.T) {
	if len(books.AllBooks) != 66 {
		t.Errorf("AllBooks count = %d; want 66", len(books.AllBooks))
	}
}

func TestAllBooks_SortOrder(t *testing.T) {
	for i, b := range books.AllBooks {
		if b.Sort != i+1 {
			t.Errorf("AllBooks[%d].Sort = %d; want %d (book: %s)", i, b.Sort, i+1, b.ID)
		}
	}
}

// --- ParseUSFMBookNames tests ---

func TestParseUSFMBookNames_WithTocMarkers(t *testing.T) {
	dir := t.TempDir()
	usfmPath := filepath.Join(dir, "01-GEN.usfm")
	content := `\id GEN EN_ULT en_English_ltr
\usfm 3.0
\ide UTF-8
\h Genesis
\toc1 The Book of Genesis
\toc2 Genesis
\toc3 Gen
\mt Genesis

\ts\*
`
	os.WriteFile(usfmPath, []byte(content), 0644)

	names := books.ParseUSFMBookNames(usfmPath)
	if names == nil {
		t.Fatal("ParseUSFMBookNames returned nil")
	}
	if names.Long != "The Book of Genesis" {
		t.Errorf("Long = %q; want %q", names.Long, "The Book of Genesis")
	}
	if names.Short != "Genesis" {
		t.Errorf("Short = %q; want %q", names.Short, "Genesis")
	}
	if names.Abbr != "Gen" {
		t.Errorf("Abbr = %q; want %q", names.Abbr, "Gen")
	}
}

func TestParseUSFMBookNames_Hindi(t *testing.T) {
	dir := t.TempDir()
	usfmPath := filepath.Join(dir, "01-GEN.usfm")
	content := "\\id GEN EN_IRV hi_Hindi_ltr\n\\usfm 3.0\n\\ide UTF-8\n\\h \u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f\n\\toc1 \u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f\n\\toc2 \u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f\n\\toc3 gen\n\\mt1 \u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f\n"
	os.WriteFile(usfmPath, []byte(content), 0644)

	names := books.ParseUSFMBookNames(usfmPath)
	if names == nil {
		t.Fatal("ParseUSFMBookNames returned nil")
	}
	if names.Long != "\u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f" {
		t.Errorf("Long = %q; want Hindi Genesis", names.Long)
	}
	if names.Short != "\u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f" {
		t.Errorf("Short = %q; want Hindi Genesis", names.Short)
	}
	if names.Abbr != "gen" {
		t.Errorf("Abbr = %q; want %q", names.Abbr, "gen")
	}
}

func TestParseUSFMBookNames_FallbackToMtAndH(t *testing.T) {
	dir := t.TempDir()
	usfmPath := filepath.Join(dir, "01-GEN.usfm")
	// No \toc markers, only \h and \mt1
	content := `\id GEN
\usfm 3.0
\h Short Name
\mt1 Long Title
`
	os.WriteFile(usfmPath, []byte(content), 0644)

	names := books.ParseUSFMBookNames(usfmPath)
	if names == nil {
		t.Fatal("ParseUSFMBookNames returned nil")
	}
	if names.Long != "Long Title" {
		t.Errorf("Long = %q; want %q (should fall back to \\mt1)", names.Long, "Long Title")
	}
	if names.Short != "Short Name" {
		t.Errorf("Short = %q; want %q (should fall back to \\h)", names.Short, "Short Name")
	}
	if names.Abbr != "" {
		t.Errorf("Abbr = %q; want empty (no \\toc3)", names.Abbr)
	}
}

func TestParseUSFMBookNames_FallbackMtWithoutNumber(t *testing.T) {
	dir := t.TempDir()
	usfmPath := filepath.Join(dir, "01-GEN.usfm")
	content := `\id GEN
\mt Some Title
`
	os.WriteFile(usfmPath, []byte(content), 0644)

	names := books.ParseUSFMBookNames(usfmPath)
	if names == nil {
		t.Fatal("ParseUSFMBookNames returned nil")
	}
	if names.Long != "Some Title" {
		t.Errorf("Long = %q; want %q (should fall back to \\mt)", names.Long, "Some Title")
	}
}

func TestParseUSFMBookNames_MissingFile(t *testing.T) {
	names := books.ParseUSFMBookNames("/nonexistent/path.usfm")
	if names != nil {
		t.Error("ParseUSFMBookNames should return nil for missing file")
	}
}

func TestParseUSFMBookNames_NoUsefulMarkers(t *testing.T) {
	dir := t.TempDir()
	usfmPath := filepath.Join(dir, "01-GEN.usfm")
	content := `\id GEN
\usfm 3.0
\c 1
\v 1 In the beginning...
`
	os.WriteFile(usfmPath, []byte(content), 0644)

	names := books.ParseUSFMBookNames(usfmPath)
	if names != nil {
		t.Error("ParseUSFMBookNames should return nil when no useful markers found")
	}
}

// --- LocalizedNameEntryWithNames tests ---

func TestLocalizedNameEntryWithNames_EnglishWithUSFM(t *testing.T) {
	usfmNames := &books.LocalizedBookNames{
		Long:  "The Book of Genesis",
		Short: "Genesis",
		Abbr:  "Gen",
	}
	key, ln := books.LocalizedNameEntryWithNames("gen", "en", "Genesis Title", usfmNames)
	if key != "book-gen" {
		t.Errorf("key = %q; want %q", key, "book-gen")
	}
	// USFM should override defaults for English
	if ln.Long["en"] != "The Book of Genesis" {
		t.Errorf("Long[en] = %q; want %q", ln.Long["en"], "The Book of Genesis")
	}
	if ln.Short["en"] != "Genesis" {
		t.Errorf("Short[en] = %q; want %q", ln.Short["en"], "Genesis")
	}
	if ln.Abbr["en"] != "Gen" {
		t.Errorf("Abbr[en] = %q; want %q", ln.Abbr["en"], "Gen")
	}
}

func TestLocalizedNameEntryWithNames_NonEnglishWithUSFM(t *testing.T) {
	usfmNames := &books.LocalizedBookNames{
		Long:  "\u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f",
		Short: "\u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f",
		Abbr:  "gen",
	}
	key, ln := books.LocalizedNameEntryWithNames("gen", "hi", "", usfmNames)
	if key != "book-gen" {
		t.Errorf("key = %q; want %q", key, "book-gen")
	}
	// Should have both English and Hindi entries
	if ln.Long["en"] != "The Book of Genesis" {
		t.Errorf("Long[en] = %q; want English fallback", ln.Long["en"])
	}
	if ln.Long["hi"] != "\u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f" {
		t.Errorf("Long[hi] = %q; want Hindi name", ln.Long["hi"])
	}
	if ln.Abbr["hi"] != "gen" {
		t.Errorf("Abbr[hi] = %q; want %q", ln.Abbr["hi"], "gen")
	}
}

func TestLocalizedNameEntryWithNames_ProjectTitleFallback(t *testing.T) {
	// No USFM names, but project title is provided
	key, ln := books.LocalizedNameEntryWithNames("gen", "hi", "\u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f", nil)
	if key != "book-gen" {
		t.Errorf("key = %q; want %q", key, "book-gen")
	}
	// Long and Short should use the project title for non-English
	if ln.Long["hi"] != "\u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f" {
		t.Errorf("Long[hi] = %q; want project title", ln.Long["hi"])
	}
	if ln.Short["hi"] != "\u0909\u0924\u094d\u092a\u0924\u094d\u0924\u093f" {
		t.Errorf("Short[hi] = %q; want project title", ln.Short["hi"])
	}
	// Abbr should only have English (no USFM toc3)
	if _, ok := ln.Abbr["hi"]; ok {
		t.Error("Abbr[hi] should not exist without USFM toc3")
	}
	if ln.Abbr["en"] != "Gen" {
		t.Errorf("Abbr[en] = %q; want English fallback", ln.Abbr["en"])
	}
}

func TestLocalizedNameEntryWithNames_EnglishFallbackOnly(t *testing.T) {
	// No USFM, no project title
	key, ln := books.LocalizedNameEntryWithNames("gen", "hi", "", nil)
	if key != "book-gen" {
		t.Errorf("key = %q; want %q", key, "book-gen")
	}
	// Should only have English entries
	if ln.Long["en"] != "The Book of Genesis" {
		t.Errorf("Long[en] = %q; want English fallback", ln.Long["en"])
	}
	if _, ok := ln.Long["hi"]; ok {
		t.Error("Long[hi] should not exist when no localized data provided")
	}
}

func TestLocalizedNameEntryWithNames_UnknownBook(t *testing.T) {
	key, _ := books.LocalizedNameEntryWithNames("xyz", "en", "Some Title", nil)
	if key != "" {
		t.Errorf("key should be empty for unknown book; got %q", key)
	}
}

func TestLocalizedNameEntryWithNames_USFMOverridesProjectTitle(t *testing.T) {
	// Both USFM and project title provided — USFM should win
	usfmNames := &books.LocalizedBookNames{
		Long:  "USFM Long Name",
		Short: "USFM Short",
	}
	_, ln := books.LocalizedNameEntryWithNames("gen", "fr", "Manifest Title", usfmNames)
	if ln.Long["fr"] != "USFM Long Name" {
		t.Errorf("Long[fr] = %q; want USFM value over manifest title", ln.Long["fr"])
	}
	if ln.Short["fr"] != "USFM Short" {
		t.Errorf("Short[fr] = %q; want USFM value over manifest title", ln.Short["fr"])
	}
}

// --- FindUSFMFile tests ---

func TestFindUSFMFile_StandardPattern(t *testing.T) {
	dir := t.TempDir()
	// Create a standard NN-CODE.usfm file
	usfmPath := filepath.Join(dir, "01-GEN.usfm")
	os.WriteFile(usfmPath, []byte("\\id GEN\n"), 0644)

	found := books.FindUSFMFile(dir, "gen")
	if found != usfmPath {
		t.Errorf("FindUSFMFile = %q; want %q", found, usfmPath)
	}
}

func TestFindUSFMFile_DirectPattern(t *testing.T) {
	dir := t.TempDir()
	// Create a CODE.usfm file (no numeric prefix)
	usfmPath := filepath.Join(dir, "GEN.usfm")
	os.WriteFile(usfmPath, []byte("\\id GEN\n"), 0644)

	found := books.FindUSFMFile(dir, "gen")
	if found != usfmPath {
		t.Errorf("FindUSFMFile = %q; want %q", found, usfmPath)
	}
}

func TestFindUSFMFile_NotFound(t *testing.T) {
	dir := t.TempDir()
	found := books.FindUSFMFile(dir, "gen")
	if found != "" {
		t.Errorf("FindUSFMFile should return empty string when file not found; got %q", found)
	}
}
