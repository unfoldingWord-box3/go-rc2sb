// Command rc2sb converts a Resource Container (RC) repository to Scripture Burrito (SB) format.
//
// Usage:
//
//	rc2sb [flags] <inDir> <outDir>
//	rc2sb --payload /path/to/en_tw <inDir> <outDir>
//	rc2sb --usfm /path/to/en_ult <inDir> <outDir>
//
// Flags:
//
//	--payload <dir>   Path to a Translation Words directory (e.g., en_tw) for TWL payload creation.
//	                  If not set, auto-detects <lang>_tw/ inside inDir.
//	--usfm <dir>      Path to a USFM directory for localized Bible book names in TSV repos.
//	                  If not set, uses manifest project titles, then English fallback.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	rc2sb "github.com/unfoldingWord/go-rc2sb"
)

func main() {
	payload := flag.String("payload", "", "path to a Translation Words directory (e.g., en_tw) for TWL payload creation")
	usfm := flag.String("usfm", "", "path to a USFM directory for localized Bible book names in TSV repos")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: rc2sb [flags] <inDir> <outDir>\n\n")
		fmt.Fprintf(os.Stderr, "Converts a Resource Container (RC) repository to Scripture Burrito (SB) format.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  inDir    Path to the RC repository (must contain manifest.yaml)\n")
		fmt.Fprintf(os.Stderr, "  outDir   Path where SB output will be written\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	inDir := flag.Arg(0)
	outDir := flag.Arg(1)

	opts := rc2sb.Options{
		PayloadPath: *payload,
		USFMPath:    *usfm,
	}

	result, err := rc2sb.Convert(context.Background(), inDir, outDir, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Converted %s (%s) with %d ingredients\n",
		result.Subject, result.Identifier, result.Ingredients)
}
