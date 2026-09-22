# gosharp

Go utilities: conversion, time, errors, MapReduce, **self-contained high-performance JSON**, collections, and more.

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
| `json` | **in-house** JSON (codec cache + typed paths; not a sonic wrapper) |
| `collection` / `retry` / `safego` / `hashx` / `randx` | slices, retry, safe goroutines, hash, random |
| `encoding/gbinary` / `regex` / `empty` / `judge` / `structutil` | binary, regex, empty checks, strings, struct tags |

## JSON

Implemented in-tree (inspired by common high-perf techniques, **not** wrapping ByteDance sonic). Sonic is only used in benchmarks for comparison.

```bash
GOTOOLCHAIN=go1.26.0 go test ./json/ -bench=. -benchmem
```

On Apple M5: Marshal beats sonic (~1.6–2×); small Unmarshal matches/beats sonic; large Unmarshal may still trail sonic’s JIT/SIMD.

## License

[Apache License 2.0](./LICENSE) · Copyright 2024-2026 lemo-ai
