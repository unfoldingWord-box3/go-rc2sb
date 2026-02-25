// Package books provides Bible book code data, localized names, and sort order.
package books

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unfoldingWord/go-rc2sb/sb"
)

// LocalizedBookNames holds localized book name data extracted from USFM markers.
type LocalizedBookNames struct {
	Long  string // from \toc1 (or fallback: \mt1, \mt)
	Short string // from \toc2 (or fallback: \h)
	Abbr  string // from \toc3
}

// BookInfo holds information about a single Bible book.
type BookInfo struct {
	ID    string // lowercase identifier (e.g., "gen")
	Code  string // uppercase USFM code (e.g., "GEN")
	Sort  int    // canonical sort order (1-66)
	Abbr  string // abbreviation (e.g., "Gen")
	Short string // short name (e.g., "Genesis")
	Long  string // long name (e.g., "The Book of Genesis")
}

// AllBooks is the ordered list of all 66 Bible books.
var AllBooks = []BookInfo{
	{ID: "gen", Code: "GEN", Sort: 1, Abbr: "Gen", Short: "Genesis", Long: "The Book of Genesis"},
	{ID: "exo", Code: "EXO", Sort: 2, Abbr: "Exo", Short: "Exodus", Long: "The Book of Exodus"},
	{ID: "lev", Code: "LEV", Sort: 3, Abbr: "Lev", Short: "Leviticus", Long: "The Book of Leviticus"},
	{ID: "num", Code: "NUM", Sort: 4, Abbr: "Num", Short: "Numbers", Long: "The Book of Numbers"},
	{ID: "deu", Code: "DEU", Sort: 5, Abbr: "Deu", Short: "Deuteronomy", Long: "The Book of Deuteronomy"},
	{ID: "jos", Code: "JOS", Sort: 6, Abbr: "Jos", Short: "Joshua", Long: "The Book of Joshua"},
	{ID: "jdg", Code: "JDG", Sort: 7, Abbr: "Jdg", Short: "Judges", Long: "The Book of Judges"},
	{ID: "rut", Code: "RUT", Sort: 8, Abbr: "Rut", Short: "Ruth", Long: "The Book of Ruth"},
	{ID: "1sa", Code: "1SA", Sort: 9, Abbr: "1Sa", Short: "First Samuel", Long: "The First Book of Samuel"},
	{ID: "2sa", Code: "2SA", Sort: 10, Abbr: "2Sa", Short: "Second Samuel", Long: "The Second Book of Samuel"},
	{ID: "1ki", Code: "1KI", Sort: 11, Abbr: "1Ki", Short: "First Kings", Long: "The First Book of Kings"},
	{ID: "2ki", Code: "2KI", Sort: 12, Abbr: "2Ki", Short: "Second Kings", Long: "The Second Book of Kings"},
	{ID: "1ch", Code: "1CH", Sort: 13, Abbr: "1Ch", Short: "First Chronicles", Long: "The First Book of the Chronicles"},
	{ID: "2ch", Code: "2CH", Sort: 14, Abbr: "2Ch", Short: "Second Chronicles", Long: "The Second Book of the Chronicles"},
	{ID: "ezr", Code: "EZR", Sort: 15, Abbr: "Ezr", Short: "Ezra", Long: "The Book of Ezra"},
	{ID: "neh", Code: "NEH", Sort: 16, Abbr: "Neh", Short: "Nehemiah", Long: "The Book of Nehemiah"},
	{ID: "est", Code: "EST", Sort: 17, Abbr: "Est", Short: "Esther", Long: "The Book of Esther"},
	{ID: "job", Code: "JOB", Sort: 18, Abbr: "Job", Short: "Job", Long: "The Book of Job"},
	{ID: "psa", Code: "PSA", Sort: 19, Abbr: "Psa", Short: "Psalms", Long: "The Book of Psalms"},
	{ID: "pro", Code: "PRO", Sort: 20, Abbr: "Pro", Short: "Proverbs", Long: "The Book of Proverbs"},
	{ID: "ecc", Code: "ECC", Sort: 21, Abbr: "Ecc", Short: "Ecclesiastes", Long: "The Book of Ecclesiastes"},
	{ID: "sng", Code: "SNG", Sort: 22, Abbr: "Sng", Short: "Song of Songs", Long: "The Song of Songs"},
	{ID: "isa", Code: "ISA", Sort: 23, Abbr: "Isa", Short: "Isaiah", Long: "The Book of Isaiah"},
	{ID: "jer", Code: "JER", Sort: 24, Abbr: "Jer", Short: "Jeremiah", Long: "The Book of Jeremiah"},
	{ID: "lam", Code: "LAM", Sort: 25, Abbr: "Lam", Short: "Lamentations", Long: "The Book of Lamentations"},
	{ID: "ezk", Code: "EZK", Sort: 26, Abbr: "Ezk", Short: "Ezekiel", Long: "The Book of Ezekiel"},
	{ID: "dan", Code: "DAN", Sort: 27, Abbr: "Dan", Short: "Daniel", Long: "The Book of Daniel"},
	{ID: "hos", Code: "HOS", Sort: 28, Abbr: "Hos", Short: "Hosea", Long: "The Book of Hosea"},
	{ID: "jol", Code: "JOL", Sort: 29, Abbr: "Jol", Short: "Joel", Long: "The Book of Joel"},
	{ID: "amo", Code: "AMO", Sort: 30, Abbr: "Amo", Short: "Amos", Long: "The Book of Amos"},
	{ID: "oba", Code: "OBA", Sort: 31, Abbr: "Oba", Short: "Obadiah", Long: "The Book of Obadiah"},
	{ID: "jon", Code: "JON", Sort: 32, Abbr: "Jon", Short: "Jonah", Long: "The Book of Jonah"},
	{ID: "mic", Code: "MIC", Sort: 33, Abbr: "Mic", Short: "Micah", Long: "The Book of Micah"},
	{ID: "nam", Code: "NAM", Sort: 34, Abbr: "Nam", Short: "Nahum", Long: "The Book of Nahum"},
	{ID: "hab", Code: "HAB", Sort: 35, Abbr: "Hab", Short: "Habakkuk", Long: "The Book of Habakkuk"},
	{ID: "zep", Code: "ZEP", Sort: 36, Abbr: "Zep", Short: "Zephaniah", Long: "The Book of Zephaniah"},
	{ID: "hag", Code: "HAG", Sort: 37, Abbr: "Hag", Short: "Haggai", Long: "The Book of Haggai"},
	{ID: "zec", Code: "ZEC", Sort: 38, Abbr: "Zec", Short: "Zechariah", Long: "The Book of Zechariah"},
	{ID: "mal", Code: "MAL", Sort: 39, Abbr: "Mal", Short: "Malachi", Long: "The Book of Malachi"},
	{ID: "mat", Code: "MAT", Sort: 40, Abbr: "Mat", Short: "Matthew", Long: "The Gospel of Matthew"},
	{ID: "mrk", Code: "MRK", Sort: 41, Abbr: "Mrk", Short: "Mark", Long: "The Gospel of Mark"},
	{ID: "luk", Code: "LUK", Sort: 42, Abbr: "Luk", Short: "Luke", Long: "The Gospel of Luke"},
	{ID: "jhn", Code: "JHN", Sort: 43, Abbr: "Jhn", Short: "John", Long: "The Gospel of John"},
	{ID: "act", Code: "ACT", Sort: 44, Abbr: "Act", Short: "Acts", Long: "The Acts of the Apostles"},
	{ID: "rom", Code: "ROM", Sort: 45, Abbr: "Rom", Short: "Romans", Long: "The Letter of Paul to the Romans"},
	{ID: "1co", Code: "1CO", Sort: 46, Abbr: "1Co", Short: "First Corinthians", Long: "The First Letter of Paul to the Corinthians"},
	{ID: "2co", Code: "2CO", Sort: 47, Abbr: "2Co", Short: "Second Corinthians", Long: "The Second Letter of Paul to the Corinthians"},
	{ID: "gal", Code: "GAL", Sort: 48, Abbr: "Gal", Short: "Galatians", Long: "The Letter of Paul to the Galatians"},
	{ID: "eph", Code: "EPH", Sort: 49, Abbr: "Eph", Short: "Ephesians", Long: "The Letter of Paul to the Ephesians"},
	{ID: "php", Code: "PHP", Sort: 50, Abbr: "Php", Short: "Philippians", Long: "The Letter of Paul to the Philippians"},
	{ID: "col", Code: "COL", Sort: 51, Abbr: "Col", Short: "Colossians", Long: "The Letter of Paul to the Colossians"},
	{ID: "1th", Code: "1TH", Sort: 52, Abbr: "1Th", Short: "First Thessalonians", Long: "The First Letter of Paul to the Thessalonians"},
	{ID: "2th", Code: "2TH", Sort: 53, Abbr: "2Th", Short: "Second Thessalonians", Long: "The Second Letter of Paul to the Thessalonians"},
	{ID: "1ti", Code: "1TI", Sort: 54, Abbr: "1Ti", Short: "First Timothy", Long: "The First Letter of Paul to Timothy"},
	{ID: "2ti", Code: "2TI", Sort: 55, Abbr: "2Ti", Short: "Second Timothy", Long: "The Second Letter of Paul to Timothy"},
	{ID: "tit", Code: "TIT", Sort: 56, Abbr: "Tit", Short: "Titus", Long: "The Letter of Paul to Titus"},
	{ID: "phm", Code: "PHM", Sort: 57, Abbr: "Phm", Short: "Philemon", Long: "The Letter of Paul to Philemon"},
	{ID: "heb", Code: "HEB", Sort: 58, Abbr: "Heb", Short: "Hebrews", Long: "The Letter to the Hebrews"},
	{ID: "jas", Code: "JAS", Sort: 59, Abbr: "Jas", Short: "James", Long: "The Letter of James"},
	{ID: "1pe", Code: "1PE", Sort: 60, Abbr: "1Pe", Short: "First Peter", Long: "The First Letter of Peter"},
	{ID: "2pe", Code: "2PE", Sort: 61, Abbr: "2Pe", Short: "Second Peter", Long: "The Second Letter of Peter"},
	{ID: "1jn", Code: "1JN", Sort: 62, Abbr: "1Jn", Short: "First John", Long: "The First Letter of John"},
	{ID: "2jn", Code: "2JN", Sort: 63, Abbr: "2Jn", Short: "Second John", Long: "The Second Letter of John"},
	{ID: "3jn", Code: "3JN", Sort: 64, Abbr: "3Jn", Short: "Third John", Long: "The Third Letter of John"},
	{ID: "jud", Code: "JUD", Sort: 65, Abbr: "Jud", Short: "Jude", Long: "The Letter of Jude"},
	{ID: "rev", Code: "REV", Sort: 66, Abbr: "Rev", Short: "Revelation", Long: "The Book of Revelation"},
}

