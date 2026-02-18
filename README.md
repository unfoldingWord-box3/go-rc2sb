# go-rc2sb

Go library for converting [Resource Container](https://resource-container.readthedocs.io/) (RC) repositories to [Scripture Burrito](https://docs.burrito.bible/) (SB) format.

## Installation

```bash
go get github.com/nichmahn/go-rc2sb
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    rc2sb "github.com/nichmahn/go-rc2sb"
)

func main() {
    ctx := context.Background()

    result, err := rc2sb.Convert(ctx, "/path/to/rc-repo", "/path/to/sb-output", rc2sb.Options{})
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Converted %s (%s) with %d ingredients\n",
        result.Subject, result.Identifier, result.Ingredients)
}
```

### With Payload (TWL)

For subjects that require a payload directory (e.g., TSV Translation Words Links needs the Translation Words repo):

```go
opts := rc2sb.Options{
    PayloadDirs: map[string]string{
        "Translation Words": "/path/to/en_tw",
    },
}

result, err := rc2sb.Convert(ctx, "/path/to/en_twl", "/path/to/output", opts)
```

## API

### `Convert(ctx, inDir, outDir, opts) (Result, error)`

Converts an RC repository to SB format.

- `ctx` - Context for cancellation
- `inDir` - Path to the RC repository (must contain `manifest.yaml`)
- `outDir` - Path where SB output will be written
- `opts` - Options including optional payload directories

Returns a `Result` with conversion metadata, or an error.

### Options

```go
type Options struct {
    PayloadDirs map[string]string // Maps subject names to payload RC repo paths
}
```

### Result

```go
type Result struct {
    Subject     string // RC subject that was converted
    Identifier  string // RC identifier (e.g., "obs", "ult", "tn")
    InDir       string // Input RC directory
    OutDir      string // Output SB directory
    Ingredients int    // Number of ingredient files
}
```

## Supported Subjects

| Subject | SB Flavor Type | Notes |
|---------|---------------|-------|
| Open Bible Stories | gloss/textStories | Copies content/ to ingredients/content/ |
| Aligned Bible | scripture/textTranslation | Strips numeric prefix from USFM filenames; abbreviation from RC identifier |
| Bible | scripture/textTranslation | Same as Aligned Bible (e.g., ULT, UST) |
| Hebrew Old Testament | scripture/textTranslation | Same as Aligned Bible (e.g., UHB) |
| Greek New Testament | scripture/textTranslation | Same as Aligned Bible (e.g., UGNT) |
| Translation Words | peripheral/x-peripheralArticles | Copies bible/{kt,other,names}/ articles |
| Translation Academy | peripheral/x-peripheralArticles | Copies nested markdown hierarchy |
| TSV Translation Notes | parascriptural/x-bcvnotes | Strips tn_ prefix from TSV filenames |
| TSV Translation Questions | parascriptural/x-bcvquestions | Strips tq_ prefix from TSV filenames |
| TSV Translation Words Links | parascriptural/x-bcvarticles | Includes TW payload when provided |
| TSV OBS Study Notes | peripheral/x-obsnotes | Single TSV file conversion |
| TSV OBS Study Questions | peripheral/x-obsquestions | Single TSV file conversion |
| TSV OBS Translation Notes | peripheral/x-obsnotes | Single TSV file conversion |
| TSV OBS Translation Questions | peripheral/x-obsquestions | Single TSV file conversion |

## Error Handling

- Missing `manifest.yaml` returns an error indicating the directory is not a valid RC repo
- Unsupported subjects return an error listing all supported subjects
- Context cancellation is checked at key points during conversion
- File I/O errors are wrapped with context and returned

## Building

```bash
go build ./...
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run a specific test
go test -run TestConvertOpenBibleStories

# Run only unit tests (no samples needed)
go test ./rc/... ./sb/... ./books/...

# Run integration tests (requires samples/ directory)
go test -run TestConvert -v
```

### Integration Tests

Integration tests use sample RC/SB pairs in the `samples/` directory (gitignored). Each test:

1. Runs `Convert()` on the sample RC input
2. Verifies the output metadata structure matches the expected SB metadata
3. Verifies internal consistency (every ingredient in metadata.json exists on disk with correct MD5 and size)

### Unit Tests

- `rc/manifest_test.go` - Manifest parsing (valid, invalid, missing)
- `sb/ingredient_test.go` - MD5/MIME/size computation
- `sb/metadata_test.go` - Metadata creation, serialization, round-trip
- `books/books_test.go` - Book lookups, localized names, sort order
- `error_test.go` - Error handling (missing manifest, unsupported subject, cancelled context)

## Architecture

```
go-rc2sb/
+-- convert.go              # Public Convert() function
+-- options.go              # Options and Result types
+-- rc/
|   +-- manifest.go         # RC manifest.yaml parsing
+-- sb/
|   +-- metadata.go         # SB metadata.json types
|   +-- ingredient.go       # Ingredient computation (MD5, MIME, size)
+-- books/
|   +-- books.go            # Bible book data (66 books, localized names)
+-- handler/
|   +-- handler.go          # Handler interface
|   +-- registry.go         # Subject -> handler registry
|   +-- common.go           # Shared helpers (file copy, metadata building)
|   +-- obs.go              # Open Bible Stories
|   +-- aligned_bible.go    # Bible/USFM handler (Aligned Bible, Bible, Hebrew OT, Greek NT)
|   +-- tw.go               # Translation Words
|   +-- ta.go               # Translation Academy
|   +-- tn.go               # TSV Translation Notes
|   +-- tq.go               # TSV Translation Questions
|   +-- twl.go              # TSV Translation Words Links (with payload)
|   +-- obs_tsv.go          # OBS TSV variants (4 types)
|   +-- subjects/
|       +-- register.go     # Registers all handlers
```

## License

See [LICENSE](LICENSE) for details.
