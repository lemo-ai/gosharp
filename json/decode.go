package json

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

type decoder struct {
	data []byte
	off  int
}

func (d *decoder) skipSpace() {
	for d.off < len(d.data) {
		c := d.data[d.off]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			d.off++
			continue
		}
		return
	}
}

func (d *decoder) peek() byte {
	d.skipSpace()
	if d.off >= len(d.data) {
		return 0
	}
	return d.data[d.off]
}

func (d *decoder) next() byte {
	d.skipSpace()
	if d.off >= len(d.data) {
		return 0
	}
	c := d.data[d.off]
	d.off++
	return c
}

func (d *decoder) expect(c byte) error {
	if d.next() != c {
		return syntaxError(d.off, "expect %q", c)
	}
	return nil
}

func Unmarshal(data []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(v)}
	}
	d := decoder{data: data}
	d.skipSpace()
	if err := decodeValue(&d, rv.Elem()); err != nil {
		return err
	}
	d.skipSpace()
	if d.off != len(d.data) {
		return syntaxError(d.off, "trailing garbage")
	}
	return nil
}

func UnmarshalFromString(str string, v any) error {
	return Unmarshal(s2b(str), v)
}

func decodeValue(d *decoder, v reflect.Value) error {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	c := getCodec(v.Type())
	return c.decode(d, v)
}

func decodeBool(d *decoder, v reflect.Value) error {
	switch d.peek() {
	case 't':
		if err := d.consumeLiteral("true"); err != nil {
			return err
		}
		v.SetBool(true)
	case 'f':
		if err := d.consumeLiteral("false"); err != nil {
			return err
		}
		v.SetBool(false)
	case 'n':
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		v.SetBool(false)
	default:
		return syntaxError(d.off, "invalid bool")
	}
	return nil
}

func decodeInt(d *decoder, v reflect.Value) error {
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		v.SetInt(0)
		return nil
	}
	n, err := d.readNumber()
	if err != nil {
		return err
	}
	if i, ok := parseIntFast(n); ok {
		v.SetInt(i)
		return nil
	}
	f, err := parseFloatBytes(n)
	if err != nil {
		return err
	}
	v.SetInt(int64(f))
	return nil
}

func decodeUint(d *decoder, v reflect.Value) error {
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		v.SetUint(0)
		return nil
	}
	n, err := d.readNumber()
	if err != nil {
		return err
	}
	u, err := strconv.ParseUint(b2s(n), 10, 64)
	if err != nil {
		f, e2 := strconv.ParseFloat(b2s(n), 64)
		if e2 != nil {
			return err
		}
		u = uint64(f)
	}
	v.SetUint(u)
	return nil
}

func decodeFloat(d *decoder, v reflect.Value) error {
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		v.SetFloat(0)
		return nil
	}
	n, err := d.readNumber()
	if err != nil {
		return err
	}
	f, err := strconv.ParseFloat(b2s(n), 64)
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}

func decodeString(d *decoder, v reflect.Value) error {
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		v.SetString("")
		return nil
	}
	s, err := d.readStringFast()
	if err != nil {
		return err
	}
	v.SetString(s)
	return nil
}

func decodeBytes(d *decoder, v reflect.Value) error {
	if d.peek() == 'n' {
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		v.SetBytes(nil)
		return nil
	}
	if d.peek() == '[' {
		// array of numbers fallback
		var tmp []byte
		tv := reflect.ValueOf(&tmp).Elem()
		elem := getCodec(reflect.TypeOf(byte(0)))
		if err := makeDecodeSlice(elem)(d, tv); err != nil {
			return err
		}
		v.SetBytes(tmp)
		return nil
	}
	s, err := d.readString()
	if err != nil {
		return err
	}
	decoded, err := decodeBase64(s)
	if err != nil {
		return err
	}
	v.SetBytes(decoded)
	return nil
}

