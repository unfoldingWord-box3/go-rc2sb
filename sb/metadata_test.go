package sb_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/nichmahn/go-rc2sb/sb"
)

func TestNewMetadata(t *testing.T) {
	m := sb.NewMetadata()
	if m.Format != "scripture burrito" {
		t.Errorf("Format = %q; want %q", m.Format, "scripture burrito")
	}
	if m.Meta.Version != "1.0.0" {
		t.Errorf("Meta.Version = %q; want %q", m.Meta.Version, "1.0.0")
	}
	if m.Confidential != false {
		t.Error("Confidential should be false")
	}
	if m.Ingredients == nil {
		t.Error("Ingredients should not be nil")
	}
}

func TestMetadata_WriteToFile(t *testing.T) {
	m := sb.NewMetadata()
	m.Format = "scripture burrito"
	m.Identification = sb.Identification{
		Name: map[string]string{"en": "Test"},
	}

	dir := t.TempDir()
	if err := m.WriteToFile(dir); err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}

	// Read back
	data, err := os.ReadFile(filepath.Join(dir, "metadata.json"))
	if err != nil {
		t.Fatalf("reading metadata.json: %v", err)
	}

	var m2 sb.Metadata
	if err := json.Unmarshal(data, &m2); err != nil {
		t.Fatalf("parsing metadata.json: %v", err)
	}

	if m2.Format != "scripture burrito" {
		t.Errorf("Format = %q; want %q", m2.Format, "scripture burrito")
	}
	if m2.Identification.Name["en"] != "Test" {
		t.Errorf("Name = %q; want %q", m2.Identification.Name["en"], "Test")
	}
}

func TestMetadata_JSONRoundTrip(t *testing.T) {
	m := sb.NewMetadata()
	m.Type = sb.Type{
		FlavorType: sb.FlavorType{
			Name: "parascriptural",
			Flavor: sb.Flavor{
				Name: "x-bcvnotes",
			},
			CurrentScope: map[string][]string{
				"GEN": {},
				"EXO": {},
			},
		},
	}
	m.Ingredients["test.tsv"] = sb.Ingredient{
		Checksum: sb.Checksum{MD5: "abc123"},
		MimeType: "text/tab-separated-values",
		Size:     1234,
		Scope:    map[string][]string{"GEN": {}},
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var m2 sb.Metadata
	if err := json.Unmarshal(data, &m2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if m2.Type.FlavorType.Name != "parascriptural" {
		t.Errorf("FlavorType.Name = %q; want %q", m2.Type.FlavorType.Name, "parascriptural")
	}
	if len(m2.Type.FlavorType.CurrentScope) != 2 {
		t.Errorf("CurrentScope len = %d; want 2", len(m2.Type.FlavorType.CurrentScope))
	}
	if ing, ok := m2.Ingredients["test.tsv"]; !ok {
		t.Error("missing ingredient test.tsv")
	} else if ing.Size != 1234 {
		t.Errorf("ingredient size = %d; want 1234", ing.Size)
	}
}
