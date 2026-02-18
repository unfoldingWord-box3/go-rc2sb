package books_test

import (
	"testing"

	"github.com/nichmahn/go-rc2sb/books"
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
