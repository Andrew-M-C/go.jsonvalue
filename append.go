package jsonvalue

import (
	"github.com/buger/jsonparser"
)

type append struct {
	v *V
	c *V // child
}

// Append starts setting a child JSON value
func (v *V) Append(child *V) *append {
	if nil == child {
		child = NewNull()
	}
	return &append{
		v: v,
		c: child,
	}
}

// AppendString is equivalent to Append(jsonvalue.NewString(s))
func (v *V) AppendString(s string) *append {
	return v.Append(NewString(s))
}

// AppendBool is equivalent to Append(jsonvalue.NewBool(b))
func (v *V) AppendBool(b bool) *append {
	return v.Append(NewBool(b))
}

// AppendInt is equivalent to Append(jsonvalue.NewInt(b))
func (v *V) AppendInt(i int) *append {
	return v.Append(NewInt(i))
}

// AppendInt64 is equivalent to Append(jsonvalue.NewInt64(b))
func (v *V) AppendInt64(i int64) *append {
	return v.Append(NewInt64(i))
}

// AppendInt32 is equivalent to Append(jsonvalue.NewInt32(b))
func (v *V) AppendInt32(i int32) *append {
	return v.Append(NewInt32(i))
}

// AppendUint is equivalent to Append(jsonvalue.NewUint(b))
func (v *V) AppendUint(u uint) *append {
	return v.Append(NewUint(u))
}

// AppendUint64 is equivalent to Append(jsonvalue.NewUint64(b))
func (v *V) AppendUint64(u uint64) *append {
	return v.Append(NewUint64(u))
}

// AppendUint32 is equivalent to Append(jsonvalue.NewUint32(b))
func (v *V) AppendUint32(u uint32) *append {
	return v.Append(NewUint32(u))
}

// AppendFloat64 is equivalent to Append(jsonvalue.NewFloat64(b))
func (v *V) AppendFloat64(f float64, prec int) *append {
	return v.Append(NewFloat64(f, prec))
}

// AppendFloat32 is equivalent to Append(jsonvalue.NewFloat32(b))
func (v *V) AppendFloat32(f float32, prec int) *append {
	return v.Append(NewFloat32(f, prec))
}

// AppendNull is equivalent to Append(jsonvalue.NewNull())
func (v *V) AppendNull() *append {
	return v.Append(NewNull())
}

// AppendObject is equivalent to Append(jsonvalue.NewObject())
func (v *V) AppendObject() *append {
	return v.Append(NewObject())
}

// AppendArray is equivalent to Append(jsonvalue.NewArray())
func (v *V) AppendArray() *append {
	return v.Append(NewArray())
}

// InTheBeginning completes the following operation of Append().
func (apd *append) InTheBeginning(params ...interface{}) (*V, error) {
	v := apd.v
	c := apd.c
	if v.valueType == jsonparser.NotExist {
		return nil, ErrValueUninitialized
	}

	// this is the last iteration
	paramCount := len(params)
	if 0 == paramCount {
		if v.valueType != jsonparser.Array {
			return nil, ErrNotArrayValue
		}

		v.arrayChildren.PushBack(c)
		return c, nil
	}

	// this is not the last iterarion
	shouldSet := false
	child, err := v.GetArray(params[0], params[1:paramCount]...)
	if err != nil {
		if err == ErrNotFound {
			shouldSet = true
			child = NewArray()
		} else {
			return nil, err
		}
	}

	child.arrayChildren.PushFront(c)
	if shouldSet {
		v.Set(child).At(params[0], params[1:paramCount]...)
	}
	return c, nil
}

// InTheEnd completes the following operation of Append().
func (apd *append) InTheEnd(params ...interface{}) (*V, error) {
	v := apd.v
	c := apd.c
	if v.valueType == jsonparser.NotExist {
		return nil, ErrValueUninitialized
	}

	// this is the last iteration
	paramCount := len(params)
	if 0 == paramCount {
		if v.valueType != jsonparser.Array {
			return nil, ErrNotArrayValue
		}

		v.arrayChildren.PushBack(c)
		return c, nil
	}

	// this is not the last iterarion
	shouldSet := false
	child, err := v.GetArray(params[0], params[1:paramCount]...)
	if err != nil {
		if err == ErrNotFound {
			shouldSet = true
			child = NewArray()
		} else {
			return nil, err
		}
	}

	child.arrayChildren.PushBack(c)
	if shouldSet {
		v.Set(child).At(params[0], params[1:paramCount]...)
	}
	return c, nil
}
