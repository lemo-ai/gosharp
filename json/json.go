// Package json is a self-contained high-performance JSON library for gosharp.
//
// Design (inspired by common high-perf JSON engines, implemented independently):
//   - compile-once type codecs cached by reflect.Type
//   - buffer pooling for marshal
//   - no HTML-escape by default (faster; closer to sonic ConfigDefault)
//   - specialized paths for primitives / slices / maps / structs
//
// It does NOT wrap bytedance/sonic. Benchmark against sonic in json_test.go.
package json

import (
	"bytes"
	encodingjson "encoding/json"
	"reflect"
)

// Pretouch compiles codecs for v's type tree ahead of time.
func Pretouch(v any) error {
	if v == nil {
		return nil
	}
	t := reflect.TypeOf(v)
	pretouchType(t)
	return nil
}

// Indent appends to dst an indented form of the JSON encoding src.
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	return encodingjson.Indent(dst, src, prefix, indent)
}
