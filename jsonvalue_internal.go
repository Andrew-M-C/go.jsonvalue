package jsonvalue

import "encoding/base64"

var internal = struct {
	b64 *base64.Encoding

	defaultMarshalOption *Opt
}{
	b64: base64.StdEncoding,

	defaultMarshalOption: emptyOptions(),
}
