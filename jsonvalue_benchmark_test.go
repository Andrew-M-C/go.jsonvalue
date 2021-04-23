package jsonvalue

import (
	"encoding/json"
	"testing"

	jsonparser "github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
)

// go test -bench=. -run=none -benchmem -benchtime=10s

var unmarshalText = []byte(`{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!"},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}}},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}`)
var jsonit = jsoniter.ConfigCompatibleWithStandardLibrary

type object struct {
	Int    int       `json:"int"`
	Float  float64   `json:"float"`
	String string    `json:"string"`
	Object *object   `json:"object,omitempty"`
	Array  []*object `json:"array,omitempty"`
}

func BenchmarkGoStdJsonStructUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		o := object{}
		json.Unmarshal(unmarshalText, &o)
	}
}

func BenchmarkGoStdJsonStructMarshal(b *testing.B) {
	o := object{}
	json.Unmarshal(unmarshalText, &o)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&o)
		if err != nil {
			b.Errorf("marshal error: %v", err)
			return
		}
	}
}

func BenchmarkGoStdJsonMapItfUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := map[string]interface{}{}
		json.Unmarshal(unmarshalText, &m)
	}
}

func BenchmarkGoStdJsonMapItfMarshal(b *testing.B) {
	m := map[string]interface{}{}
	json.Unmarshal(unmarshalText, &m)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&m)
		if err != nil {
			b.Errorf("marshal error: %v", err)
			return
		}
	}
}

func BenchmarkJsoniterrStructUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		o := object{}
		jsonit.Unmarshal(unmarshalText, &o)
	}
}

func BenchmarkJsoniterrStructMarshal(b *testing.B) {
	o := object{}
	jsonit.Unmarshal(unmarshalText, &o)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&o)
		if err != nil {
			b.Errorf("marshal error: %v", err)
			return
		}
	}
}

func BenchmarkJsoniterrMapItfUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := map[string]interface{}{}
		jsonit.Unmarshal(unmarshalText, &m)
	}
}

func BenchmarkJsoniterrMapItfMarshal(b *testing.B) {
	m := map[string]interface{}{}
	jsonit.Unmarshal(unmarshalText, &m)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&m)
		if err != nil {
			b.Errorf("marshal error: %v", err)
			return
		}
	}
}

func BenchmarkJsoniterrUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsoniter.Get(unmarshalText)
	}
}

// func BenchmarkJsoniterrMarshal(b *testing.B) {
// 	j := jsoniter.Get(unmarshalText)
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		j.ToString()
// 	}
// }

func BenchmarkJsonparsrUnmarshal(b *testing.B) {
	var objEach func([]byte, []byte, jsonparser.ValueType, int) error
	var arrEach func([]byte, jsonparser.ValueType, int, error)

	objEach = func(k, v []byte, t jsonparser.ValueType, _ int) (noErr error) {
		switch t {
		default:
			// do nothing
		case jsonparser.Array:
			jsonparser.ArrayEach(v, arrEach)
		case jsonparser.Object:
			jsonparser.ObjectEach(v, objEach)
		}
		return
	}

	arrEach = func(v []byte, t jsonparser.ValueType, _ int, _ error) {
		switch t {
		default:
			// do nothing
		case jsonparser.Array:
			jsonparser.ArrayEach(v, arrEach)
		case jsonparser.Object:
			jsonparser.ObjectEach(v, objEach)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsonparser.ObjectEach(unmarshalText, objEach)
	}
}

func BenchmarkJsonvalueUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unmarshal(unmarshalText)
	}
}

func BenchmarkJsonvalueMarshal(b *testing.B) {
	j, _ := Unmarshal(unmarshalText)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := j.Marshal()
		if err != nil {
			b.Errorf("marshal error: %v", err)
			return
		}
	}
}
