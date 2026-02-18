package rc2sb

// Options configures the RC to SB conversion.
type Options struct {
	// PayloadPath is the path to a Translation Words directory (e.g., "/path/to/en_tw")
	// used when converting TSV Translation Words Links repos.
	// If set, the bible/ subdirectory within this path is copied to ingredients/payload/
	// in the SB output, and rc:// links in the TWL TSV files are rewritten to
	// relative ./payload/ paths.
	//
	// If empty, the TWL handler auto-detects a <lang>_tw/ subdirectory inside
	// the input RC repo directory (where <lang> is the manifest's language identifier).
	// If neither is found, no payload is created and TSV files are copied as-is.
	PayloadPath string
}

// Result holds information about a completed conversion.
type Result struct {
	// Subject is the RC subject that was converted.
	Subject string

	// Identifier is the RC identifier (e.g., "obs", "ult", "tn").
	Identifier string

	// InDir is the input RC directory that was converted.
	InDir string

	// OutDir is the output SB directory that was created.
	OutDir string

	// Ingredients is the number of ingredient files in the SB output.
	Ingredients int
}
