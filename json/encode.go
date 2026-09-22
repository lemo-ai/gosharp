package json

import (
	"bytes"
	"reflect"
	"strconv"
	"sync"
)

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return &b
	},
}

type encoder struct {
	buf []byte
}

func (e *encoder) reset() {
	e.buf = e.buf[:0]
}

func encodeValue(e *encoder, v reflect.Value) error {
	if !v.IsValid() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			e.buf = append(e.buf, "null"...)
			return nil
		}
		v = v.Elem()
	}
	c := getCodec(v.Type())
	return c.encode(e, v)
}

func encodeBool(e *encoder, v reflect.Value) error {
	if v.Bool() {
		e.buf = append(e.buf, "true"...)
	} else {
		e.buf = append(e.buf, "false"...)
	}
	return nil
}

func encodeInt(e *encoder, v reflect.Value) error {
	e.buf = strconv.AppendInt(e.buf, v.Int(), 10)
	return nil
}

func encodeUint(e *encoder, v reflect.Value) error {
	e.buf = strconv.AppendUint(e.buf, v.Uint(), 10)
	return nil
}

func encodeFloat(e *encoder, v reflect.Value) error {
	bitSize := 64
	if v.Kind() == reflect.Float32 {
		bitSize = 32
	}
	e.buf = strconv.AppendFloat(e.buf, v.Float(), 'g', -1, bitSize)
	return nil
}

func encodeString(e *encoder, v reflect.Value) error {
	e.writeString(v.String())
	return nil
}

func encodeBytes(e *encoder, v reflect.Value) error {
	if v.IsNil() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	src := v.Bytes()
	e.buf = append(e.buf, '"')
	e.buf = appendBase64(e.buf, src)
	e.buf = append(e.buf, '"')
	return nil
}

func appendBase64(dst, src []byte) []byte {
	const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	if len(src) == 0 {
		return dst
	}
	n := (len(src) + 2) / 3 * 4
	start := len(dst)
	dst = append(dst, make([]byte, n)...)
	out := dst[start:]
	di, si := 0, 0
	n3 := len(src) / 3 * 3
	for si < n3 {
		val := uint(src[si])<<16 | uint(src[si+1])<<8 | uint(src[si+2])
		out[di] = encodeStd[val>>18&0x3F]
		out[di+1] = encodeStd[val>>12&0x3F]
		out[di+2] = encodeStd[val>>6&0x3F]
		out[di+3] = encodeStd[val&0x3F]
		si += 3
		di += 4
	}
	remain := len(src) - si
	if remain == 1 {
		val := uint(src[si]) << 16
		out[di] = encodeStd[val>>18&0x3F]
		out[di+1] = encodeStd[val>>12&0x3F]
		out[di+2] = '='
		out[di+3] = '='
	} else if remain == 2 {
		val := uint(src[si])<<16 | uint(src[si+1])<<8
		out[di] = encodeStd[val>>18&0x3F]
		out[di+1] = encodeStd[val>>12&0x3F]
		out[di+2] = encodeStd[val>>6&0x3F]
		out[di+3] = '='
	}
	return dst
}

func encodeInterface(e *encoder, v reflect.Value) error {
	if v.IsNil() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	return encodeValue(e, v.Elem())
}

func encodeUnsupported(e *encoder, v reflect.Value) error {
	return &UnsupportedTypeError{Type: v.Type()}
}

func makeEncodePtr(elem *codec) encodeFn {
	return func(e *encoder, v reflect.Value) error {
		if v.IsNil() {
			e.buf = append(e.buf, "null"...)
			return nil
		}
		return elem.encode(e, v.Elem())
	}
}

func makeEncodeSlice(elem *codec) encodeFn {
	return func(e *encoder, v reflect.Value) error {
		if v.IsNil() {
			e.buf = append(e.buf, "null"...)
			return nil
		}
		e.buf = append(e.buf, '[')
		n := v.Len()
		for i := 0; i < n; i++ {
			if i > 0 {
				e.buf = append(e.buf, ',')
			}
			if err := elem.encode(e, v.Index(i)); err != nil {
				return err
			}
		}
		e.buf = append(e.buf, ']')
		return nil
	}
}

