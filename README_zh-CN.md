# gosharp

Go 常用工具库：类型转换、时间、错误栈、MapReduce、**自研高性能 JSON**、切片工具等。

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
| `encoding/gbinary` | 二进制编解码 |
| `json` | **自研**高性能 JSON（类型 codec 缓存 + 特化路径，不依赖 sonic） |
| `regex` | 带编译缓存的正则 |
| `empty` | 空值 / nil 判断 |
| `judge` / `stringutil` | 字符串工具 |
| `structutil` | 结构体 tag 反射 |
| `collection` | 泛型切片工具 |
| `retry` | 重试 |
| `safego` | 安全 goroutine |
| `hashx` | 哈希快捷方法 |
| `randx` | 随机数 / 随机串 |

## JSON（自研）

`json` 包是独立实现，**不引用** `bytedance/sonic`。设计参考了高性能 JSON 引擎的常见思路：

- 按 `reflect.Type` 编译并缓存 codec
- `[]int` / `[]string` / `map[string]string` 等特化编解码
- buffer pool、无 HTML escape（默认）、unsafe 零拷贝写回

压测对比（Apple M5，含 sonic 仅作对照依赖）：

```bash
GOTOOLCHAIN=go1.26.0 go test ./json/ -bench=. -benchmem
```

典型结果（越高越好的相对速度，以 sonic 为 1.0）：

| 场景 | gosharp vs sonic |
|------|------------------|
| 小对象 Marshal | **约 2× 更快** |
| 小对象 Unmarshal | **持平或略快** |
| 大对象 Marshal | **约 1.6× 更快** |
| 大对象 Unmarshal | 仍略慢于 sonic（JIT/SIMD 优势） |

## 许可证

[Apache License 2.0](./LICENSE) · Copyright 2024-2026 lemo-ai
