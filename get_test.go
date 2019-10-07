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

	return
}
