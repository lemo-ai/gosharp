package json

import (
	"reflect"
)

// Specialized codecs for common slice/map types — these dominate unmarshal allocs.

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
	var s []int
	if v.CanAddr() {
		s = *(*[]int)(v.Addr().UnsafePointer())
	} else {
		s = v.Interface().([]int)
	}
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

func decodeIntSlice(d *decoder, v reflect.Value) error {
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		setSliceInt(v, nil)
		return nil
	}
	if err := d.expect('['); err != nil {
		return err
	}
	out := make([]int, 0, estimateArrayLen(d, 8))
	if d.peek() == ']' {
		d.off++
		setSliceInt(v, out)
		return nil
	}
	for {
		n, err := d.readNumber()
		if err != nil {
			return err
		}
		i, ok := parseIntFast(n)
		if !ok {
			f, err := parseFloatBytes(n)
			if err != nil {
				return err
			}
			i = int64(f)
		}
		out = append(out, int(i))
		c := d.peek()
		if c == ']' {
			d.off++
			break
		}
		if err := d.expect(','); err != nil {
			return err
		}
	}
	setSliceInt(v, out)
	return nil
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

func encodeInt64Slice(e *encoder, v reflect.Value) error {
	if v.IsNil() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	var s []int64
	if v.CanAddr() {
		s = *(*[]int64)(v.Addr().UnsafePointer())
	} else {
		s = v.Interface().([]int64)
	}
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
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		setSliceInt64(v, nil)
		return nil
	}
	if err := d.expect('['); err != nil {
		return err
	}
	out := make([]int64, 0, estimateArrayLen(d, 8))
	if d.peek() == ']' {
		d.off++
		setSliceInt64(v, out)
		return nil
	}
	for {
		n, err := d.readNumber()
		if err != nil {
			return err
		}
		i, ok := parseIntFast(n)
		if !ok {
			f, err := parseFloatBytes(n)
			if err != nil {
				return err
			}
			i = int64(f)
		}
		out = append(out, i)
		c := d.peek()
		if c == ']' {
			d.off++
			break
		}
		if err := d.expect(','); err != nil {
			return err
		}
	}
	setSliceInt64(v, out)
	return nil
}

func encodeStringSlice(e *encoder, v reflect.Value) error {
	if v.IsNil() {
		e.buf = append(e.buf, "null"...)
		return nil
	}
	var s []string
	if v.CanAddr() {
		s = *(*[]string)(v.Addr().UnsafePointer())
	} else {
		s = v.Interface().([]string)
	}
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
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		setSliceString(v, nil)
		return nil
	}
	if err := d.expect('['); err != nil {
		return err
	}
	out := make([]string, 0, estimateArrayLen(d, 4))
	if d.peek() == ']' {
		d.off++
		setSliceString(v, out)
		return nil
	}
	for {
		s, err := d.readString()
		if err != nil {
			return err
		}
		out = append(out, s)
		c := d.peek()
		if c == ']' {
			d.off++
			break
		}
		if err := d.expect(','); err != nil {
			return err
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
	var m map[string]string
	if v.CanAddr() {
		m = *(*map[string]string)(v.Addr().UnsafePointer())
	} else {
		m = v.Interface().(map[string]string)
	}
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
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		setMapStringString(v, nil)
		return nil
	}
	if err := d.expect('{'); err != nil {
		return err
	}
	m := make(map[string]string, 8)
	if d.peek() == '}' {
		d.off++
		setMapStringString(v, m)
		return nil
	}
	for {
		k, err := d.readString()
		if err != nil {
			return err
		}
		if err := d.expect(':'); err != nil {
			return err
		}
		val, err := d.readString()
		if err != nil {
			return err
		}
		m[k] = val
		c := d.peek()
		if c == '}' {
			d.off++
			break
		}
		if err := d.expect(','); err != nil {
			return err
		}
	}
	setMapStringString(v, m)
	return nil
}
