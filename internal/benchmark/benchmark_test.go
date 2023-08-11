package benchmark_test

import (
	"testing"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// go test -bench=. -run=none -benchmem -benchtime=2s

var unmarshalText = []byte(`{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!"},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}}},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}`)

func init() {
	jsonvalue.SetDefaultMarshalOptions(jsonvalue.OptUTF8())
}

func Benchmark_Unmarshal_Jsonvalue(b *testing.B) {
	origB := unmarshalText
	for i := 0; i < b.N; i++ {
		_, _ = jsonvalue.Unmarshal(origB)
	}
}
