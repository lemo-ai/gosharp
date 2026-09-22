package json

import "errors"

var errWasNull = errors.New("json: null")

// beginArrayOrNull consumes '[' or reports null. On success, off points at first
// element or ']'. Returns errWasNull when the value is null.
func (d *decoder) beginArrayOrNull() error {
	d.skipSpace()
	if d.off >= len(d.data) {
		return syntaxError(d.off, "unexpected EOF")
	}
	switch d.data[d.off] {
	case 'n':
		if err := d.consumeLiteral("null"); err != nil {
			return err
		}
		return errWasNull
	case '[':
		d.off++
		d.skipSpace()
		if d.off >= len(d.data) {
			return syntaxError(d.off, "unterminated array")
		}
		return nil
	default:
		return syntaxError(d.off, "expect '['")
	}
}

// afterValueInArray skips spaces and consumes ',' or ']'. Returns the delimiter.
func (d *decoder) afterValueInArray() (byte, error) {
	d.skipSpace()
	if d.off >= len(d.data) {
		return 0, syntaxError(d.off, "unterminated array")
	}
	c := d.data[d.off]
	d.off++
	switch c {
	case ']', ',':
		return c, nil
	default:
		return 0, syntaxError(d.off, "expect ',' or ']'")
	}
}

// readStringFast is optimized for compact JSON without escapes (common path).
func (d *decoder) readStringFast() (string, error) {
	d.skipSpace()
	if d.off >= len(d.data) || d.data[d.off] != '"' {
		return "", syntaxError(d.off, "expect '\"'")
	}
	d.off++
	data := d.data
	start := d.off
	i := start
	for i < len(data) {
		c := data[i]
		if c == '"' {
			s := b2s(data[start:i])
			d.off = i + 1
			return s, nil
		}
		if c == '\\' || c < 0x20 {
			d.off = start
			return d.readStringSlow(start)
		}
		i++
	}
	return "", syntaxError(d.off, "unterminated string")
}
