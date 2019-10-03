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
