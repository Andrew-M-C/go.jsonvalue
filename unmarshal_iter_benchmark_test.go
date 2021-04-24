package jsonvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

// go test -bench=. -run=none -benchmem

func BenchmarkGoStdParseFloat(b *testing.B) {
	s := "-12345.67890"

	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(s, 64)
	}
}

func BenchmarkGoStdParseInt(b *testing.B) {
	s := "1234567890"

	for i := 0; i < b.N; i++ {
		strconv.ParseUint(s, 10, 64)
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

func BenchmarkJsitParseFloat(b *testing.B) {
	raw := []byte("-12345.67890")

	for i := 0; i < b.N; i++ {
		jsoniter.Get(raw).ToFloat64()
	}
}

func BenchmarkJsitParseInt(b *testing.B) {
	raw := []byte("-1234567890")

	for i := 0; i < b.N; i++ {
		jsoniter.Get(raw).ToInt64()
	}
}

func BenchmarkJsonParseString(b *testing.B) {
	// orig := "你👨‍👩‍👧‍👧你"
	orig := fmt.Sprintf(
		"%c%c%c%c%c%c%c%c%c",
		0x2F804, 0x1F468, 0x200D, 0x1F469, 0x200D, 0x1F467, 0x200D, 0x1F467, 0x4F60,
	)

	raw := []byte(fmt.Sprintf(`{"string":"%s"}`, orig))

	data := struct {
		String string `json:"string"`
	}{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// err := jsonit.Unmarshal(raw, &s)
		json.Unmarshal(raw, &data)
	}
}

func BenchmarkJsitParseString(b *testing.B) {
	// orig := "你👨‍👩‍👧‍👧你"
	orig := fmt.Sprintf(
		"%c%c%c%c%c%c%c%c%c",
		0x2F804, 0x1F468, 0x200D, 0x1F469, 0x200D, 0x1F467, 0x200D, 0x1F467, 0x4F60,
	)

	raw := []byte(fmt.Sprintf(`{"string":"%s"}`, orig))

	data := struct {
		String string `json:"string"`
	}{}

	jsonit := jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// err := jsonit.Unmarshal(raw, &s)
		jsonit.Unmarshal(raw, &data)
	}
}

func BenchmarkIterParseString(b *testing.B) {
	orig := fmt.Sprintf(
		":::::string%c%c%c%c%c%c%c%c%c",
		0x2F804, 0x1F468, 0x200D, 0x1F469, 0x200D, 0x1F467, 0x200D, 0x1F467, 0x4F60,
	)

	origB := []byte(orig)

	itLst := make([]*iter, b.N)
	for i := 0; i < b.N; i++ {
		bytes := make([]byte, len(origB))
		copy(bytes, origB)
		itLst[i] = &iter{b: bytes}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// err := jsonit.Unmarshal(raw, &s)
		it := itLst[0]
		it.parseStrFromBytesBackward(0, len(it.b))
	}
}

func BenchmarkJsitParseSimpleString(b *testing.B) {
	raw := []byte(`{"string":"hello"}`)

	data := struct {
		String string `json:"string"`
	}{}

	jsonit := jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// err := jsonit.Unmarshal(raw, &s)
		jsonit.Unmarshal(raw, &data)
	}
}

func BenchmarkIterParseSimpleString(b *testing.B) {
	orig := "hello, world, hello!"

	origB := []byte(orig)

	itLst := make([]*iter, b.N)
	for i := 0; i < b.N; i++ {
		bytes := make([]byte, len(origB))
		copy(bytes, origB)
		itLst[i] = &iter{b: bytes}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// err := jsonit.Unmarshal(raw, &s)
		it := itLst[0]
		it.parseStrFromBytesBackward(0, len(it.b))
	}
}

func BenchmarkJsitParseObject(b *testing.B) {
	raw := unmarshalText

	data := struct {
		String string `json:"string"`
	}{}

	jsonit := jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// err := jsonit.Unmarshal(raw, &s)
		jsonit.Unmarshal(raw, &data)
	}
}

func BenchmarkIterParseObject(b *testing.B) {
	origB := unmarshalText

	itLst := make([]*iter, b.N)
	for i := 0; i < b.N; i++ {
		bytes := make([]byte, len(origB))
		copy(bytes, origB)
		itLst[i] = &iter{b: bytes}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// err := jsonit.Unmarshal(raw, &s)
		it := itLst[i]
		unmarshalWithIter(it, 0, len(it.b))
	}
}

func BenchmarkJsvlParseObject(b *testing.B) {
	raw := unmarshalText
	for i := 0; i < b.N; i++ {
		Unmarshal(raw)
	}
}
