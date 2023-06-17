package jsonvalue

// ================ INSERT ================

// MustInserter is just like Inserter, but not returning sub-value or error.
type MustInserter interface {
	// After completes the following operation of Insert(). It inserts value AFTER specified position.
	//
	// The last parameter identifies the postion where a new JSON is inserted after, it should ba an interger, no matter signed or unsigned.
	// If the position is zero or positive interger, it tells the index of an array. If the position is negative, it tells the backward index of an array.
	//
	// For example, 0 represents the first, and -2 represents the second last.
	//
	// After 结束并完成 Insert() 函数的后续插入操作，表示插入到指定位置的前面。
	//
	// 在 Before 函数的最后一个参数指定了被插入的 JSON 数组的位置，这个参数应当是一个整型（有无符号类型均可）。
	// 如果这个值等于0或者正整数，那么它指定的是在 JSON 数组中的位置（从0开始）。如果这个值是负数，那么它指定的是 JSON 数组中从最后一个位置开始算起的位置。
	//
	// 举例说明：0 表示第一个位置，而 -2 表示倒数第二个位置。
	After(firstParam interface{}, otherParams ...interface{})

	// Before completes the following operation of Insert(). It inserts value BEFORE specified position.
	//
	// The last parameter identifies the postion where a new JSON is inserted after, it should ba an interger, no matter signed or unsigned.
	// If the position is zero or positive interger, it tells the index of an array. If the position is negative, it tells the backward index of an array.
	//
	// For example, 0 represents the first, and -2 represents the second last.
	//
	// Before 结束并完成 Insert() 函数的后续插入操作，表示插入到指定位置的后面。
	//
	// 在 Before 函数的最后一个参数指定了被插入的 JSON 数组的位置，这个参数应当是一个整型（有无符号类型均可）。
	// 如果这个值等于0或者正整数，那么它指定的是在 JSON 数组中的位置（从0开始）。如果这个值是负数，那么它指定的是 JSON 数组中从最后一个位置开始算起的位置。
	//
	// 举例说明：0 表示第一个位置，而 -2 表示倒数第二个位置。
	Before(firstParam interface{}, otherParams ...interface{})
}

type mInsert struct {
	inserter Inserter
}

// MustInsert is just like Insert, but not returning sub-value or error.
func (v *V) MustInsert(child any) MustInserter {
	ins := v.Insert(child)
	return &mInsert{
		inserter: ins,
	}
}

// MustInsertString is just like InsertString, but not returning sub-value or error.
func (v *V) MustInsertString(s string) MustInserter {
	return v.MustInsert(NewString(s))
}

// MustInsertBool is just like InsertBool, but not returning sub-value or error.
func (v *V) MustInsertBool(b bool) MustInserter {
	return v.MustInsert(NewBool(b))
}

// MustInsertInt is just like InsertInt, but not returning sub-value or error.
func (v *V) MustInsertInt(i int) MustInserter {
	return v.MustInsert(NewInt(i))
}

// MustInsertInt64 is just like InsertInt64, but not returning sub-value or error.
func (v *V) MustInsertInt64(i int64) MustInserter {
	return v.MustInsert(NewInt64(i))
}

// MustInsertInt32 is just like InsertInt32, but not returning sub-value or error.
func (v *V) MustInsertInt32(i int32) MustInserter {
	return v.MustInsert(NewInt32(i))
}

// MustInsertUint is just like InsertUint, but not returning sub-value or error.
func (v *V) MustInsertUint(u uint) MustInserter {
	return v.MustInsert(NewUint(u))
}

// MustInsertUint64 is just like InsertUint64, but not returning sub-value or error.
func (v *V) MustInsertUint64(u uint64) MustInserter {
	return v.MustInsert(NewUint64(u))
}

// MustInsertUint32 is just like InsertUint32, but not returning sub-value or error.
func (v *V) MustInsertUint32(u uint32) MustInserter {
	return v.MustInsert(NewUint32(u))
}

// MustInsertFloat64 is just like InsertFloat64, but not returning sub-value or error.
func (v *V) MustInsertFloat64(f float64) MustInserter {
	return v.MustInsert(NewFloat64(f))
}

// MustInsertFloat32 is just like InsertFloat32, but not returning sub-value or error.
func (v *V) MustInsertFloat32(f float32) MustInserter {
	return v.MustInsert(NewFloat32(f))
}

// MustInsertNull is just like InsertNull, but not returning sub-value or error.
func (v *V) MustInsertNull() MustInserter {
	return v.MustInsert(NewNull())
}

