package jsonvalue

import (
	"encoding"
	"encoding/base64"
	"encoding/json"
	"reflect"
)

var internal = struct {
	base64 *base64.Encoding

	defaultMarshalOption *Opt

	types struct {
		JSONMarshaler reflect.Type
		TextMarshaler reflect.Type
	}
}{}

func init() {
	internal.base64 = base64.StdEncoding
	internal.defaultMarshalOption = emptyOptions()

	internal.types.JSONMarshaler = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	internal.types.TextMarshaler = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
}

type pool interface {
	Get() *V
}

type globalPool struct{}

func (globalPool) Get() *V {
	return &V{}
}