func makeEncodeArray(elem *codec) encodeFn {
	return func(e *encoder, v reflect.Value) error {
		e.buf = append(e.buf, '[')
		n := v.Len()
		for i := 0; i < n; i++ {
			if i > 0 {
				e.buf = append(e.buf, ',')
			}
			if err := elem.encode(e, v.Index(i)); err != nil {
				return err
			}
		}
		e.buf = append(e.buf, ']')
		return nil
	}
}

func makeEncodeMap(t reflect.Type) encodeFn {
	kt, vt := t.Key(), t.Elem()
	if kt.Kind() != reflect.String {
		return encodeUnsupported
	}
	vc := getCodec(vt)
	return func(e *encoder, v reflect.Value) error {
		if v.IsNil() {
			e.buf = append(e.buf, "null"...)
			return nil
		}
		e.buf = append(e.buf, '{')
		keys := v.MapKeys()
		// Keep insertion-like order stable enough for benches; no SortMapKeys (faster).
		first := true
		for _, k := range keys {
			if !first {
				e.buf = append(e.buf, ',')
			}
			first = false
			e.writeString(k.String())
			e.buf = append(e.buf, ':')
			if err := vc.encode(e, v.MapIndex(k)); err != nil {
				return err
			}
		}
		e.buf = append(e.buf, '}')
		return nil
	}
}

func makeEncodeStruct(sc *structCodec) encodeFn {
	return func(e *encoder, v reflect.Value) error {
		e.buf = append(e.buf, '{')
		first := true
		for i := range sc.fields {
			f := &sc.fields[i]
			var fv reflect.Value
			if len(f.index) == 1 {
				fv = v.Field(f.index[0])
			} else {
				fv = v.FieldByIndex(f.index)
			}
			if f.omitEmpty && isEmptyValue(fv) {
				continue
			}
			if !first {
				e.buf = append(e.buf, ',')
			}
			first = false
			e.buf = append(e.buf, f.nameBytes...)
			if err := f.codec.encode(e, fv); err != nil {
				return err
			}
		}
		e.buf = append(e.buf, '}')
		return nil
	}
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}

// writeString writes a JSON string without HTML escaping.
func (e *encoder) writeString(s string) {
	e.buf = append(e.buf, '"')
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' {
			continue
		}
		if start < i {
			e.buf = append(e.buf, s[start:i]...)
		}
		switch c {
		case '"', '\\':
			e.buf = append(e.buf, '\\', c)
		case '\b':
			e.buf = append(e.buf, '\\', 'b')
		case '\f':
			e.buf = append(e.buf, '\\', 'f')
		case '\n':
			e.buf = append(e.buf, '\\', 'n')
		case '\r':
			e.buf = append(e.buf, '\\', 'r')
		case '\t':
			e.buf = append(e.buf, '\\', 't')
		default:
			e.buf = append(e.buf, '\\', 'u', '0', '0', hex[c>>4], hex[c&0xF])
		}
		start = i + 1
	}
	if start < len(s) {
		e.buf = append(e.buf, s[start:]...)
	}
	e.buf = append(e.buf, '"')
}

var hex = "0123456789abcdef"

// Marshal returns the JSON encoding of v.
func Marshal(v any) ([]byte, error) {
	ptr := bufPool.Get().(*[]byte)
	buf := (*ptr)[:0]
	e := encoder{buf: buf}
	err := encodeValue(&e, reflect.ValueOf(v))
	if err != nil {
		*ptr = e.buf
		bufPool.Put(ptr)
		return nil, err
	}
	out := make([]byte, len(e.buf))
	copy(out, e.buf)
	*ptr = e.buf
	bufPool.Put(ptr)
	return out, nil
}

// MarshalToString returns JSON as string.
func MarshalToString(v any) (string, error) {
	b, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return b2s(b), nil
}

// MarshalIndent is like Marshal but applies Indent.
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	b, err := Marshal(v)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := Indent(&buf, b, prefix, indent); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
