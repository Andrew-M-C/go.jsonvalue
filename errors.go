package jsonvalue

// Error is equavilent to string and used to create some error constants in this package.
// Error constants: http://godoc.org/github.com/Andrew-M-C/go.jsonvalue/#pkg-constants
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// ErrNilParameter identifies input paremeter is nil
	//
	// ErrNilParameter 表示参数为空
	ErrNilParameter = Error("nil parameter")

	// ErrValueUninitialized identifies that a V object is not initialized
	//
	// ErrValueUninitialized 表示当前的 jsonvalue 实例未初始化
	ErrValueUninitialized = Error("jsonvalue instance is not initialized")

	// ErrRawBytesUnrecignized identifies all unexpected raw bytes
	//
	// ErrRawBytesUnrecignized 表示无法识别的序列文本
	ErrRawBytesUnrecignized = Error("unrecognized raw text")

	// ErrNotValidBoolValue shows that a value starts with 't' or 'f' is not eventually a bool value
	//
	// ErrNotValidBoolValue 表示当前值不是一个合法的布尔值
	ErrNotValidBoolValue = Error("not a valid bool object")

	// ErrNotValidNulllValue shows that a value starts with 'n' is not eventually a bool value
	//
	// ErrNotValidNulllValue 表示当前不是一个 null 值类型的 JSON
	ErrNotValidNulllValue = Error("not a valid null object")

	// ErrOutOfRange identifies that given position for a JSON array is out of range
	//
	// ErrOutOfRange 表示请求数组成员超出数组范围
	ErrOutOfRange = Error("out of range")

	// ErrNotFound shows that given target is not found in Delete()
	//
	// ErrNotFound 表示目标无法找到
	ErrNotFound = Error("target not found")

	// ErrTypeNotMatch shows that value type is not same as GetXxx()
	//
	// ErrTypeNotMatch 表示指定的对象不匹配
	ErrTypeNotMatch = Error("not match given type")

	// ErrNotArrayValue shows that operation target value is not an array
	//
	// ErrNotArrayValue 表示当前不是一个数组类型 JSON
	ErrNotArrayValue = Error("not an array typed value")

	// ErrIllegalString shows that it is not a legal JSON string typed value
	//
	// ErrIllegalString 表示字符串不合法
	ErrIllegalString = Error("illegal string")
)
