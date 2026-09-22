# gosharp

Go 常用工具库，提供类型转换、时间处理、错误栈、并发编排、高性能 JSON、切片工具、常用算法等能力。

> English: see [README.md](./README.md)

## 安装

```bash
go get github.com/lemo-ai/gosharp@latest
```

要求 **Go 1.26.0+**。

## 包一览

| 包 | 说明 |
|---|---|
| `convert` | 任意类型互转 |
| `gtime` | 时间封装 |
| `gerror` | 带堆栈错误 |
| `mr` | MapReduce / ForEach / Finish |
| `json` | 高性能 JSON 编解码 |
| `encoding/gbinary` | 二进制编解码 |
| `regex` | 带缓存的正则 |
| `empty` | 空值 / nil 判断 |
| `judge` / `stringutil` | 字符串工具 |
| `structutil` | 结构体 tag 反射 |
| `collection` | 泛型切片工具 |
| `retry` | 重试 |
| `safego` | 安全 goroutine |
| `hashx` | 哈希快捷方法 |
| `randx` | 随机数 / 随机串 |
| `algo` | 常用算法（数学 / 搜索排序 / 字符串 / 并查集 / LRU） |

## 示例

完整可运行示例见 [`examples/`](./examples/)：

```bash
go run ./examples/convert
go run ./examples/json
go run ./examples/algo
go run ./examples/mr
# 其余见 examples/README.md
```

## 快速上手

```go
import (
	"github.com/lemo-ai/gosharp/algo"
	"github.com/lemo-ai/gosharp/collection"
	"github.com/lemo-ai/gosharp/convert"
	"github.com/lemo-ai/gosharp/json"
)

_ = convert.Int("42")
_, _ = json.Marshal(map[string]any{"ok": true})
_ = collection.Unique([]int{1, 2, 2, 3})
_ = algo.GCD(54, 24)
```

## 说明

- `gtime.SetTimeZone` 只影响本包默认时区（`gtime.Location()`），不会修改 `time.Local`。
- `stringutil` 委托给 `judge`，新代码推荐直接使用 `judge`。

## 许可证

[Apache License 2.0](./LICENSE)

Copyright 2024-2026 lemo-ai
