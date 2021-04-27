package jsonvalue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	// jsonvalue103 "github.com/Andrew-M-C/go.jsonvalue.1.0.3"
	jsonparser "github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
)

// go test -bench=. -run=none -benchmem -benchtime=10s

var unmarshalText = []byte(`{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!"},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}}},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}`)
var jsonit = jsoniter.ConfigCompatibleWithStandardLibrary

func generateLongObject() []byte {
	buff := bytes.Buffer{}

	buff.WriteRune('{')

	for i := 0; i < 100; i++ {
		if i > 0 {
			buff.WriteRune(',')
		}

		buff.WriteRune('"')
		buff.WriteString(strconv.FormatInt(int64(i), 10))
		buff.WriteString(`":`)
		buff.Write(unmarshalText)
	}

	buff.WriteRune(('}'))

	return buff.Bytes()
}

type object struct {
	Int    int       `json:"int"`
	Float  float64   `json:"float"`
	String string    `json:"string"`
	Object *object   `json:"object,omitempty"`
	Array  []*object `json:"array,omitempty"`
}

func Benchmark_Unmarshal_GoStdJsonStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		o := object{}
		json.Unmarshal(unmarshalText, &o)
	}
}

func Benchmark__Marshal__GoStdJsonStruct(b *testing.B) {
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

func Benchmark_Unmarshal_GoStdJsonMapItf_blob(b *testing.B) {
	raw := generateLongObject()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := map[string]interface{}{}
		json.Unmarshal(raw, &m)
	}
}

func Benchmark_Unmarshal_GoStdJsonMapItf(b *testing.B) {
	raw := unmarshalText
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := map[string]interface{}{}
		json.Unmarshal(raw, &m)
	}
}

func Benchmark__Marshal__GoStdJsonMapItf(b *testing.B) {
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

func Benchmark_Unmarshal_JsoniterStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		o := object{}
		jsonit.Unmarshal(unmarshalText, &o)
	}
}

func Benchmark__Marshal__JsoniterStruct(b *testing.B) {
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

func Benchmark_Unmarshal_JsoniterMapItf_blob(b *testing.B) {
	raw := generateLongObject()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := map[string]interface{}{}
		jsonit.Unmarshal(raw, &m)
	}
}

func Benchmark_Unmarshal_JsoniterMapItf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := map[string]interface{}{}
		jsonit.Unmarshal(unmarshalText, &m)
	}
}

func Benchmark__Marshal__JsoniterMapItf(b *testing.B) {
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

func Benchmark____Get____Jsoniter(b *testing.B) {
	raw := unmarshalText
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsoniter.Get(raw)
	}
}

func Benchmark____Get____Jsoniter_AndGetParsedValue(b *testing.B) {
	raw := unmarshalText
	any := jsoniter.Get(raw)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		any.Get("object", "object", "object", "array", 1)
	}
}

func Benchmark____Get____Jsoniter_ReadOneChain(b *testing.B) {
	raw := unmarshalText
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		any := jsoniter.Get(raw)
		any.Get("object", "object", "object", "array", 1)
	}
}

func Benchmark____Get____Jsoniter_ReadLevelOne(b *testing.B) {
	raw := unmarshalText
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		any := jsoniter.Get(raw)
		keys := any.Keys()
		for _, k := range keys {
			any.Get(k)
		}
	}
}

func Benchmark____Get____Jsoniter_blob_ReadOneChain(b *testing.B) {
	raw := generateLongObject()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		any := jsoniter.Get(raw)
		any.Get("object", "object", "object", "array", 1)
	}
}

func Benchmark_Unmarshal_Jsonparser_Full_Blob(b *testing.B) {
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

	raw := generateLongObject()
	b.ResetTimer()

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
		jsonparser.ObjectEach(raw, objEach)
	}
}

func Benchmark_Unmarshal_Jsonparser_ReadLevelOne_Blob(b *testing.B) {
	objEach := func(k, v []byte, t jsonparser.ValueType, _ int) (noErr error) {
		return
	}

	raw := generateLongObject()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonparser.ObjectEach(raw, objEach)
	}
}

func Benchmark_Unmarshal_Jsonparser_ReadOneChain(b *testing.B) {
	raw := unmarshalText
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonparser.Get(raw, "object", "object", "object", "array", "[1]")
	}
}

// func Benchmark_Unmarshal_Jsonvalue_v1_0_3(b *testing.B) {
// 	origB := unmarshalText
// 	for i := 0; i < b.N; i++ {
// 		jsonvalue103.Unmarshal(origB)
// 	}
// }

func Benchmark_Unmarshal_Jsonvalue(b *testing.B) {
	origB := unmarshalText
	for i := 0; i < b.N; i++ {
		Unmarshal(origB)
	}
}

func Benchmark_Unmarshal_Jsonvalue_GetByInterface(b *testing.B) {
	v, _ := Unmarshal(unmarshalText)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v.Get("object", "object", "object", "array", 1)
	}
}

func Benchmark_Unmarshal_Jsonvalue_GetByInterfaceCaseless(b *testing.B) {
	v, _ := Unmarshal(unmarshalText)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v.Caseless().Get("object", "object", "object", "array", 1)
	}
}

func Benchmark_Unmarshal_Jsonvalue_GetByString(b *testing.B) {
	v, _ := Unmarshal(unmarshalText)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sGet(v, "object", "object", "object", "array", index(1))
	}
}

func Benchmark_Unmarshal_Jsonvalue_blob(b *testing.B) {
	// origB := unmarshalText
	origB := generateLongObject()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Unmarshal(origB)
	}
}

func Benchmark_Unmarshal_JsonvalueNoCopy(b *testing.B) {
	// origB := unmarshalText
	origB := generateLongObject()

	lst := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		bytes := make([]byte, len(origB))
		copy(bytes, origB)
		lst[i] = bytes
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// err := jsonit.Unmarshal(raw, &s)
		UnmarshalNoCopy(lst[i])
	}
}

func Benchmark__Marshal__Jsonvalue(b *testing.B) {
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

// ================
// functions below is a experimental SGet() function, which use a function that acccept string
// parameters only instead of interface{}.

func index(i int) string {
	return fmt.Sprintf("[%d]", i)
}

func isIndex(s string) (int, bool) {
	le := len(s)
	if le == 0 {
		return 0, false
	}

	if s[0] == '[' && s[le-1] == ']' {
		// continue
	} else {
		return 0, false
	}

	i, err := strconv.ParseInt(s[1:le-1], 10, 64)
	if err != nil {
		return 0, false
	}
	return int(i), true
}

func sGet(v *V, firstParam string, otherParams ...string) (*V, error) {
	child, err := sGetInCurrValue(v, firstParam)
	if err != nil {
		return nil, err
	}

	if len(otherParams) == 0 {
		return child, nil
	}
	return sGet(child, otherParams[0], otherParams[1:]...)
}

func sGetInCurrValue(v *V, param string) (*V, error) {
	if v.valueType == jsonparser.Array {
		// integer expected
		pos, ok := isIndex(param)
		if !ok {
			return nil, fmt.Errorf("%w, expected a number string but got '%s'", ErrNotValidNumberValue, param)
		}
		child, _ := v.childAtIndex(pos)
		if nil == child {
			return nil, ErrOutOfRange
		}
		return child, nil

	} else if v.valueType == jsonparser.Object {
		// string expected
		child, exist := v.getFromObjectChildren(false, param)
		if !exist {
			return nil, ErrNotFound
		}
		return child, nil

	} else {
		return nil, fmt.Errorf("%v type does not supports Get()", v.valueType)
	}
}
