# saferio

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

`saferio` is a small Go library for safely handling temporary files and logs in [monitoring-forge](https://github.com/monitoring-forge) Mackerel plugins.

It provides helpers to read and write JSON files while reducing the risk of symlink-based attacks and partial writes.

## Features

- **Atomic JSON writes**: `WriteJSON` writes to a temporary file and renames it into place, minimizing the chance of leaving a partially written file.
- **Safe file opening**: `OpenRD` / `OpenAD` open files with `O_NOFOLLOW` on Unix-like systems to help prevent symlink attacks.
- **Cross-platform**: Separate implementations for Unix and Windows.
- **Small surface**: Only three public functions plus platform-specific open helpers.

## Usage

```go
package main

import (
	"log"

	"github.com/monitoring-forge/saferio"
)

func main() {
	dir := "/var/lib/mackerel-plugin"

	// Write metrics state atomically.
	state := map[string]any{"counter": 42}
	if err := saferio.WriteJSON(dir, "mackerel-plugin-foo.json", state); err != nil {
		log.Fatal(err)
	}

	// Read it back safely.
	var loaded map[string]any
	if err := saferio.ReadJSON(dir, "mackerel-plugin-foo.json", &loaded); err != nil {
		log.Fatal(err)
	}

	// Check for a file presence.
	if saferio.FileExists(dir, "mackerel-plugin-foo.json") {
		log.Println("state file exists")
	}
}
```

## API

### `FileExists(dir, filename string) bool`

Returns `true` if the file exists in `dir`.

### `WriteJSON(dir, filename string, v any) error`

Encodes `v` as JSON to a temporary file in `dir` and renames it to `filename`.
If `filename` already exists, it attempts to replace it atomically.

### `ReadJSON(dir, filename string, v any) error`

Opens `filename` safely (with `O_NOFOLLOW` on Unix) and decodes its JSON content into `v`.

If the file is empty, `ReadJSON` returns `io.EOF`. Callers that want to treat an empty file as "no state" can check for this error explicitly.

### `OpenRD(filename string) (*os.File, error)`

Opens a file for reading with `O_RDONLY`. On Unix-like systems, `O_NOFOLLOW` is added.

### `OpenAD(filename string) (*os.File, error)`

Opens a file for appending with `O_WRONLY | O_CREATE | O_APPEND`. On Unix-like systems, `O_NOFOLLOW` is added.

## Platform-specific behavior

| Function | Unix-like | Windows |
|----------|-----------|---------|
| `OpenRD` | `O_RDONLY \| O_NOFOLLOW` | `O_RDONLY` |
| `OpenAD` | `O_WRONLY \| O_CREATE \| O_APPEND \| O_NOFOLLOW` | `O_WRONLY \| O_CREATE \| O_APPEND` |

Windows does not support `O_NOFOLLOW`, so the flag is omitted on that platform.

## License

MIT License - see [LICENSE](LICENSE) for details.
