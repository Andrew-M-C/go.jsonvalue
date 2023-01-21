package jsonvalue

import (
	"encoding/base64"
)

var internal = struct {
	b64 *base64.Encoding

	defaultMarshalOption *Opt
}{}

func init() {
	internal.b64 = base64.StdEncoding
	internal.defaultMarshalOption = emptyOptions()
}