// bookByID is a lookup map from lowercase identifier to BookInfo.
var bookByID map[string]*BookInfo

// bookByCode is a lookup map from uppercase code to BookInfo.
var bookByCode map[string]*BookInfo

func init() {
	bookByID = make(map[string]*BookInfo, len(AllBooks))
	bookByCode = make(map[string]*BookInfo, len(AllBooks))
	for i := range AllBooks {
		b := &AllBooks[i]
		bookByID[b.ID] = b
		bookByCode[b.Code] = b
	}
}

// ByID returns the BookInfo for a lowercase identifier (e.g., "gen"), or nil if not found.
func ByID(id string) *BookInfo {
	return bookByID[strings.ToLower(id)]
}

// ByCode returns the BookInfo for an uppercase code (e.g., "GEN"), or nil if not found.
func ByCode(code string) *BookInfo {
	return bookByCode[strings.ToUpper(code)]
}

// IsBookID returns true if the given identifier is a recognized Bible book.
func IsBookID(id string) bool {
	return bookByID[strings.ToLower(id)] != nil
}

// LocalizedNameEntry returns the SB LocalizedName for a book identifier.
func LocalizedNameEntry(id string) (string, sb.LocalizedName) {
	b := ByID(id)
	if b == nil {
		return "", sb.LocalizedName{}
	}
	key := "book-" + b.ID
	return key, sb.LocalizedName{
		Abbr:  map[string]string{"en": b.Abbr},
		Short: map[string]string{"en": b.Short},
		Long:  map[string]string{"en": b.Long},
	}
}

