# gosharp

A Go utility library for conversion, time, errors, concurrency, high-performance JSON, collections, crypto helpers, and common algorithms.

中文说明见 [README_zh-CN.md](./README_zh-CN.md)。

## Install

```bash
go get github.com/lemo-ai/gosharp@latest
```

Requires **Go 1.26.0+**.

## Packages

| Package | Purpose |
|---|---|
| `convert` / `gtime` / `gerror` / `mr` | conversion, time, stacked errors, MapReduce |
| `json` | high-performance JSON |
| `collection` / `retry` / `safego` / `hashx` / `randx` | slices, retry, safe goroutines, hash, random |
| `cryptox` | AES-GCM, ChaCha20-Poly1305, RSA-OAEP/PSS, HMAC, bcrypt, PBKDF2 |
| `encoding/gbinary` / `regex` / `empty` / `judge` / `structutil` | binary, regex, empty checks, strings, struct tags |
| `algo` | algorithms: math, search/sort, string, UnionFind, LRU |

## Examples

Runnable demos live under [`examples/`](./examples/):

```bash
go run ./examples/convert
go run ./examples/json
go run ./examples/algo
go run ./examples/cryptox
```

See [examples/README.md](./examples/README.md) for the full list.

## License

[Apache License 2.0](./LICENSE)

Copyright 2024-2026 lemo-ai
