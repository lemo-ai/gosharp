# gosharp

Go 常用工具库，提供类型转换、时间处理、错误栈、并发编排、高性能 JSON、切片与哈希等能力。

> English: see [README.md](./README.md)

## 安装

```bash
go get github.com/lemo-ai/gosharp@latest
```

要求 **Go 1.26.0+**。

## 包一览

| 包 | 说明 |
|---|---|
| `convert` | 任意类型互转（数字、字符串、切片、map、struct、时间等） |
| `gtime` | 时间封装，支持常见字符串解析与自定义格式 |
| `gerror` | 带堆栈的错误，支持 `Wrap` / `Cause` / `errors.Is` |
| `mr` | MapReduce / ForEach / Finish 并发编排 |
| `json` | 高性能 JSON 编解码（类型 codec 缓存与常见类型特化路径） |
| `encoding/gbinary` | 大端 / 小端二进制编解码与位操作 |
| `regex` | 带编译缓存的正则 API |
| `empty` | 空值 / nil 判断 |
| `judge` / `stringutil` | 字符串判定与简单处理 |
| `structutil` | 结构体字段 / tag 反射工具 |
| `collection` | 泛型切片工具（Contains / Unique / Filter / Chunk 等） |
| `retry` | 可配置重试（次数 / 延迟 / 退避） |
| `safego` | 带 panic recover 的安全 goroutine |
| `hashx` | MD5 / SHA1 / SHA256 / FNV 快捷方法 |
| `randx` | 安全随机数 / 随机字符串 |

## 快速示例

### 类型转换

```go
import "github.com/lemo-ai/gosharp/convert"

n := convert.Int("42")
ss := convert.Strings([]string{"a", "b"})
```

### JSON

```go
import "github.com/lemo-ai/gosharp/json"

b, err := json.Marshal(map[string]any{"ok": true})
var out map[string]any
err = json.Unmarshal(b, &out)
_ = json.Pretouch(MyType{}) // 预热热点类型
```

### 带堆栈错误

```go
import (
	"errors"
	"io"

	"github.com/lemo-ai/gosharp/gerror"
)

err := gerror.Wrap(io.EOF, "read failed")
_ = errors.Is(err, io.EOF) // true
```

### 切片工具

```go
import "github.com/lemo-ai/gosharp/collection"

collection.Unique([]int{1, 2, 2, 3})     // [1 2 3]
collection.Chunk([]int{1, 2, 3, 4, 5}, 2) // [[1 2] [3 4] [5]]
```

## 说明

- `gtime.SetTimeZone` 只影响本包默认时区（`gtime.Location()`），不会修改进程级 `time.Local`。
- `stringutil` 委托给 `judge`，新代码推荐直接使用 `judge`。

## 许可证

[Apache License 2.0](./LICENSE)

Copyright 2024-2026 lemo-ai
