package jsonvalue

// MustSetter is just like Setter, but not returning sub-value or error.
type MustSetter interface {
	// At completes the following operation of Set(). It defines position of value in Set() and return the new value set.
	//
	// The usage of At() is perhaps the most important. This function will recursively search for child value, and set the
	// new value specified by Set() or SetXxx() series functions. Please unfold and read the following examples, they are important.
	//
	// At 完成 Set() 函数的后续操作并设置相应的子成员。其参数指定了应该在哪个位置设置子成员，并且返回被设置的子成员对象。
	//
	// 该函数的用法恐怕是 jsonvalue 中最重要的内容了：该函数会按照给定的可变参数递归地一层一层查找 JSON 值的子成员，并且设置到指定的位置上。
	// 设置的逻辑说明起来比较抽象，请打开以下的例子以了解，这非常重要。
	At(firstParam interface{}, otherParams ...interface{})
}

type mSetter struct {
	setter Setter
}

// MustSet is just like Set, but not returning sub-value or error.
func (v *V) MustSet(child any) MustSetter {
	setter := v.Set(child)
	return &mSetter{
		setter: setter,
	}
}

// SetString is equivalent to Set(jsonvalue.NewString(s))
//
// SetString 等效于 Set(jsonvalue.NewString(s))
func (v *V) MustSetString(s string) MustSetter {
	return v.MustSet(NewString(s))
}

// SetBytes is equivalent to Set(NewString(base64.StdEncoding.EncodeToString(b)))
//
// SetBytes 等效于 Set(NewString(base64.StdEncoding.EncodeToString(b)))
func (v *V) MustSetBytes(b []byte) MustSetter {
	s := internal.b64.EncodeToString(b)
	return v.MustSetString(s)
}

// SetBool is equivalent to Set(jsonvalue.NewBool(b))
//
// SetBool 等效于 Set(jsonvalue.NewBool(b))
func (v *V) MustSetBool(b bool) MustSetter {
	return v.MustSet(NewBool(b))
}

// SetInt is equivalent to Set(jsonvalue.NewInt(b))
//
// SetInt 等效于 Set(jsonvalue.NewInt(b))
func (v *V) MustSetInt(i int) MustSetter {
	return v.MustSet(NewInt(i))
}

// SetInt64 is equivalent to Set(jsonvalue.NewInt64(b))
//
// SetInt64 等效于 Set(jsonvalue.NewInt64(b))
func (v *V) MustSetInt64(i int64) MustSetter {
	return v.MustSet(NewInt64(i))
}

// SetInt32 is equivalent to Set(jsonvalue.NewInt32(b))
//
// SetInt32 等效于 Set(jsonvalue.NewInt32(b))
func (v *V) MustSetInt32(i int32) MustSetter {
	return v.MustSet(NewInt32(i))
}

// SetUint is equivalent to Set(jsonvalue.NewUint(b))
//
// SetUint 等效于 Set(jsonvalue.NewUint(b))
func (v *V) MustSetUint(u uint) MustSetter {
	return v.MustSet(NewUint(u))
}

// SetUint64 is equivalent to Set(jsonvalue.NewUint64(b))
//
// SetUint64 is equivalent to Set(jsonvalue.NewUint64(b))
func (v *V) MustSetUint64(u uint64) MustSetter {
	return v.MustSet(NewUint64(u))
}

// SetUint32 is equivalent to Set(jsonvalue.NewUint32(b))
//
// SetUint32 等效于 Set(jsonvalue.NewUint32(b))
func (v *V) MustSetUint32(u uint32) MustSetter {
	return v.MustSet(NewUint32(u))
}

// SetFloat64 is equivalent to Set(jsonvalue.NewFloat64(b))
//
// SetFloat64 等效于 Set(jsonvalue.NewFloat64(b))
func (v *V) MustSetFloat64(f float64) MustSetter {
	return v.MustSet(NewFloat64(f))
}

// SetFloat32 is equivalent to Set(jsonvalue.NewFloat32(b))
//
// SetFloat32 等效于 Set(jsonvalue.NewFloat32(b))
func (v *V) MustSetFloat32(f float32) MustSetter {
	return v.MustSet(NewFloat32(f))
}

// SetNull is equivalent to Set(jsonvalue.NewNull())
//
// SetNull 等效于 Set(jsonvalue.NewNull())
func (v *V) MustSetNull() MustSetter {
	return v.MustSet(NewNull())
}

// SetObject is equivalent to Set(jsonvalue.NewObject())
//
// SetObject 等效于 Set(jsonvalue.NewObject())
func (v *V) MustSetObject() MustSetter {
	return v.MustSet(NewObject())
}

// SetArray is equivalent to Set(jsonvalue.NewArray())
//
// SetArray 等效于 Set(jsonvalue.NewArray())
func (v *V) MustSetArray() MustSetter {
	return v.MustSet(NewArray())
}

func (s *mSetter) At(firstParam any, otherParams ...any) {
	_, _ = s.setter.At(firstParam, otherParams...)
}
