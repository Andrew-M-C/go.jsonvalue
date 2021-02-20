package jsonvalue

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

func (v *V) delFromObjectChildren(key string) (exist bool) {
	_, exist = v.children.object[key]
	if exist {
		delete(v.children.object, key)
		v.delCaselessKey(key)
		return true
	}

	lowerKey := strings.ToLower(key)
	keys, exist := v.children.lowerCaseKeys[lowerKey]
	if !exist {
		return false
	}

	for actualKey := range keys {
		_, exist = v.children.object[actualKey]
		if exist {
			delete(v.children.object, actualKey)
			v.delCaselessKey(actualKey)
			return true
		}
	}

	return false
}

// Delete deletes specified JSON value. Forexample, parameters ("data", "list") identifies deleting value in data.list.
// While ("list", 1) means deleting 2nd (count from one) element from the "list" array.
//
// Delete 从 JSON 中删除参数指定的对象。比如参数 ("data", "list") 表示删除 data.list 值；参数 ("list", 1) 则表示删除 list
// 数组的第2（从1算起）个值。
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
		v.children.array.Remove(e)
		return nil
	}

	// else, this is an object value
	return fmt.Errorf("%v type does not supports Delete()", v.valueType)
}
