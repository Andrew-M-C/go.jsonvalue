package beta

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// Import convert json value from a marsalable parameter to *V. This a experimental function.
//
// Import 将符合 encoding/json 的 struct 转为 *jsonvalue.V 类型。不经过 encoding/json，并且支持 Option.
func Import(src interface{}, opts ...jsonvalue.Option) (*jsonvalue.V, error) {
	v, fu, err := validateValAndReturnParser(reflect.ValueOf(src), ext{})
	if err != nil {
		return &jsonvalue.V{}, err
	}
	res, err := fu(v, ext{})
	if err != nil {
		return &jsonvalue.V{}, err
	}
	return res, nil
}

// parserFunc 处理对应 reflect.Value 的函数
type parserFunc func(v reflect.Value, ex ext) (*jsonvalue.V, error)

type ext struct {
	omitempty bool
	toString  bool
}

// validateValAndReturnParser 检查入参合法性并返回响应的处理函数
func validateValAndReturnParser(v reflect.Value, ex ext) (out reflect.Value, fu parserFunc, err error) {
	out = v

	switch v.Kind() {
	default:
		// 	fallthrough
		// case reflect.Invalid, reflect.Complex64, reflect.Complex128:
		// 	fallthrough
		// case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		err = fmt.Errorf("jsonvalue: unsupported type: %v", v.Type())

	case reflect.Bool:
		fu = parseBoolValue

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		fu = parseIntValue

	case reflect.Uintptr, reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		fu = parseUintValue

	case reflect.Float32, reflect.Float64:
		fu = parseFloatValue

	case reflect.Array:
		fu = parseArrayValue

	case reflect.Interface:
		return validateValAndReturnParser(v.Elem(), ex)

	case reflect.Map:
		switch v.Type().Key().Kind() {
		default:
			err = fmt.Errorf("unsupported key type for a map: %v", v.Type().Key())
		case reflect.String:
			fu = parseStringMapValue
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			fu = parseIntMapValue
		case reflect.Uintptr, reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			fu = parseUintMapValue
		}

	case reflect.Ptr:
		fu = parsePtrValue

	case reflect.Slice:
		if v.Type() == reflect.TypeOf([]byte{}) {
			fu = parseBytesValue
		} else if v.Type() == reflect.TypeOf(json.RawMessage{}) {
			fu = parseJSONRawMessageValue
		} else {
			fu = parseSliceValue
		}

	case reflect.String:
		fu = parseStringValue

	case reflect.Struct:
		fu = parseStructValue
	}

	return
}

func parseBoolValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	b := v.Bool()
	if !b && ex.omitempty {
		return nil, nil
	}
	if ex.toString {
		return jsonvalue.NewString(fmt.Sprint(b)), nil
	}
	return jsonvalue.NewBool(b), nil
}

func parseIntValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	i := v.Int()
	if i == 0 && ex.omitempty {
		return nil, nil
	}
	if ex.toString {
		return jsonvalue.NewString(strconv.FormatInt(i, 10)), nil
	}
	return jsonvalue.NewInt64(i), nil
}

func parseUintValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	u := v.Uint()
	if u == 0 && ex.omitempty {
		return nil, nil
	}
	if ex.toString {
		return jsonvalue.NewString(strconv.FormatUint(u, 10)), nil
	}
	return jsonvalue.NewUint64(u), nil
}

func parseFloatValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	f := v.Float()
	if f == 0.0 && ex.omitempty {
		return nil, nil
	}
	if ex.toString {
		return jsonvalue.NewString(strconv.FormatFloat(f, 'f', -1, 64)), nil
	}
	return jsonvalue.NewFloat64(f), nil
}

func parseArrayValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	res := jsonvalue.NewArray()
	le := v.Len()

	for i := 0; i < le; i++ {
		vv := v.Index(i)
		vv, fu, err := validateValAndReturnParser(vv, ext{})
		if err != nil {
			return nil, err
		}
		child, err := fu(vv, ext{})
		if err != nil {
			return nil, err
		}
		res.Append(child).InTheEnd()
	}

	return res, nil
}