// CodeFromProjectID returns the uppercase book code for a project identifier.
// For example, "gen" -> "GEN".
func CodeFromProjectID(id string) string {
	b := ByID(id)
	if b != nil {
		return b.Code
	}
	// Fallback: uppercase the id
	return strings.ToUpper(id)
}

// LocalizedNameEntryWithNames builds a LocalizedName using the following priority:
//  1. USFM toc markers (if usfmNames is non-nil and fields are non-empty)
//  2. Manifest project title (if non-empty)
//  3. English fallback from AllBooks
//
// The lang parameter specifies the language tag for the localized names (e.g., "hi", "en").
// English fallback names are always included under the "en" key.
func LocalizedNameEntryWithNames(id string, lang string, projectTitle string, usfmNames *LocalizedBookNames) (string, sb.LocalizedName) {
	b := ByID(id)
	if b == nil {
		return "", sb.LocalizedName{}
	}
	key := "book-" + b.ID

	ln := sb.LocalizedName{
		Abbr:  make(map[string]string),
		Short: make(map[string]string),
		Long:  make(map[string]string),
	}

	// Always include English fallback names
	ln.Abbr["en"] = b.Abbr
	ln.Short["en"] = b.Short
	ln.Long["en"] = b.Long

	// If the language is English, we only need to apply overrides from USFM/manifest
	// on top of the English defaults. If non-English, add localized entries under the lang key.
	if lang == "en" {
		// For English, USFM toc values override the hardcoded English defaults
		if usfmNames != nil {
			if usfmNames.Long != "" {
				ln.Long["en"] = usfmNames.Long
			}
			if usfmNames.Short != "" {
				ln.Short["en"] = usfmNames.Short
			}
			if usfmNames.Abbr != "" {
				ln.Abbr["en"] = usfmNames.Abbr
			}
		}
		return key, ln
	}

	// For non-English languages, build localized entries under the lang key.
	// Priority: USFM toc markers > project title > (omit localized key)

	// Long name: \toc1 > projectTitle > (English only)
	localLong := ""
	if usfmNames != nil && usfmNames.Long != "" {
		localLong = usfmNames.Long
	} else if projectTitle != "" {
		localLong = projectTitle
	}
	if localLong != "" {
		ln.Long[lang] = localLong
	}

	// Short name: \toc2 > projectTitle > (English only)
	localShort := ""
	if usfmNames != nil && usfmNames.Short != "" {
		localShort = usfmNames.Short
	} else if projectTitle != "" {
		localShort = projectTitle
	}
	if localShort != "" {
		ln.Short[lang] = localShort
	}

	// Abbreviation: \toc3 > (English only)
	if usfmNames != nil && usfmNames.Abbr != "" {
		ln.Abbr[lang] = usfmNames.Abbr
	}

	return key, ln
}

