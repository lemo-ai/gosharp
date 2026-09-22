# gosharp

A Go utility library for type conversion, time helpers, stacked errors, concurrency helpers, high-performance JSON, collections, hashing, and more.

中文说明见 [README_zh-CN.md](./README_zh-CN.md)。

## Install

```bash
go get github.com/lemo-ai/gosharp@latest
```

Requires **Go 1.26.0+**.

## Packages

| Package | Purpose |
|---|---|
| `convert` | Convert between common Go types |
| `gtime` | Time wrapper with flexible parsing/formatting |
| `gerror` | Errors with stack traces (`Wrap` / `Cause` / `errors.Is`) |
| `mr` | MapReduce / ForEach / Finish concurrency helpers |
| `json` | High-performance JSON encode/decode |
| `encoding/gbinary` | Little/big-endian binary encode/decode |
| `regex` | Regex APIs with compiled-pattern cache |
| `empty` | Empty / nil checks |
| `judge` / `stringutil` | String checks and light transforms |
| `structutil` | Struct field / tag reflection helpers |
| `collection` | Generic slice helpers |
| `retry` | Retry with delay / backoff |
| `safego` | Goroutines with panic recovery |
| `hashx` | MD5 / SHA / FNV helpers |
| `randx` | Crypto & fast random helpers |

## Quick examples

```go
import (
	"github.com/lemo-ai/gosharp/collection"
	"github.com/lemo-ai/gosharp/convert"
	"github.com/lemo-ai/gosharp/json"
)

_ = convert.Int("42")
_, _ = json.Marshal(map[string]any{"ok": true})
_ = collection.Unique([]int{1, 2, 2, 3})
```

## Notes

- `gtime.SetTimeZone` only affects this package (`gtime.Location()`), not process-wide `time.Local`.
- `stringutil` delegates to `judge`; prefer `judge` in new code.

## License

[Apache License 2.0](./LICENSE)

Copyright 2024-2026 lemo-ai
