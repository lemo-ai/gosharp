package json

import (
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

// codec is a compiled encode/decode pair for one concrete type.
type codec struct {
	encode encodeFn
	decode decodeFn
}

type encodeFn func(e *encoder, v reflect.Value) error
type decodeFn func(d *decoder, v reflect.Value) error

var (
	codecCache sync.Map // map[reflect.Type]*atomic.Pointer[codec]
)

func getCodec(t reflect.Type) *codec {
	if t == nil {
		return nil
	}
	if p, ok := codecCache.Load(t); ok {
		return p.(*atomic.Pointer[codec]).Load()
	}
	ap := &atomic.Pointer[codec]{}
	actual, loaded := codecCache.LoadOrStore(t, ap)
	if loaded {
		ap = actual.(*atomic.Pointer[codec])
		if c := ap.Load(); c != nil {
			return c
		}
	}
	c := compileCodec(t)
	ap.Store(c)
	return c
}

func pretouchType(t reflect.Type) {
	if t == nil {
		return
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	_ = getCodec(t)
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i).Type
			pretouchType(ft)
		}
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Pointer:
		pretouchType(t.Elem())
		if t.Kind() == reflect.Map {
			pretouchType(t.Key())
		}
	}
}

type structFieldCodec struct {
	name      string
	nameBytes []byte // `"name":`
	index     []int
	omitEmpty bool
	codec     *codec
	typ       reflect.Type
}

type structCodec struct {
	fields    []structFieldCodec
	fieldMap  map[string]*structFieldCodec
	encode    encodeFn
	decode    decodeFn
}

func compileStruct(t reflect.Type) *structCodec {
	sc := &structCodec{
		fieldMap: make(map[string]*structFieldCodec),
	}
	var walk func(t reflect.Type, index []int)
	walk = func(t reflect.Type, index []int) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" && !f.Anonymous { // unexported
				continue
			}
			tag := f.Tag.Get("json")
			if tag == "-" {
				continue
			}
			name, opts := parseTag(tag)
			if name == "" {
				name = f.Name
			}
			idx := append(append([]int{}, index...), i)
			if f.Anonymous && name == f.Name && f.Type.Kind() == reflect.Struct && opts == "" {
				walk(f.Type, idx)
				continue
			}
			if f.Anonymous && f.Type.Kind() == reflect.Pointer && f.Type.Elem().Kind() == reflect.Struct && name == f.Name {
				walk(f.Type.Elem(), idx)
				continue
			}
			sfc := structFieldCodec{
				name:      name,
				nameBytes: []byte(`"` + name + `":`),
				index:     idx,
				omitEmpty: optsContains(opts, "omitempty"),
				codec:     getCodec(f.Type),
				typ:       f.Type,
			}
			sc.fields = append(sc.fields, sfc)
		}
	}
	walk(t, nil)
	// Rebuild fieldMap after all appends so pointers stay stable.
	sc.fieldMap = make(map[string]*structFieldCodec, len(sc.fields))
	for i := range sc.fields {
		sc.fieldMap[sc.fields[i].name] = &sc.fields[i]
	}
	sc.encode = makeEncodeStruct(sc)
	sc.decode = makeDecodeStruct(sc)
	return sc
}

func parseTag(tag string) (name, opts string) {
	if tag == "" {
		return "", ""
	}
	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' {
			return tag[:i], tag[i+1:]
		}
	}
	return tag, ""
}

func optsContains(opts, want string) bool {
	for opts != "" {
		var o string
		o, opts, _ = cutComma(opts)
		if o == want {
			return true
		}
	}
	return false
}

func cutComma(s string) (before, after string, found bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return s[:i], s[i+1:], true
		}
	}
	return s, "", false
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
