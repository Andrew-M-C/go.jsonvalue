package jsonvalue

import (
	"encoding/json"
	"testing"
)

func TestBasicFunction(t *testing.T) {
	raw := `{"message":"hello, 世界","float":1234.123456789123456789,"true":true,"false":false,"null":null,"obj":{"msg":"hi"},"arr":["你好","world",null],"uint":1234,"int":-1234}`

	v, err := Unmarshal([]byte(raw))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	t.Logf("OK: %+v", v)

	b, err := v.Marshal()
	if err != nil {
		t.Errorf("marshal failed: %v", err)
		return
	}
	t.Logf("marshal: '%s'", string(b))

	// can it be unmarshal back?
	j := make(map[string]interface{})
	err = json.Unmarshal(b, &j)
	if err != nil {
		t.Errorf("cannot unmarshal back: %v", err)
		return
	}
	b, _ = json.Marshal(&j)
	t.Logf("marshal back: %v", string(b))

	// just for reference
	// {
	// 	v := make(map[string]interface{})
	// 	json.Unmarshal([]byte(raw), &v)
	// 	b, _ := json.Marshal(&v)
	// 	t.Logf("official: %v", string(b))
	// }

	return
}

func TestMiscCharacters(t *testing.T) {
	s := "\"/\b\f\t\r\n<>&你好世界\\n"
	expected := "\"\\\"\\/\\b\\f\\t\\r\\n\\u003C\\u003E\\u0026\\u4F60\\u597D\\u4E16\\u754C\\\\n\""
	v := NewString(s)
	raw, err := v.MarshalString()
	if err != nil {
		t.Errorf("MarshalString() failed: %v", err)
		return
	}

	t.Logf("marshaled: '%s'", raw)
	if raw != expected {
		t.Errorf("marshal does not acted as expected <%s>", raw)
		t.Errorf("%s <-- raw", raw)
		t.Errorf("%s <-- expected", expected)
		return
	}
}

func TestUTF16(t *testing.T) {
	orig := "你👨‍👩‍👧‍👧你"

	v := NewObject()
	v.SetString(orig).At("string")

	data := struct {
		String string `json:"string"`
	}{}

	s := v.MustMarshalString()
	t.Logf("marshaled string '%s': '%s'", orig, s)
	if orig == s {
		t.Errorf("marshaled string should not equal!")
		return
	}

	b := v.MustMarshal()
	err := json.Unmarshal(b, &data)
	if err != nil {
		t.Errorf("unmarshal string '%s' failed: %v", orig, err)
		return
	}

	if data.String != orig {
		t.Errorf("unmarshaled string not expected! <%s>", string(b))
		return
	}
}

func TestPercentage(t *testing.T) {
	s := "%"
	expectedA := "\"\\u0025\""
	expectedB := "\"%\""
	v := NewString(s)
	raw, err := v.MarshalString()
	if err != nil {
		t.Errorf("MarshalString() failed: %v", err)
		return
	}

	t.Log("marshaled: '" + raw + "'")
	if raw != expectedA && raw != expectedB {
		t.Errorf("marshal does not acted as expected")
		return
	}
}

func TestMiscInt(t *testing.T) {
	var err error
	var checkCount int
	checkInt := func(i, expected int) {
		defer func() {
			checkCount++
		}()
		if err != nil {
			t.Errorf("%02d: Unexpected error: %v", checkCount, err)
			return
		}
		if i != expected {
			t.Errorf("%02d: i(%d) != %d", checkCount, i, expected)
			return
		}
	}

	raw := `[1,2,3,4,5,6,7]`
	v, err := UnmarshalString(raw)
	checkInt(0, 0)

	i, err := v.GetInt(uint(2))
	checkInt(i, 3)

	_, err = v.GetInt(int64(2))
	checkInt(i, 3)

	_, err = v.GetInt(uint64(2))
	checkInt(i, 3)

	_, err = v.GetInt(int32(2))
	checkInt(i, 3)

	_, err = v.GetInt(uint32(2))
	checkInt(i, 3)

	_, err = v.GetInt(int16(2))
	checkInt(i, 3)

	_, err = v.GetInt(uint16(2))
	checkInt(i, 3)

	_, err = v.GetInt(int8(2))
	checkInt(i, 3)

	_, err = v.GetInt(uint8(2))
	checkInt(i, 3)
}
