package jsonvalue

import (
	"encoding"
	"encoding/json"
)

var (
	_ json.Marshaler   = (*V)(nil)
	_ json.Unmarshaler = (*V)(nil)

	_ encoding.BinaryMarshaler   = (*V)(nil)
	_ encoding.BinaryUnmarshaler = (*V)(nil)
)

// MarshalJSON implements json.Marshaler
func (v *V) MarshalJSON() ([]byte, error) {
	return v.Marshal()
}

// UnmarshalJSON implements json.Unmarshaler
func (v *V) UnmarshalJSON(b []byte) error {
	res, err := Unmarshal(b)
	if err != nil {
		return err
	}
	*v = *res
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler
func (v *V) MarshalBinary() ([]byte, error) {
	return v.Marshal()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (v *V) UnmarshalBinary(b []byte) error {
	res, err := Unmarshal(b)
	if err != nil {
		return err
	}
	*v = *res
	return nil
}
