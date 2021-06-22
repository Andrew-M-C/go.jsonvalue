package jsonvalue

import (
	"encoding/json"
)

// Export convert jsonvalue to another type of parameter. The target parameter type should match the type of *V.
func (v *V) Export(dst interface{}) error {
	b, err := v.Marshal()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dst)
}

// Import convert json value from a marsalable parameter to *V
func Import(src interface{}) (*V, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	return Unmarshal(b)
}