// ParseUSFMBookNames reads the first 20 lines of a USFM file and extracts
// \toc1, \toc2, \toc3 markers for localized book names. Falls back to \mt1/\mt
// for the long name and \h for the short name if toc markers are missing.
// Returns nil if the file doesn't exist or contains no useful markers.
func ParseUSFMBookNames(filePath string) *LocalizedBookNames {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var toc1, toc2, toc3, h, mt string

	scanner := bufio.NewScanner(f)
	lineCount := 0
	for scanner.Scan() && lineCount < 20 {
		lineCount++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if val := extractUSFMMarker(line, `\toc1`); val != "" {
			toc1 = val
		} else if val := extractUSFMMarker(line, `\toc2`); val != "" {
			toc2 = val
		} else if val := extractUSFMMarker(line, `\toc3`); val != "" {
			toc3 = val
		} else if val := extractUSFMMarker(line, `\h`); val != "" {
			h = val
		} else if val := extractUSFMMarker(line, `\mt1`); val != "" {
			mt = val
		} else if val := extractUSFMMarker(line, `\mt`); val != "" && mt == "" {
			// \mt without number, only use if \mt1 wasn't found
			mt = val
		}
	}

	// Build result with fallbacks
	longName := toc1
	if longName == "" {
		longName = mt
	}
	shortName := toc2
	if shortName == "" {
		shortName = h
	}

	// Return nil if nothing useful was found
	if longName == "" && shortName == "" && toc3 == "" {
		return nil
	}

	return &LocalizedBookNames{
		Long:  longName,
		Short: shortName,
		Abbr:  toc3,
	}
}

// FindUSFMFile searches for a USFM file matching a book code in a directory.
// It looks for patterns like "NN-CODE.usfm" (e.g., "01-GEN.usfm") or "CODE.usfm".
// Returns the full path if found, or empty string if not found.
func FindUSFMFile(usfmDir string, bookID string) string {
	code := CodeFromProjectID(bookID)

	// Try NN-CODE.usfm pattern first (most common)
	matches, err := filepath.Glob(filepath.Join(usfmDir, fmt.Sprintf("*-%s.usfm", code)))
	if err == nil && len(matches) > 0 {
		return matches[0]
	}

	// Try CODE.usfm
	direct := filepath.Join(usfmDir, code+".usfm")
	if _, err := os.Stat(direct); err == nil {
		return direct
	}

	// Try lowercase variants
	matches, err = filepath.Glob(filepath.Join(usfmDir, fmt.Sprintf("*-%s.usfm", strings.ToLower(code))))
	if err == nil && len(matches) > 0 {
		return matches[0]
	}

	return ""
}

// extractUSFMMarker extracts the value after a USFM marker like "\toc1 VALUE".
// Returns empty string if the line doesn't start with the marker.
func extractUSFMMarker(line, marker string) string {
	// The marker must be at the start of the line and followed by a space
	if !strings.HasPrefix(line, marker+" ") && line != marker {
		return ""
	}
	val := strings.TrimPrefix(line, marker)
	return strings.TrimSpace(val)
}
