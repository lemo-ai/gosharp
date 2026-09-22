# gosharp

Go 常用工具库，封装类型转换、时间处理、错误栈、MapReduce、正则缓存、二进制编解码等能力。

> English: see [README.md](./README.md)

## 安装

```bash
go get github.com/lemo-ai/gosharp@latest
```

要求 Go 1.21+。

## 包一览

| 包 | 说明 |
|---|---|
| `convert` | 任意类型互转（数字、字符串、切片、map、struct、时间等） |
| `gtime` | 时间封装，支持常见字符串解析与自定义格式 |
| `gerror` | 带堆栈的错误，支持 `Wrap` / `Cause` / `errors.Is` |
| `mr` | MapReduce / ForEach / Finish 并发编排 |
| `encoding/gbinary` | 大端/小端二进制编解码与位操作 |
| `json` | 基于 json-iterator 的 JSON 编解码 |
| `regex` | 带编译缓存的正则 API |
| `empty` | 空值 / nil 判断 |
| `judge` / `stringutil` | 字符串判定与简单处理 |
| `structutil` | 结构体字段 / tag 反射工具 |

## 快速示例

### 类型转换

```go
package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/convert"
)

func main() {
	fmt.Println(convert.Int("42")) // 42
	fmt.Println(convert.Strings([]string{"a", "b"}))
	fmt.Println(convert.Map(struct {
		Name string `json:"name"`
	}{Name: "gosharp"}))
}
```

### 带堆栈错误

```go
import (
	"errors"
	"fmt"
	"io"

	"github.com/lemo-ai/gosharp/gerror"
)

err := gerror.Wrap(io.EOF, "read failed")
fmt.Println(errors.Is(err, io.EOF)) // true
fmt.Printf("%+v\n", err)            // 错误信息 + 堆栈
```

### MapReduce

```go
import "github.com/lemo-ai/gosharp/mr"

v, err := mr.MapReduce(
	func(source chan<- interface{}) {
		for i := 1; i <= 5; i++ {
			source <- i
		}
	},
	func(item interface{}, writer mr.Writer, cancel func(error)) {
		writer.Write(item.(int) * 2)
	},
	func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		sum := 0
		for v := range pipe {
			sum += v.(int)
		}
		writer.Write(sum)
	},
)
```

## 许可证

[Apache License 2.0](./LICENSE)

Copyright 2024-2026 lemo-ai

部分实现参考了 [GoFrame](https://github.com/gogf/gf)、[go-zero](https://github.com/zeromicro/go-zero) 等开源项目的设计思路，并在本仓库中做了裁剪与修正。

## 说明

- `gtime.SetTimeZone` 只影响本包默认时区（见 `gtime.Location()`），**不会**修改进程级 `time.Local`。
- `stringutil` 已委托给 `judge`，两者 API 等价，推荐新代码直接使用 `judge`。
- 欢迎补充测试与 issue。
