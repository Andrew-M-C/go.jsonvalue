package jsonvalue

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

func (v *V) delFromObjectChildren(caseless bool, key string) (exist bool) {
	_, exist = v.children.object[key]
	if exist {
		delete(v.children.object, key)
		v.delCaselessKey(key)
		return true
	}

	if !caseless {
		return false
	}

	v.initCaselessStorage()

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
	return v.delete(false, firstParam, otherParams...)
}

func (v *V) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	paramCount := len(otherParams)
	if paramCount == 0 {
		return v.deleteInCurrValue(caseless, firstParam)
	}

	child, err := v.get(caseless, firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	// if child == nil {
	// 	return ErrNotFound
	// }

	return child.delete(caseless, otherParams[paramCount-1])
}

func (v *V) deleteInCurrValue(caseless bool, param interface{}) error {
	if v.valueType == jsonparser.Object {
		// string expected
		key, err := intfToString(param)
		if err != nil {
			return err
		}

		if exist := v.delFromObjectChildren(caseless, key); !exist {
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

		pos = v.posAtIndexForRead(pos)
		if pos < 0 {
			return ErrOutOfRange
		}
		v.deleteInArr(pos)
		return nil
	}

	// else, this is an object value
	return fmt.Errorf("%v type does not supports Delete()", v.valueType)
}

func (v *V) deleteInArr(pos int) {
	le := len(v.children.array)
	v.children.array[pos] = nil
	copy(v.children.array[pos:], v.children.array[pos+1:])
	v.children.array = v.children.array[:le-1]
}
