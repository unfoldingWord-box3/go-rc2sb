# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go library (`github.com/nichmahn/go-rc2sb`) for converting Resource Container (RC) repositories to Scripture Burrito (SB) repositories. RC is a format used by unfoldingWord for Bible translation resources (spec: rc0.2). SB (Scripture Burrito) is the newer standardized format (spec: 1.0.0).

## Build & Test Commands

```bash
go build ./...           # Build all packages
go test ./...            # Run all tests
go test ./... -v         # Run tests with verbose output
go test -run TestName    # Run a specific test
go vet ./...             # Lint/static analysis

# Run only unit tests (no samples needed)
go test ./rc/... ./sb/... ./books/...

# Run integration tests (requires samples/ directory)
go test -run TestConvert -v
```

## Architecture

### Public API

The single entry point is `Convert()` in `convert.go`:

```go
func Convert(ctx context.Context, inDir string, outDir string, opts Options) (Result, error)
```

- Reads `manifest.yaml` from `inDir`, determines the subject, looks up the handler, runs conversion, writes `metadata.json` to `outDir`.
- `Options.PayloadDirs` maps subject names to auxiliary RC repo paths (e.g., TWL needs Translation Words).

### Package Structure

```
go-rc2sb/
├── convert.go              # Public Convert() function, orchestration
├── options.go              # Options and Result types
├── rc/
│   └── manifest.go         # RC manifest.yaml parsing (DublinCore, projects)
├── sb/
│   ├── metadata.go         # SB metadata.json types and JSON serialization
│   └── ingredient.go       # Ingredient computation (MD5, MIME type, size)
├── books/
│   └── books.go            # Bible book data (66 books, localized names, codes)
├── handler/
│   ├── handler.go          # Handler interface definition
│   ├── registry.go         # Subject -> handler registry (Register/Lookup)
│   ├── common.go           # Shared helpers (file copy, metadata building, copyright)
│   ├── obs.go              # Open Bible Stories handler
│   ├── aligned_bible.go    # Bible/USFM handler (Aligned Bible, Bible, Hebrew OT, Greek NT)
│   ├── tw.go               # Translation Words handler
│   ├── ta.go               # Translation Academy handler
│   ├── tn.go               # TSV Translation Notes handler
│   ├── tq.go               # TSV Translation Questions handler
│   ├── twl.go              # TSV Translation Words Links handler (with payload)
│   ├── obs_tsv.go          # Generic OBS TSV handler (4 variants)
│   └── subjects/
│       └── register.go     # Registers all 14 handlers via init()
```

### Key Design Patterns

- **Handler pattern**: Each subject type implements the `Handler` interface (`Subject() string`, `Convert(...)`). Handlers are registered in `handler/subjects/register.go` via `init()`.
- **Blank import for registration**: `convert.go` imports `_ "github.com/nichmahn/go-rc2sb/handler/subjects"` to trigger handler registration.
- **Shared helpers in `handler/common.go`**: `BuildBaseMetadata()`, `BuildCopyright()`, `CopyFileAndComputeIngredient()`, `CopyFileWithScope()`, `CopyLicenseIngredient()`, `copyTreeToIngredients()`.

### Subject -> SB Type Mapping

| Subject | FlavorType/Flavor | IdAuthority | Abbreviation |
|---------|-------------------|-------------|-------------|
| Open Bible Stories | gloss/textStories | BurritoTruck | OBS |
| Aligned Bible | scripture/textTranslation | uWBurritos | (from RC identifier) |
| Bible | scripture/textTranslation | uWBurritos | (from RC identifier) |
| Hebrew Old Testament | scripture/textTranslation | uWBurritos | (from RC identifier) |
| Greek New Testament | scripture/textTranslation | uWBurritos | (from RC identifier) |
| Translation Words | peripheral/x-peripheralArticles | uWBurritos | TW |
| Translation Academy | peripheral/x-peripheralArticles | uWBurritos | TA |
| TSV Translation Notes | parascriptural/x-bcvnotes | uWBurritos | TN |
| TSV Translation Questions | parascriptural/x-bcvquestions | uWBurritos | TQ |
| TSV Translation Words Links | parascriptural/x-bcvarticles | uWBurritos | TW |
| TSV OBS Study Notes | peripheral/x-obsnotes | BurritoTruck | OBSSN |
| TSV OBS Study Questions | peripheral/x-obsquestions | BurritoTruck | OBSSQ |
| TSV OBS Translation Notes | peripheral/x-obsnotes | BurritoTruck | OBSTN |
| TSV OBS Translation Questions | peripheral/x-obsquestions | BurritoTruck | OBSTQ |

### RC Format (Input)
- **manifest.yaml**: Dublin Core metadata (conformsto: rc0.2), project list, language, versioning
- **media.yaml**: Optional media format definitions (PDF, audio, video URLs)
- **content/**: Resource files in various formats depending on type

### SB Format (Output)
- **metadata.json**: Scripture Burrito metadata with identification, languages, type/flavor, and an `ingredients` map listing every file with its MD5 checksum, MIME type, and size
- **ingredients/**: All content files organized under this directory

### Key Conversion Logic

1. **Metadata**: Transform `manifest.yaml` (Dublin Core) into `metadata.json` (Scripture Burrito schema) — map identifiers, versions, languages, project info
2. **File relocation**: Copy content files into `ingredients/` directory, adjusting paths per resource type (e.g., strip `tn_` prefix from TSV filenames, strip numeric prefix from USFM filenames)
3. **Checksum computation**: SB metadata.json requires MD5 checksums, MIME types, and byte sizes for every ingredient file
4. **Content preservation**: File contents (Markdown, USFM, TSV) are unchanged between formats
5. **Payload resolution**: TWL extracts `rc://*/tw/dict/bible/{category}/{article}` links from the TWLink TSV column and copies matched TW articles to `ingredients/payload/`

### Testing

- **Integration tests** (`convert_test.go`): One test per subject type (11 total). Requires `samples/` directory (gitignored) with RC/SB pairs. Tests compare structural metadata (flavor type, scope keys, abbreviation, language, ingredient keys) and verify internal consistency (every ingredient exists on disk with correct MD5 and size).
- **Unit tests**: `rc/manifest_test.go`, `sb/ingredient_test.go`, `sb/metadata_test.go`, `books/books_test.go`
- **Error handling tests** (`error_test.go`): Missing manifest, unsupported subject, cancelled context, invalid YAML

### Dependencies

- `gopkg.in/yaml.v3` — YAML parsing for RC manifest files
