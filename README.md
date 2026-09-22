# gosharp

A small Go utility library for type conversion, time helpers, stacked errors, MapReduce, regex caching, and binary encoding.

中文说明见 [README_zh-CN.md](./README_zh-CN.md)。

## Install

```bash
go get github.com/lemo-ai/gosharp@latest
```

Requires Go 1.21+.

## Packages

| Package | Purpose |
|---|---|
| `convert` | Convert between common Go types |
| `gtime` | Time wrapper with flexible parsing/formatting |
| `gerror` | Errors with stack traces (`Wrap` / `Cause` / `errors.Is`) |
| `mr` | MapReduce / ForEach / Finish concurrency helpers |
| `encoding/gbinary` | Little/big-endian binary encode/decode |
| `json` | JSON helpers on top of json-iterator |
| `regex` | Regex APIs with compiled-pattern cache |
| `empty` | Empty / nil checks |
| `judge` / `stringutil` | String checks and light transforms |
| `structutil` | Struct field / tag reflection helpers |

## License

[Apache License 2.0](./LICENSE)

Copyright 2024-2026 lemo-ai

Some designs were inspired by [GoFrame](https://github.com/gogf/gf) and [go-zero](https://github.com/zeromicro/go-zero), adapted and fixed in this repository.

## Notes

- `gtime.SetTimeZone` only affects this package's default location (`gtime.Location()`); it does **not** mutate process-wide `time.Local`.
- `stringutil` delegates to `judge`; prefer `judge` in new code.
