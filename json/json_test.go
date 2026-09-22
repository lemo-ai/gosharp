package json_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/bytedance/sonic"
	jsoniter "github.com/json-iterator/go"
	gosharpjson "github.com/lemo-ai/gosharp/json"
	"github.com/stretchr/testify/assert"
)

type benchPayload struct {
	ID      int64             `json:"id"`
	Name    string            `json:"name"`
	Active  bool              `json:"active"`
	Score   float64           `json:"score"`
	Tags    []string          `json:"tags"`
	Meta    map[string]string `json:"meta"`
	Nested  *benchNested      `json:"nested"`
	Numbers []int             `json:"numbers"`
}

type benchNested struct {
	Title string `json:"title"`
	Count int    `json:"count"`
}

func samplePayload() benchPayload {
	return benchPayload{
		ID:     42,
		Name:   "gosharp-benchmark-payload",
		Active: true,
		Score:  99.5,
		Tags:   []string{"go", "json", "sonic", "util"},
		Meta: map[string]string{
			"env":  "test",
			"lang": "zh-CN",
		},
		Nested: &benchNested{
			Title: "nested",
			Count: 7,
		},
		Numbers: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
}

func largePayload() benchPayload {
	tags := make([]string, 64)
	numbers := make([]int, 256)
	meta := make(map[string]string, 32)
	for i := range tags {
		tags[i] = "tag-" + string(rune('a'+i%26))
	}
	for i := range numbers {
		numbers[i] = i
	}
	for i := 0; i < 32; i++ {
		meta["k"+string(rune('a'+i%26))] = "value-padding-xxxxxxxxxxxx"
	}
	return benchPayload{
		ID:      10086,
		Name:    "gosharp-large-benchmark-payload-with-more-text-content",
		Active:  true,
		Score:   12345.6789,
		Tags:    tags,
		Meta:    meta,
		Nested:  &benchNested{Title: "nested-large", Count: 99},
		Numbers: numbers,
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	in := samplePayload()
	_ = gosharpjson.Pretouch(in)
	b, err := gosharpjson.Marshal(in)
	assert.NoError(t, err)
	assert.True(t, gosharpjson.Valid(b))

	var out benchPayload
	assert.NoError(t, gosharpjson.Unmarshal(b, &out))
	assert.Equal(t, in.ID, out.ID)
	assert.Equal(t, in.Name, out.Name)
	assert.Equal(t, in.Active, out.Active)
	assert.Equal(t, in.Tags, out.Tags)
	assert.Equal(t, in.Meta, out.Meta)
	assert.Equal(t, in.Nested, out.Nested)
	assert.Equal(t, in.Numbers, out.Numbers)

	indented, err := gosharpjson.MarshalIndent(in, "", "  ")
	assert.NoError(t, err)
	assert.Contains(t, string(indented), "\n")

	s, err := gosharpjson.MarshalToString(in)
	assert.NoError(t, err)
	var out2 benchPayload
	assert.NoError(t, gosharpjson.UnmarshalFromString(s, &out2))
	assert.Equal(t, in.Name, out2.Name)
}

func TestEncoderDecoder(t *testing.T) {
	var buf bytes.Buffer
	enc := gosharpjson.NewEncoder(&buf)
	assert.NoError(t, enc.Encode(samplePayload()))

	dec := gosharpjson.NewDecoder(&buf)
	var out benchPayload
	assert.NoError(t, dec.Decode(&out))
	assert.Equal(t, samplePayload().Name, out.Name)
}

func TestCompatWithStdlib(t *testing.T) {
	in := samplePayload()
	ours, err := gosharpjson.Marshal(in)
	assert.NoError(t, err)
	var viaStd benchPayload
	assert.NoError(t, json.Unmarshal(ours, &viaStd))
	assert.Equal(t, in.Name, viaStd.Name)

	stdBytes, err := json.Marshal(in)
	assert.NoError(t, err)
	var viaOurs benchPayload
	assert.NoError(t, gosharpjson.Unmarshal(stdBytes, &viaOurs))
	assert.Equal(t, in.ID, viaOurs.ID)
}

func BenchmarkMarshal(b *testing.B) {
	v := samplePayload()
	_ = gosharpjson.Pretouch(v)
	_ = sonic.Pretouch(reflect.TypeOf(v))

	b.Run("stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := jsoniter.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("sonic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := sonic.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gosharp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := gosharpjson.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUnmarshal(b *testing.B) {
	v := samplePayload()
	raw, err := json.Marshal(v)
	if err != nil {
		b.Fatal(err)
	}
	_ = gosharpjson.Pretouch((*benchPayload)(nil))
	_ = sonic.Pretouch(reflect.TypeOf((*benchPayload)(nil)))

	b.Run("stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchPayload
			if err := json.Unmarshal(raw, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchPayload
			if err := jsoniter.Unmarshal(raw, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("sonic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchPayload
			if err := sonic.Unmarshal(raw, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gosharp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchPayload
			if err := gosharpjson.Unmarshal(raw, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMarshalLarge(b *testing.B) {
	v := largePayload()
	_ = gosharpjson.Pretouch(v)
	_ = sonic.Pretouch(reflect.TypeOf(v))

	b.Run("stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := jsoniter.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("sonic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := sonic.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gosharp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := gosharpjson.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUnmarshalLarge(b *testing.B) {
	v := largePayload()
	raw, err := json.Marshal(v)
	if err != nil {
		b.Fatal(err)
	}
	_ = gosharpjson.Pretouch((*benchPayload)(nil))
	_ = sonic.Pretouch(reflect.TypeOf((*benchPayload)(nil)))

	b.Run("stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchPayload
			if err := json.Unmarshal(raw, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchPayload
			if err := jsoniter.Unmarshal(raw, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("sonic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchPayload
			if err := sonic.Unmarshal(raw, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gosharp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchPayload
			if err := gosharpjson.Unmarshal(raw, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
}
