package jsonvalue

import (
	"encoding/json"
	"testing"
)

// go test -bench=. -run=none -benchmem -benchtime=10s

var unmarshalText = []byte(`{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!"},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}}},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}`)

type object struct {
	Int    int       `json:"int"`
	Float  float64   `json:"float"`
	String string    `json:"string"`
	Object *object   `json:"object,omitempty"`
	Array  []*object `json:"array,omitempty"`
}

func BenchmarkUnmarshalGoMapInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := map[string]interface{}{}
		err := json.Unmarshal(unmarshalText, &m)
		if err != nil {
			b.Errorf("unmarshal error: %v", err)
			return
		}
	}
	return
}

func BenchmarkUnmarshalJsonvalue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Unmarshal(unmarshalText)
		if err != nil {
			b.Errorf("unmarshal error: %v", err)
			return
		}
	}
}

func BenchmarkUnmarshalStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		o := object{}
		err := json.Unmarshal(unmarshalText, &o)
		if err != nil {
			b.Errorf("unmarshal error: %v", err)
			return
		}
	}
}
