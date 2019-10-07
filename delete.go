package jsonvalue

import (
	"fmt"

	"github.com/buger/jsonparser"
)

// Delete deletes specified JSON value
func (v *V) Delete(firstParam interface{}, otherParams ...interface{}) error {
	if 0 == len(otherParams) {
		return v.deleteInCurrValue(firstParam)
	}

	child, err := v.Get(firstParam)
	if err != nil {
		return err
	}
	if nil == child {
		return ErrNotFound
	}

	return child.Delete(otherParams[0], otherParams[1:])
}

func (v *V) deleteInCurrValue(param interface{}) error {
	if v.valueType == jsonparser.Object {
		// string expected
		key, err := intfToString(param)
		if err != nil {
			return err
		}

		if _, exist := v.objectChildren[key]; false == exist {
			return ErrNotFound
		}
		delete(v.objectChildren, key)
		return nil
	}

	if v.valueType == jsonparser.Array {
		// interger expected
		pos, err := intfToInt(param)
		if err != nil {
			return err
		}

		e := v.elementAtIndex(pos)
		if nil == e {
			return ErrOutOfRange
		}
		v.arrayChildren.Remove(e)
		return nil
	}

	// else, this is an object value
	return fmt.Errorf("%v type does not supports Delete()", v.valueType)
}
