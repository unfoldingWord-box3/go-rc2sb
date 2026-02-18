package rc2sb

// Options configures the RC to SB conversion.
type Options struct {
	// PayloadDirs maps subject names to paths of additional RC repos
	// that provide payload data.
	// Example: {"Translation Words": "/path/to/en_tw"} for TWL conversion.
	PayloadDirs map[string]string
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
