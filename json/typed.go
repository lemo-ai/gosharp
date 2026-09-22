package json

import (
	"reflect"
)

func compileCodec(t reflect.Type) *codec {
	c := &codec{}
	switch t.Kind() {
	case reflect.Bool:
		c.encode = encodeBool
		c.decode = decodeBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		c.encode = encodeInt
		c.decode = decodeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		c.encode = encodeUint
		c.decode = decodeUint
	case reflect.Float32, reflect.Float64:
		c.encode = encodeFloat
		c.decode = decodeFloat
	case reflect.String:
		c.encode = encodeString
		c.decode = decodeString
	case reflect.Interface:
		c.encode = encodeInterface
		c.decode = decodeInterface
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			c.encode = encodeBytes
			c.decode = decodeBytes
		case reflect.Int:
			c.encode = encodeIntSlice
			c.decode = decodeIntSlice
		case reflect.Int64:
			c.encode = encodeInt64Slice
			c.decode = decodeInt64Slice
		case reflect.String:
			c.encode = encodeStringSlice
			c.decode = decodeStringSlice
		default:
			elem := getCodec(t.Elem())
			c.encode = makeEncodeSlice(elem)
			c.decode = makeDecodeSlice(elem)
		}
	case reflect.Array:
		elem := getCodec(t.Elem())
		c.encode = makeEncodeArray(elem)
		c.decode = makeDecodeArray(elem)
	case reflect.Map:
		if t.Key().Kind() == reflect.String && t.Elem().Kind() == reflect.String {
			c.encode = encodeMapStringString
			c.decode = decodeMapStringString
		} else {
			c.encode = makeEncodeMap(t)
			c.decode = makeDecodeMap(t)
		}
	case reflect.Struct:
		sc := compileStruct(t)
		c.encode = sc.encode
		c.decode = sc.decode
	case reflect.Pointer:
		elem := getCodec(t.Elem())
		c.encode = makeEncodePtr(elem)
		c.decode = makeDecodePtr(elem)
	default:
		c.encode = encodeUnsupported
		c.decode = decodeUnsupported
	}
	return c
}

func encodeIntSlice(e *encoder, v reflect.Value) error {
	if v.IsNil() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	s := sliceIntOf(v)
	e.buf = append(e.buf, '[')
	for i, n := range s {
		if i > 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = appendInt(e.buf, int64(n))
	}
	e.buf = append(e.buf, ']')
	return nil
}

// decodeIntSlice is a single-pass, allocation-light decoder for []int.
func decodeIntSlice(d *decoder, v reflect.Value) error {
	if err := d.beginArrayOrNull(); err != nil {
		if err == errWasNull {
			setSliceInt(v, nil)
			return nil
		}
		return err
	}
	data := d.data
	off := d.off
	if data[off] == ']' {
		d.off = off + 1
		setSliceInt(v, make([]int, 0))
		return nil
	}
	out := make([]int, 0, 16)
	for {
		neg := false
		if off < len(data) && data[off] == '-' {
			neg = true
			off++
		}
		if off >= len(data) || data[off] < '0' || data[off] > '9' {
			return syntaxError(off, "invalid number")
		}
		var n int64
		for off < len(data) {
			c := data[off]
			if c < '0' || c > '9' {
				break
			}
			n = n*10 + int64(c-'0')
			off++
		}
		if neg {
			n = -n
		}
		out = append(out, int(n))
		if off >= len(data) {
			return syntaxError(off, "unterminated array")
		}
		c := data[off]
		off++
		if c == ']' {
			break
		}
		if c != ',' {
			return syntaxError(off, "expect ',' or ']'")
		}
	}
	d.off = off
	setSliceInt(v, out)
	return nil
}

func encodeInt64Slice(e *encoder, v reflect.Value) error {
	if v.IsNil() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	s := sliceInt64Of(v)
	e.buf = append(e.buf, '[')
	for i, n := range s {
		if i > 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = appendInt(e.buf, n)
	}
	e.buf = append(e.buf, ']')
	return nil
}

func decodeInt64Slice(d *decoder, v reflect.Value) error {
	if err := d.beginArrayOrNull(); err != nil {
		if err == errWasNull {
			setSliceInt64(v, nil)
			return nil
		}
		return err
	}
	data := d.data
	off := d.off
	if data[off] == ']' {
		d.off = off + 1
		setSliceInt64(v, make([]int64, 0))
		return nil
	}
	out := make([]int64, 0, 16)
	for {
		neg := false
		if off < len(data) && data[off] == '-' {
			neg = true
			off++
		}
		if off >= len(data) || data[off] < '0' || data[off] > '9' {
			return syntaxError(off, "invalid number")
		}
		var n int64
		for off < len(data) {
			c := data[off]
			if c < '0' || c > '9' {
				break
			}
			n = n*10 + int64(c-'0')
			off++
		}
		if neg {
			n = -n
		}
		out = append(out, n)
		if off >= len(data) {
			return syntaxError(off, "unterminated array")
		}
		c := data[off]
		off++
		if c == ']' {
			break
		}
		if c != ',' {
			return syntaxError(off, "expect ',' or ']'")
		}
	}
	d.off = off
	setSliceInt64(v, out)
	return nil
}

