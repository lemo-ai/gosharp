package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

// ConfigCompatibleWithStandardLibrary tries to be compatible
// with standard library behavior.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Marshal adapts to encoding/json Marshal API.
//
// Marshal returns the JSON encoding of v.
// Refer to https://pkg.go.dev/encoding/json#Marshal for more information.
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalToString convenient method to write as string instead of []byte.
func MarshalToString(v interface{}) (string, error) {
	return json.MarshalToString(v)
}

// MarshalIndent same as json.MarshalIndent.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// UnmarshalFromString is a convenient method to read from string instead of []byte.
func UnmarshalFromString(str string, v interface{}) error {
	return json.UnmarshalFromString(str, v)
}

// Unmarshal adapts to encoding/json Unmarshal API.
//
// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// Refer to https://pkg.go.dev/encoding/json#Unmarshal for more information.
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// NewEncoder same as json.NewEncoder.
func NewEncoder(writer io.Writer) *jsoniter.Encoder {
	return json.NewEncoder(writer)
}

// NewDecoder adapts to encoding/json NewDecoder API.
//
// NewDecoder returns a new decoder that reads from r.
func NewDecoder(reader io.Reader) *jsoniter.Decoder {
	return json.NewDecoder(reader)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
	return jsoniter.Valid(data)
}
