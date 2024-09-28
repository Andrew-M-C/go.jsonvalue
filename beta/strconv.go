package beta

import (
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// Deprecated: beta.Import is released, please use jsonvalue.Import.
func Import(src any) (*jsonvalue.V, error) {
	return jsonvalue.Import(src)
}
