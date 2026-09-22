package convert

import "reflect"

// assignConverted converts value to dst's type via Convert and assigns if possible.
// Returns true when assignment succeeded without panicking.
func assignConverted(dst reflect.Value, value interface{}) bool {
	if !dst.IsValid() || !dst.CanSet() {
		return false
	}
	converted := Convert(value, dst.Type().String())
	cv := reflect.ValueOf(converted)
	if !cv.IsValid() {
		return false
	}
	if cv.Type().AssignableTo(dst.Type()) {
		dst.Set(cv)
		return true
	}
	if cv.CanConvert(dst.Type()) {
		dst.Set(cv.Convert(dst.Type()))
		return true
	}
	return false
}

// convertToType converts value to the given reflect type without panicking.
func convertToType(value interface{}, typ reflect.Type) (reflect.Value, bool) {
	converted := Convert(value, typ.String())
	cv := reflect.ValueOf(converted)
	if !cv.IsValid() {
		return reflect.Value{}, false
	}
	if cv.Type().AssignableTo(typ) {
		return cv, true
	}
	if cv.CanConvert(typ) {
		return cv.Convert(typ), true
	}
	src := reflect.ValueOf(value)
	if src.IsValid() && src.CanConvert(typ) {
		return src.Convert(typ), true
	}
	return reflect.Value{}, false
}
