package jsonvalue

import (
	"strconv"
	"testing"
)

// go test -bench=. -run=none -benchmem

func BenchmarkGoStdParseFloat(b *testing.B) {
	s := "-12345.67890"

	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(s, 64)
	}
}

func BenchmarkGoStdParseInt(b *testing.B) {
	s := "-1234567890"

	for i := 0; i < b.N; i++ {
		strconv.ParseInt(s, 10, 64)
	}
}

func BenchmarkIterParseFloat(b *testing.B) {
	it := &iter{b: []byte("-12345.67890")}

	for i := 0; i < b.N; i++ {
		it.parseNumber(0)
	}
}

func BenchmarkIterParseInt(b *testing.B) {
	it := &iter{b: []byte("-1234567890")}

	for i := 0; i < b.N; i++ {
		it.parseNumber(0)
	}
}
