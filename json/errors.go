package json

import (
	"fmt"
	"reflect"
)

// SyntaxError is a syntax error in JSON.
type SyntaxError struct {
	Offset  int64
	Message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("json: %s at offset %d", e.Message, e.Offset)
}

// UnsupportedTypeError is returned when marshaling unsupported types.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "json: unsupported type: " + e.Type.String()
}

// InvalidUnmarshalError is returned by Unmarshal for invalid targets.
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}