func decodeInterface(d *decoder, v reflect.Value) error {
	val, err := d.readAny()
	if err != nil {
		return err
	}
	if val == nil {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	v.Set(reflect.ValueOf(val))
	return nil
}

func decodeUnsupported(d *decoder, v reflect.Value) error {
	return &UnsupportedTypeError{Type: v.Type()}
}

func makeDecodePtr(elem *codec) decodeFn {
	return func(d *decoder, v reflect.Value) error {
		if d.peek() == 'n' {
			if err := d.consumeLiteral("null"); err != nil {
				return err
			}
			v.Set(reflect.Zero(v.Type()))
			return nil
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return elem.decode(d, v.Elem())
	}
}

func makeDecodeSlice(elem *codec) decodeFn {
	return func(d *decoder, v reflect.Value) error {
		if d.peek() == 'n' {
			if err := d.consumeLiteral("null"); err != nil {
				return err
			}
			v.Set(reflect.Zero(v.Type()))
			return nil
		}
		if err := d.expect('['); err != nil {
			return err
		}
		typ := v.Type()
		if d.peek() == ']' {
			d.off++
			v.Set(reflect.MakeSlice(typ, 0, 0))
			return nil
		}
		const initCap = 8
		slice := reflect.MakeSlice(typ, initCap, initCap)
		n := 0
		for {
			if n == slice.Cap() {
				news := reflect.MakeSlice(typ, slice.Cap()*2, slice.Cap()*2)
				reflect.Copy(news, slice)
				slice = news
			}
			if err := elem.decode(d, slice.Index(n)); err != nil {
				return err
			}
			n++
			c := d.peek()
			if c == ']' {
				d.off++
				break
			}
			if err := d.expect(','); err != nil {
				return err
			}
		}
		v.Set(slice.Slice(0, n))
		return nil
	}
}

func makeDecodeArray(elem *codec) decodeFn {
	return func(d *decoder, v reflect.Value) error {
		if err := d.expect('['); err != nil {
			return err
		}
		n := v.Len()
		for i := 0; i < n; i++ {
			if i > 0 {
				if err := d.expect(','); err != nil {
					return err
				}
			}
			if err := elem.decode(d, v.Index(i)); err != nil {
				return err
			}
		}
		return d.expect(']')
	}
}

func makeDecodeMap(t reflect.Type) decodeFn {
	if t.Key().Kind() != reflect.String {
		return decodeUnsupported
	}
	vc := getCodec(t.Elem())
	return func(d *decoder, v reflect.Value) error {
		if d.peek() == 'n' {
			if err := d.consumeLiteral("null"); err != nil {
				return err
			}
			v.Set(reflect.Zero(v.Type()))
			return nil
		}
		if err := d.expect('{'); err != nil {
			return err
		}
		m := reflect.MakeMap(t)
		if d.peek() == '}' {
			d.off++
			v.Set(m)
			return nil
		}
		for {
			key, err := d.readString()
			if err != nil {
				return err
			}
			if err := d.expect(':'); err != nil {
				return err
			}
			elemV := reflect.New(t.Elem()).Elem()
			if err := vc.decode(d, elemV); err != nil {
				return err
			}
			m.SetMapIndex(reflect.ValueOf(key), elemV)
			c := d.peek()
			if c == '}' {
				d.off++
				break
			}
			if err := d.expect(','); err != nil {
				return err
			}
		}
		v.Set(m)
		return nil
	}
}

func makeDecodeStruct(sc *structCodec) decodeFn {
	return func(d *decoder, v reflect.Value) error {
		if d.peek() == 'n' {
			if err := d.consumeLiteral("null"); err != nil {
				return err
			}
			return nil
		}
		if err := d.expect('{'); err != nil {
			return err
		}
		if d.peek() == '}' {
			d.off++
			return nil
		}
		for {
			key, err := d.readStringFast()
			if err != nil {
				return err
			}
			if err := d.expect(':'); err != nil {
				return err
			}
			if f, ok := sc.fieldMap[key]; ok {
				var fv reflect.Value
				if len(f.index) == 1 {
					fv = v.Field(f.index[0])
				} else {
					fv = v.FieldByIndex(f.index)
				}
				if err := f.codec.decode(d, fv); err != nil {
					return err
				}
			} else {
				if err := d.skipValue(); err != nil {
					return err
				}
			}
			c := d.peek()
			if c == '}' {
				d.off++
				break
			}
			if err := d.expect(','); err != nil {
				return err
			}
		}
		return nil
	}
}

func (d *decoder) consumeLiteral(lit string) error {
	if d.off+len(lit) > len(d.data) || b2s(d.data[d.off:d.off+len(lit)]) != lit {
		return syntaxError(d.off, "invalid literal %s", lit)
	}
	d.off += len(lit)
	return nil
}

func (d *decoder) readNumber() ([]byte, error) {
	d.skipSpace()
	start := d.off
	if d.off < len(d.data) && (d.data[d.off] == '-' || d.data[d.off] == '+') {
		d.off++
	}
	for d.off < len(d.data) {
		c := d.data[d.off]
		if (c >= '0' && c <= '9') || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-' {
			d.off++
			continue
		}
		break
	}
	if start == d.off {
		return nil, syntaxError(d.off, "invalid number")
	}
	return d.data[start:d.off], nil
}

func (d *decoder) readString() (string, error) {
	if err := d.expect('"'); err != nil {
		return "", err
	}
	start := d.off
	for d.off < len(d.data) {
		c := d.data[d.off]
		if c == '"' {
			s := b2s(d.data[start:d.off])
			d.off++
			return s, nil
		}
		if c == '\\' || c < 0x20 {
			return d.readStringSlow(start)
		}
		d.off++
	}
	return "", syntaxError(d.off, "unterminated string")
}

func (d *decoder) readStringSlow(start int) (string, error) {
	var b []byte
	if start < d.off {
		b = append(b, d.data[start:d.off]...)
	}
	for d.off < len(d.data) {
		c := d.data[d.off]
		if c == '"' {
			d.off++
			return b2s(b), nil
		}
		if c == '\\' {
			d.off++
			if d.off >= len(d.data) {
				return "", syntaxError(d.off, "unterminated escape")
			}
			esc := d.data[d.off]
			d.off++
			switch esc {
			case '"', '\\', '/':
				b = append(b, esc)
			case 'b':
				b = append(b, '\b')
			case 'f':
				b = append(b, '\f')
			case 'n':
				b = append(b, '\n')
			case 'r':
				b = append(b, '\r')
			case 't':
				b = append(b, '\t')
			case 'u':
				if d.off+4 > len(d.data) {
					return "", syntaxError(d.off, "invalid \\u escape")
				}
				r, err := strconv.ParseUint(b2s(d.data[d.off:d.off+4]), 16, 16)
				if err != nil {
					return "", err
				}
				d.off += 4
				rr := rune(r)
				if utf16.IsSurrogate(rr) {
					// handle surrogate pair
					if d.off+6 <= len(d.data) && d.data[d.off] == '\\' && d.data[d.off+1] == 'u' {
						r2, err := strconv.ParseUint(b2s(d.data[d.off+2:d.off+6]), 16, 16)
						if err == nil {
							d.off += 6
							rr = utf16.DecodeRune(rr, rune(r2))
						}
					}
				}
				var tmp [utf8.UTFMax]byte
				n := utf8.EncodeRune(tmp[:], rr)
				b = append(b, tmp[:n]...)
			default:
				return "", syntaxError(d.off, "invalid escape")
			}
			continue
		}
		if c < 0x20 {
			return "", syntaxError(d.off, "invalid control char")
		}
		b = append(b, c)
		d.off++
	}
	return "", syntaxError(d.off, "unterminated string")
}

func (d *decoder) readAny() (any, error) {
	switch d.peek() {
	case 'n':
		return nil, d.consumeLiteral("null")
	case 't':
		return true, d.consumeLiteral("true")
	case 'f':
		return false, d.consumeLiteral("false")
	case '"':
		return d.readString()
	case '{':
		return d.readObjectAny()
	case '[':
		return d.readArrayAny()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		n, err := d.readNumber()
		if err != nil {
			return nil, err
		}
		if bytes.IndexByte(n, '.') >= 0 || bytes.IndexByte(n, 'e') >= 0 || bytes.IndexByte(n, 'E') >= 0 {
			return strconv.ParseFloat(b2s(n), 64)
		}
		i, err := strconv.ParseInt(b2s(n), 10, 64)
		if err == nil {
			return float64(i), nil // encoding/json compatible: numbers -> float64
		}
		return strconv.ParseFloat(b2s(n), 64)
	default:
		return nil, syntaxError(d.off, "invalid value")
	}
}

func (d *decoder) readObjectAny() (any, error) {
	if err := d.expect('{'); err != nil {
		return nil, err
	}
	m := make(map[string]any)
	if d.peek() == '}' {
		d.off++
		return m, nil
	}
	for {
		k, err := d.readString()
		if err != nil {
			return nil, err
		}
		if err := d.expect(':'); err != nil {
			return nil, err
		}
		v, err := d.readAny()
		if err != nil {
			return nil, err
		}
		m[k] = v
		c := d.peek()
		if c == '}' {
			d.off++
			return m, nil
		}
		if err := d.expect(','); err != nil {
			return nil, err
		}
	}
}

func (d *decoder) readArrayAny() (any, error) {
	if err := d.expect('['); err != nil {
		return nil, err
	}
	var a []any
	if d.peek() == ']' {
		d.off++
		return a, nil
	}
	for {
		v, err := d.readAny()
		if err != nil {
			return nil, err
		}
		a = append(a, v)
		c := d.peek()
		if c == ']' {
			d.off++
			return a, nil
		}
		if err := d.expect(','); err != nil {
			return nil, err
		}
	}
}

func (d *decoder) skipValue() error {
	switch d.peek() {
	case '"':
		_, err := d.readString()
		return err
	case 't':
		return d.consumeLiteral("true")
	case 'f':
		return d.consumeLiteral("false")
	case 'n':
		return d.consumeLiteral("null")
	case '{':
		d.off++
		if d.peek() == '}' {
			d.off++
			return nil
		}
		for {
			if _, err := d.readString(); err != nil {
				return err
			}
			if err := d.expect(':'); err != nil {
				return err
			}
			if err := d.skipValue(); err != nil {
				return err
			}
			c := d.peek()
			if c == '}' {
				d.off++
				return nil
			}
			if err := d.expect(','); err != nil {
				return err
			}
		}
	case '[':
		d.off++
		if d.peek() == ']' {
			d.off++
			return nil
		}
		for {
			if err := d.skipValue(); err != nil {
				return err
			}
			c := d.peek()
			if c == ']' {
				d.off++
				return nil
			}
			if err := d.expect(','); err != nil {
				return err
			}
		}
	default:
		_, err := d.readNumber()
		return err
	}
}

func decodeBase64(s string) ([]byte, error) {
	const decodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var dec [256]byte
	for i := range dec {
		dec[i] = 0xFF
	}
	for i := 0; i < len(decodeStd); i++ {
		dec[decodeStd[i]] = byte(i)
	}
	if len(s)%4 != 0 {
		return nil, fmt.Errorf("invalid base64")
	}
	outLen := len(s) / 4 * 3
	if len(s) >= 1 && s[len(s)-1] == '=' {
		outLen--
	}
	if len(s) >= 2 && s[len(s)-2] == '=' {
		outLen--
	}
	out := make([]byte, 0, outLen)
	for i := 0; i < len(s); i += 4 {
		var n uint32
		pad := 0
		for j := 0; j < 4; j++ {
			c := s[i+j]
			if c == '=' {
				pad++
				n <<= 6
				continue
			}
			v := dec[c]
			if v == 0xFF {
				return nil, fmt.Errorf("invalid base64")
			}
			n = n<<6 | uint32(v)
		}
		out = append(out, byte(n>>16))
		if pad < 2 {
			out = append(out, byte(n>>8))
		}
		if pad < 1 {
			out = append(out, byte(n))
		}
	}
	return out, nil
}

func syntaxError(off int, format string, args ...any) error {
	return &SyntaxError{Offset: int64(off), Message: fmt.Sprintf(format, args...)}
}

// Valid reports whether data is valid JSON.
func Valid(data []byte) bool {
	d := decoder{data: data}
	if _, err := d.readAny(); err != nil {
		return false
	}
	d.skipSpace()
	return d.off == len(d.data)
}

// Encoder writes JSON values to an output stream.
type Encoder struct {
	w   io.Writer
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(v any) error {
	if enc.err != nil {
		return enc.err
	}
	b, err := Marshal(v)
	if err != nil {
		enc.err = err
		return err
	}
	if _, err = enc.w.Write(b); err != nil {
		enc.err = err
		return err
	}
	_, err = enc.w.Write([]byte{'\n'})
	enc.err = err
	return err
}

// Decoder reads JSON values from an input stream.
type Decoder struct {
	r   io.Reader
	buf []byte
	err error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (dec *Decoder) Decode(v any) error {
	if dec.err != nil {
		return dec.err
	}
	if len(dec.buf) == 0 {
		b, err := io.ReadAll(dec.r)
		if err != nil {
			dec.err = err
			return err
		}
		dec.buf = b
	}
	d := decoder{data: dec.buf}
	d.skipSpace()
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(v)}
	}
	if err := decodeValue(&d, rv.Elem()); err != nil {
		dec.err = err
		return err
	}
	dec.buf = dec.buf[d.off:]
	return nil
}
