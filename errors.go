package jsonvalue

// Error is equavilent to string and used to create some error constants in this package.
// Error constants: http://godoc.org/github.com/Andrew-M-C/go.jsonvalue/#pkg-constants
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// ErrNilParameter identifies input paremeter is nil
	ErrNilParameter = Error("nil parameter")
	// ErrValueUninitialized identifies that a V object is not initialized
	ErrValueUninitialized = Error("jsonvalue instance is not initialized")
	// ErrRawBytesUnrecignized identifies all unexpected raw bytes
	ErrRawBytesUnrecignized = Error("unrecognized raw text")
	// ErrNotValidBoolValue shows that a value starts with 't' or 'f' is not eventually a bool value
	ErrNotValidBoolValue = Error("not a valid bool object")
	// ErrNotValidNulllValue shows that a value starts with 'n' is not eventually a bool value
	ErrNotValidNulllValue = Error("not a valid null object")
	// ErrOutOfRange identifies that given position for a JSON array is out of range
	ErrOutOfRange = Error("out of range")
	// ErrNotFound shows that given target is not found in Delete()
	ErrNotFound = Error("target not found")
	// ErrTypeNotMatch shows that value type is not same as GetXxx()
	ErrTypeNotMatch = Error("not match given type")
	// ErrNotArrayValue shows that operation target value is not an array
	ErrNotArrayValue = Error("not an array typed value")
)