func encodeStringSlice(e *encoder, v reflect.Value) error {
	if v.IsNil() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	s := sliceStringOf(v)
	e.buf = append(e.buf, '[')
	for i, str := range s {
		if i > 0 {
			e.buf = append(e.buf, ',')
		}
		e.writeString(str)
	}
	e.buf = append(e.buf, ']')
	return nil
}

func decodeStringSlice(d *decoder, v reflect.Value) error {
	if err := d.beginArrayOrNull(); err != nil {
		if err == errWasNull {
			setSliceString(v, nil)
			return nil
		}
		return err
	}
	if d.data[d.off] == ']' {
		d.off++
		setSliceString(v, make([]string, 0))
		return nil
	}
	out := make([]string, 0, 8)
	for {
		s, err := d.readStringFast()
		if err != nil {
			return err
		}
		out = append(out, s)
		if d.off >= len(d.data) {
			return syntaxError(d.off, "unterminated array")
		}
		c := d.data[d.off]
		d.off++
		if c == ']' {
			break
		}
		if c != ',' {
			return syntaxError(d.off, "expect ',' or ']'")
		}
	}
	setSliceString(v, out)
	return nil
}

func encodeMapStringString(e *encoder, v reflect.Value) error {
	if v.IsNil() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	m := mapStringStringOf(v)
	e.buf = append(e.buf, '{')
	first := true
	for k, val := range m {
		if !first {
			e.buf = append(e.buf, ',')
		}
		first = false
		e.writeString(k)
		e.buf = append(e.buf, ':')
		e.writeString(val)
	}
	e.buf = append(e.buf, '}')
	return nil
}

func decodeMapStringString(d *decoder, v reflect.Value) error {
	d.skipSpace()
	if d.off < len(d.data) && d.data[d.off] == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		setMapStringString(v, nil)
		return nil
	}
	if err := d.expect('{'); err != nil {
		return err
	}
	d.skipSpace()
	if d.off < len(d.data) && d.data[d.off] == '}' {
		d.off++
		setMapStringString(v, make(map[string]string))
		return nil
	}
	m := make(map[string]string, 8)
	for {
		k, err := d.readStringFast()
		if err != nil {
			return err
		}
		d.skipSpace()
		if d.off >= len(d.data) || d.data[d.off] != ':' {
			return syntaxError(d.off, "expect ':'")
		}
		d.off++
		val, err := d.readStringFast()
		if err != nil {
			return err
		}
		m[k] = val
		d.skipSpace()
		if d.off >= len(d.data) {
			return syntaxError(d.off, "unterminated object")
		}
		c := d.data[d.off]
		d.off++
		if c == '}' {
			break
		}
		if c != ',' {
			return syntaxError(d.off, "expect ',' or '}'")
		}
	}
	setMapStringString(v, m)
	return nil
}

func sliceIntOf(v reflect.Value) []int {
	if v.CanAddr() {
		return *(*[]int)(v.Addr().UnsafePointer())
	}
	return v.Interface().([]int)
}

func sliceInt64Of(v reflect.Value) []int64 {
	if v.CanAddr() {
		return *(*[]int64)(v.Addr().UnsafePointer())
	}
	return v.Interface().([]int64)
}

func sliceStringOf(v reflect.Value) []string {
	if v.CanAddr() {
		return *(*[]string)(v.Addr().UnsafePointer())
	}
	return v.Interface().([]string)
}

func mapStringStringOf(v reflect.Value) map[string]string {
	if v.CanAddr() {
		return *(*map[string]string)(v.Addr().UnsafePointer())
	}
	return v.Interface().(map[string]string)
}

func setSliceInt(v reflect.Value, out []int) {
	if v.CanAddr() {
		*(*[]int)(v.Addr().UnsafePointer()) = out
		return
	}
	v.Set(reflect.ValueOf(out))
}

func setSliceInt64(v reflect.Value, out []int64) {
	if v.CanAddr() {
		*(*[]int64)(v.Addr().UnsafePointer()) = out
		return
	}
	v.Set(reflect.ValueOf(out))
}

func setSliceString(v reflect.Value, out []string) {
	if v.CanAddr() {
		*(*[]string)(v.Addr().UnsafePointer()) = out
		return
	}
	v.Set(reflect.ValueOf(out))
}

func setMapStringString(v reflect.Value, out map[string]string) {
	if v.CanAddr() {
		*(*map[string]string)(v.Addr().UnsafePointer()) = out
		return
	}
	v.Set(reflect.ValueOf(out))
}
