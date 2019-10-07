package jsonvalue

import (
	"fmt"
)

var (
	// ErrNilParameter identifies input paremeter is nil
	ErrNilParameter = fmt.Errorf("nil parameter")
	// ErrValueUninitialized identifies that a V object is not initialized
	ErrValueUninitialized = fmt.Errorf("jsonvalue instance is not initialized")
	// ErrRawBytesUnrecignized identifies all unexpected raw bytes
	ErrRawBytesUnrecignized = fmt.Errorf("unrecognized raw text")
	// ErrNotValidBoolValue shows that a value starts with 't' or 'f' is not eventually a bool value
	ErrNotValidBoolValue = fmt.Errorf("not a valid bool object")
	// ErrNotValidNulllValue shows that a value starts with 'n' is not eventually a bool value
	ErrNotValidNulllValue = fmt.Errorf("not a valid null object")
	// ErrOutOfRange identifies that given position for a JSON array is out of range
	ErrOutOfRange = fmt.Errorf("out of range")
	// ErrNotFound shows that given target is not found in Delete()
	ErrNotFound = fmt.Errorf("target not found")
	// ErrTypeNotMatch shows that value type is not same as GetXxx()
	ErrTypeNotMatch = fmt.Errorf("not match given type")
	// ErrNotArrayValue shows that operation target value is not an array
	ErrNotArrayValue = fmt.Errorf("not an array typed value")
)