func parseMapValue(v reflect.Value, ex ext, keyFunc func(key reflect.Value) string) (*jsonvalue.V, error) {
	if v.IsNil() {
		return parseNullValue(v, ex)
	}

	keys := v.MapKeys()
	if len(keys) == 0 {
		if ex.omitempty {
			return nil, nil
		}
		return jsonvalue.NewObject(), nil
	}

	res := jsonvalue.NewObject()

	for _, kk := range keys {
		vv := v.MapIndex(kk)
		vv, fu, err := validateValAndReturnParser(vv, ext{})
		if err != nil {
			return res, err
		}
		child, err := fu(vv, ext{})
		if err != nil {
			return res, err
		}
		res.Set(child).At(keyFunc(kk))
	}

	return res, nil
}

func parseStringMapValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	return parseMapValue(v, ex, func(k reflect.Value) string {
		return k.String()
	})
}

func parseIntMapValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	return parseMapValue(v, ex, func(k reflect.Value) string {
		return strconv.FormatInt(k.Int(), 10)
	})
}

func parseUintMapValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	return parseMapValue(v, ex, func(k reflect.Value) string {
		return strconv.FormatUint(k.Uint(), 10)
	})
}

func parsePtrValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	if v.IsNil() {
		return parseNullValue(v, ex)
	}

	v, fu, err := validateValAndReturnParser(v.Elem(), ex)
	if err != nil {
		return nil, err
	}

	return fu(v, ex)
}

func parseSliceValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	if v.IsNil() || v.Len() == 0 {
		if ex.omitempty {
			return nil, nil
		}
		return jsonvalue.NewArray(), nil
	}

	return parseArrayValue(v, ex)
}

func parseBytesValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	b := v.Interface().([]byte)
	if len(b) == 0 && ex.omitempty {
		return nil, nil
	}

	return jsonvalue.NewBytes(b), nil
}

func parseJSONRawMessageValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	raw := v.Interface().(json.RawMessage)
	if len(raw) == 0 && ex.omitempty {
		return nil, nil
	}

	return jsonvalue.Unmarshal(raw)
}

func parseStringValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	str := v.String()
	if str == "" && ex.omitempty {
		return nil, nil
	}

	return jsonvalue.NewString(str), nil
}

func parseStructValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	t := v.Type()
	numField := t.NumField()

	res := jsonvalue.NewObject()

	for i := 0; i < numField; i++ {
		vv := v.Field(i)
		tt := t.Field(i)

		kv, err := parseStructFieldValue(vv, tt)
		if err != nil {
			return nil, err
		}

		for k, v := range kv {
			res.Set(v).At(k)
		}
	}

	return res, nil
}

func parseNullValue(v reflect.Value, ex ext) (*jsonvalue.V, error) {
	if ex.omitempty {
		return nil, nil
	}
	return jsonvalue.NewNull(), nil
}

func parseStructFieldValue(fv reflect.Value, ft reflect.StructField) (m map[string]*jsonvalue.V, err error) {
	m = map[string]*jsonvalue.V{}

	if ft.Anonymous {
		numField := fv.NumField()
		for i := 0; i < numField; i++ {
			ffv := fv.Field(i)
			fft := ft.Type.Field(i)

			mm, err := parseStructFieldValue(ffv, fft)
			if err != nil {
				return nil, err
			}
			for k, v := range mm {
				m[k] = v
			}
		}
		return m, nil
	}

	if !fv.CanInterface() {
		return
	}

	fieldName, ex := readFieldTag(ft, "json")
	if fieldName == "-" {
		return
	}

	fv, fu, err := validateValAndReturnParser(fv, ex)
	if err != nil {
		return m, fmt.Errorf("parsing field '%s' error: %w", fieldName, err)
	}

	child, err := fu(fv, ex)
	if err != nil {
		return m, fmt.Errorf("parsing field '%s' error: %w", fieldName, err)
	}
	if child != nil {
		m[fieldName] = child
	}

	return m, nil
}

func readFieldTag(ft reflect.StructField, name string) (field string, ex ext) {
	tg := ft.Tag.Get(name)

	if tg == "" {
		return ft.Name, ext{}
	}

	parts := strings.Split(tg, ",")
	for i, s := range parts {
		parts[i] = strings.TrimSpace(s)
		if i > 0 {
			if s == "omitempty" {
				ex.omitempty = true
			} else if s == "string" {
				ex.toString = true
			}
		}
	}

	field = parts[0]
	if field == "" {
		field = ft.Name
	}
	return
}
