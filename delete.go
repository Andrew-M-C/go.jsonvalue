package jsonvalue

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

func (v *V) delFromObjectChildren(key string) (exist bool) {
	_, exist = v.objectChildren[key]
	if exist {
		delete(v.objectChildren, key)
		v.delCaselessKey(key)
		return true
	}

	lowerKey := strings.ToLower(key)
	keys, exist := v.lowerCaseKeys[lowerKey]
	if !exist {
		return false
	}

	for actualKey := range keys {
		_, exist = v.objectChildren[actualKey]
		if exist {
			delete(v.objectChildren, actualKey)
			v.delCaselessKey(actualKey)
			return true
		}
	}

	return false
}

// Delete deletes specified JSON value
func (v *V) Delete(firstParam interface{}, otherParams ...interface{}) error {
	paramCount := len(otherParams)
	if 0 == paramCount {
		return v.deleteInCurrValue(firstParam)
	}

	child, err := v.Get(firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	// if nil == child {
	// 	return ErrNotFound
	// }

	return child.Delete(otherParams[paramCount-1])
}

func (v *V) deleteInCurrValue(param interface{}) error {
	if v.valueType == jsonparser.Object {
		// string expected
		key, err := intfToString(param)
		if err != nil {
			return err
		}

		if exist := v.delFromObjectChildren(key); !exist {
			return ErrNotFound
		}
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
