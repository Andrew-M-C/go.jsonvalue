package jsonvalue

import (
	"os"
	"testing"
)

func check(t *testing.T, err error, function string, b bool) {
	if err != nil {
		t.Errorf("%s() error: %v", function, err)
		os.Exit(1)
	}
	if false == b {
		t.Errorf("%s() failed", function)
		os.Exit(1)
	}
	return
}

func TestGet(t *testing.T) {
	full := `{"data":{"message":["hello","world",true,null],"author":"Andrew","year":2019,"YYYY.MM":2019.12,"negative":-1234}}`
	o, err := UnmarshalString(full)
	if err != nil {
		t.Errorf("UnmarshalString failed: %v", err)
		return
	}
	b, _ := o.Marshal()
	t.Logf("unmarshal back: %s", string(b))

	{
		s, err := o.GetString("data", "author")
		check(t, err, "GetString", s == "Andrew")
	}
	{
		i, err := o.GetInt("data", "year")
		check(t, err, "GetInt", i == 2019)
	}
	{
		i, err := o.GetUint("data", "year")
		check(t, err, "GetUint", i == 2019)
	}
	{
		i, err := o.GetInt64("data", "year")
		check(t, err, "GetInt64", i == 2019)
	}
	{
		i, err := o.GetUint64("data", "year")
		check(t, err, "GetUint64", i == 2019)
	}
	{
		i, err := o.GetInt32("data", "negative")
		check(t, err, "GetInt32", i == -1234)
	}
	{
		i, err := o.GetInt32("data", "negATive") // caseless
		check(t, err, "GetInt32_caseless", i == -1234)
	}
	{
		i, err := o.GetUint32("data", "year")
		check(t, err, "GetUint64", i == 2019)
	}
	{
		f, err := o.GetFloat64("data", "YYYY.MM")
		check(t, err, "GetFloat64", f == 2019.12)
	}
	{
		f, err := o.GetFloat32("data", "YYYY.MM")
		check(t, err, "GetFloat32", f == 2019.12)
	}
	{
		err := o.GetNull("data", "message", -1)
		check(t, err, "GetNull", true)
	}
	{
		b, err := o.GetBool("data", "message", 2)
		check(t, err, "GetNull", b == true)
	}
	{
		s, err := o.GetString("data", "message", 0)
		check(t, err, "GetNull", s == "hello")
	}
	{
		s, err := o.GetString("data", "message", -3)
		check(t, err, "GetString", s == "world")
	}
	{
		l := o.Len()
		check(t, nil, "Len", l == 1)

		v, _ := o.Get("data", "message")
		l = v.Len()
		check(t, nil, "Len", l == 4)

		v, _ = o.Get("data", "author")
		l = v.Len()
		check(t, nil, "Len", l == 0)
	}
	{
		v, err := o.GetObject("data")
		check(t, err, "GetObject", v.IsObject())
	}

	return
}

func TestMiscError(t *testing.T) {
	var checkCount int
	shouldError := func(err error) {
		defer func() {
			checkCount++
		}()
		if err == nil {
			t.Errorf("%02d - error expected but not caught", checkCount)
			return
		}
		t.Logf("expected error string: %v", err)
		return
	}

	{
		var err error
		raw := `{"array":[0,1,2,3],"string":"hello, world","number":1234.12345}`
		v, _ := UnmarshalString(raw)

		// param error
		_, err = v.GetInt("array", true)
		shouldError(err)
		_, err = v.GetString(true)
		shouldError(err)

		// out of range
		_, err = v.Get("array", 100)
		shouldError(err)

		// not support
		_, err = v.Get("string", "hello")
		shouldError(err)

		// GetString
		_, err = v.GetString("number")
		shouldError(err)
		_, err = v.GetString("not exist")
		shouldError(err)

		// GetInt... and GetUint..
		_, err = v.GetInt("string")
		shouldError(err)
		_, err = v.GetUint("string")
		shouldError(err)
		_, err = v.GetInt64("string")
		shouldError(err)
		_, err = v.GetUint64("string")
		shouldError(err)
		_, err = v.GetInt32("string")
		shouldError(err)
		_, err = v.GetUint32("string")
		shouldError(err)
		_, err = v.GetFloat64("string")
		shouldError(err)
		_, err = v.GetFloat32("string")
		shouldError(err)

		// number not exist
		_, err = v.GetString("not exist")
		shouldError(err)
		_, err = v.GetInt("not exist")
		shouldError(err)
		_, err = v.GetUint("not exist")
		shouldError(err)
		_, err = v.GetInt64("not exist")
		shouldError(err)
		_, err = v.GetUint64("not exist")
		shouldError(err)
		_, err = v.GetInt32("not exist")
		shouldError(err)
		_, err = v.GetUint32("not exist")
		shouldError(err)
		_, err = v.GetFloat64("not exist")
		shouldError(err)
		_, err = v.GetFloat32("not exist")
		shouldError(err)

		// GetObject and GetArray
		_, err = v.GetObject("string")
		shouldError(err)
		_, err = v.GetArray("string")
		shouldError(err)
		_, err = v.GetObject("not exist")
		shouldError(err)
		_, err = v.GetArray("not exist")
		shouldError(err)

		// GetBool and GetNull
		_, err = v.GetBool("string")
		shouldError(err)
		err = v.GetNull("string")
		shouldError(err)
		_, err = v.GetBool("not exist")
		shouldError(err)
		err = v.GetNull("not exist")
		shouldError(err)
	}

}