// MustInsertObject is just like InsertObject, but not returning sub-value or error.
func (v *V) MustInsertObject() MustInserter {
	return v.MustInsert(NewObject())
}

// MustInsertArray is just like InsertArray, but not returning sub-value or error.
func (v *V) MustInsertArray() MustInserter {
	return v.MustInsert(NewArray())
}

func (ins *mInsert) Before(firstParam any, otherParams ...any) {
	_, _ = ins.inserter.Before(firstParam, otherParams...)
}

func (ins *mInsert) After(firstParam any, otherParams ...any) {
	_, _ = ins.inserter.After(firstParam, otherParams...)
}

// ================ APPEND ================

// MustAppender is just like Appender, but not returning sub-value or error.
type MustAppender interface {
	// InTheBeginning completes the following operation of Append().
	//
	// InTheBeginning 函数将 Append 函数指定的 JSON 值，添加到参数指定的数组的最前端
	InTheBeginning(params ...interface{})

	// InTheEnd completes the following operation of Append().
	//
	// InTheEnd 函数将 Append 函数指定的 JSON 值，添加到参数指定的数组的最后面
	InTheEnd(params ...interface{})
}

type mAppender struct {
	appender Appender
}

// Append starts appending a child JSON value to a JSON array.
//
// Append 开始将一个 JSON 值添加到一个数组中。需结合 InTheEnd() 和 InTheBeginning() 函数使用。
func (v *V) MustAppend(child any) MustAppender {
	appd := v.Append(child)
	return &mAppender{
		appender: appd,
	}
}

// MustAppendString is just like AppendString, but not returning sub-value or error.
func (v *V) MustAppendString(s string) MustAppender {
	return v.MustAppend(NewString(s))
}

// MustAppendBytes is just like AppendBytes, but not returning sub-value or error.
func (v *V) MustAppendBytes(b []byte) MustAppender {
	return v.MustAppend(NewBytes(b))
}

// MustAppendBool is just like AppendBool, but not returning sub-value or error.
func (v *V) MustAppendBool(b bool) MustAppender {
	return v.MustAppend(NewBool(b))
}

// MustAppendInt is just like AppendInt, but not returning sub-value or error.
func (v *V) MustAppendInt(i int) MustAppender {
	return v.MustAppend(NewInt(i))
}

// MustAppendInt64 is just like AppendInt64, but not returning sub-value or error.
func (v *V) MustAppendInt64(i int64) MustAppender {
	return v.MustAppend(NewInt64(i))
}

// MustAppendInt32 is just like AppendInt32, but not returning sub-value or error.
func (v *V) MustAppendInt32(i int32) MustAppender {
	return v.MustAppend(NewInt32(i))
}

// MustAppendUint is just like AppendUint, but not returning sub-value or error.
func (v *V) MustAppendUint(u uint) MustAppender {
	return v.MustAppend(NewUint(u))
}

// MustAppendUint64 is just like AppendUint64, but not returning sub-value or error.
func (v *V) MustAppendUint64(u uint64) MustAppender {
	return v.MustAppend(NewUint64(u))
}

// MustAppendUint32 is just like AppendUint32, but not returning sub-value or error.
func (v *V) MustAppendUint32(u uint32) MustAppender {
	return v.MustAppend(NewUint32(u))
}

// MustAppendFloat64 is just like AppendFloat64, but not returning sub-value or error.
func (v *V) MustAppendFloat64(f float64) MustAppender {
	return v.MustAppend(NewFloat64(f))
}

// MustAppendFloat32 is just like AppendFloat32, but not returning sub-value or error.
func (v *V) MustAppendFloat32(f float32) MustAppender {
	return v.MustAppend(NewFloat32(f))
}

// MustAppendNull is just like AppendNull, but not returning sub-value or error.
func (v *V) MustAppendNull() MustAppender {
	return v.MustAppend(NewNull())
}

// MustAppendObject is just like AppendObject, but not returning sub-value or error.
func (v *V) MustAppendObject() MustAppender {
	return v.MustAppend(NewObject())
}

// MustAppendArray is just like AppendArray, but not returning sub-value or error.
func (v *V) MustAppendArray() MustAppender {
	return v.MustAppend(NewArray())
}

func (apd *mAppender) InTheBeginning(params ...any) {
	_, _ = apd.appender.InTheBeginning(params...)
}

func (apd *mAppender) InTheEnd(params ...any) {
	_, _ = apd.appender.InTheEnd(params...)
}

// ================ DELETE ================

// MustDelete is just like Delete, but not returning error.
func (v *V) MustDelete(firstParam any, otherParams ...any) {
	_ = v.Delete(firstParam, otherParams...)
}
