package jsonvalue

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Export convert jsonvalue to another type of parameter. The target parameter type should match the type of *V.
//
// Export 将 *V 转到符合原生 encoding/json 的一个 struct 中。
func (v *V) Export(dst any) error {
	b, err := v.Marshal()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dst)
}

// Import convert json value from a marsalable parameter to *V. This a experimental function.
//
// Import 将符合 encoding/json 的 struct 转为 *V 类型。不经过 encoding/json，并且支持 Option.
func Import(src any, opts ...Option) (*V, error) {
	opt := combineOptions(opts)
	ext := ext{}
	ext.ignoreOmitempty = opt.ignoreJsonOmitempty
	v, fu, err := validateValAndReturnParser(reflect.ValueOf(src), ext)
	if err != nil {
		return &V{}, err
	}
	res, err := fu(v, ext)
	if err != nil {
		return &V{}, err
	}
	return res, nil
}

// parserFunc 处理对应 reflect.Value 的函数
type parserFunc func(v reflect.Value, ex ext) (*V, error)

type ext struct {
	// standard encoding/json tag
	omitempty bool
	toString  bool

	// extended jsonvalue options
	ignoreOmitempty bool
}

func (e ext) shouldOmitEmpty() bool {
	return e.omitempty && !e.ignoreOmitempty
}

// validateValAndReturnParser 检查入参合法性并返回相应的处理函数
func validateValAndReturnParser(v reflect.Value, ex ext) (out reflect.Value, fu parserFunc, err error) {
	out = v

	switch v.Kind() {
	default:
		// 	fallthrough
		// case reflect.Invalid, reflect.Complex64, reflect.Complex128:
		// 	fallthrough
		// case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		err = fmt.Errorf("jsonvalue: unsupported type: %v", v.Type())

	case reflect.Invalid:
		fu = parseInvalidValue

	case reflect.Bool:
		fu = parseBoolValue

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		fu = parseIntValue

	case reflect.Uintptr, reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		fu = parseUintValue

	case reflect.Float32:
		fu = parseFloat32Value

	case reflect.Float64:
		fu = parseFloat64Value

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

func parseInvalidValue(_ reflect.Value, ex ext) (*V, error) {
	if ex.shouldOmitEmpty() {
		return nil, nil
	}
	return NewNull(), nil
}

func parseBoolValue(v reflect.Value, ex ext) (*V, error) {
	b := v.Bool()
	if !b && ex.shouldOmitEmpty() {
		return nil, nil
	}
	if ex.toString {
		return NewString(fmt.Sprint(b)), nil
	}
	return NewBool(b), nil
}

func parseIntValue(v reflect.Value, ex ext) (*V, error) {
	i := v.Int()
	if i == 0 && ex.shouldOmitEmpty() {
		return nil, nil
	}
	if ex.toString {
		return NewString(strconv.FormatInt(i, 10)), nil
	}
	return NewInt64(i), nil
}

func parseUintValue(v reflect.Value, ex ext) (*V, error) {
	u := v.Uint()
	if u == 0 && ex.shouldOmitEmpty() {
		return nil, nil
	}
	if ex.toString {
		return NewString(strconv.FormatUint(u, 10)), nil
	}
	return NewUint64(u), nil
}

func parseFloat64Value(v reflect.Value, ex ext) (*V, error) {
	f := v.Float()
	if f == 0.0 && ex.shouldOmitEmpty() {
		return nil, nil
	}
	if ex.toString {
		f64 := NewFloat64(f)
		return NewString(f64.MustMarshalString()), nil
	}
	return NewFloat64(f), nil
}

func parseFloat32Value(v reflect.Value, ex ext) (*V, error) {
	f := v.Float()
	if f == 0.0 && ex.shouldOmitEmpty() {
		return nil, nil
	}
	if ex.toString {
		f32 := NewFloat32(float32(f))
		return NewString(f32.MustMarshalString()), nil
	}
	return NewFloat32(float32(f)), nil
}

func parseArrayValue(v reflect.Value, ex ext) (*V, error) {
	res := NewArray()
	le := v.Len()

	for i := 0; i < le; i++ {
		vv := v.Index(i)
		vv, fu, err := validateValAndReturnParser(vv, ex)
		if err != nil {
			return nil, err
		}
		child, err := fu(vv, ex)
		if err != nil {
			return nil, err
		}
		res.Append(child).InTheEnd()
	}

	return res, nil
}

func parseMapValue(v reflect.Value, ex ext, keyFunc func(key reflect.Value) string) (*V, error) {
	if v.IsNil() {
		return parseNullValue(v, ex)
	}

	keys := v.MapKeys()
	if len(keys) == 0 {
		if ex.shouldOmitEmpty() {
			return nil, nil
		}
		return NewObject(), nil
	}

	res := NewObject()

	for _, kk := range keys {
		vv := v.MapIndex(kk)
		vv, fu, err := validateValAndReturnParser(vv, ex)
		if err != nil {
			return res, err
		}
		child, err := fu(vv, ex)
		if err != nil {
			return res, err
		}
		res.Set(child).At(keyFunc(kk))
	}

	return res, nil
}

func parseStringMapValue(v reflect.Value, ex ext) (*V, error) {
	return parseMapValue(v, ex, func(k reflect.Value) string {
		return k.String()
	})
}

func parseIntMapValue(v reflect.Value, ex ext) (*V, error) {
	return parseMapValue(v, ex, func(k reflect.Value) string {
		return strconv.FormatInt(k.Int(), 10)
	})
}

func parseUintMapValue(v reflect.Value, ex ext) (*V, error) {
	return parseMapValue(v, ex, func(k reflect.Value) string {
		return strconv.FormatUint(k.Uint(), 10)
	})
}

func parsePtrValue(v reflect.Value, ex ext) (*V, error) {
	if v.IsNil() {
		return parseNullValue(v, ex)
	}

	v, fu, err := validateValAndReturnParser(v.Elem(), ex)
	if err != nil {
		return nil, err
	}

	return fu(v, ex)
}

func parseSliceValue(v reflect.Value, ex ext) (*V, error) {
	if v.IsNil() || v.Len() == 0 {
		if ex.shouldOmitEmpty() {
			return nil, nil
		}
		return NewArray(), nil
	}

	return parseArrayValue(v, ex)
}

func parseBytesValue(v reflect.Value, ex ext) (*V, error) {
	b := v.Interface().([]byte)
	if len(b) == 0 && ex.shouldOmitEmpty() {
		return nil, nil
	}

	return NewBytes(b), nil
}

func parseJSONRawMessageValue(v reflect.Value, ex ext) (*V, error) {
	raw := v.Interface().(json.RawMessage)
	if len(raw) == 0 && ex.shouldOmitEmpty() {
		return nil, nil
	}

	return Unmarshal(raw)
}

func parseStringValue(v reflect.Value, ex ext) (*V, error) {
	str := v.String()
	if str == "" && ex.shouldOmitEmpty() {
		return nil, nil
	}

	return NewString(str), nil
}

func parseStructValue(v reflect.Value, ex ext) (*V, error) {
	t := v.Type()
	numField := t.NumField()

	res := NewObject()

	for i := 0; i < numField; i++ {
		vv := v.Field(i)
		tt := t.Field(i)

		kv, err := parseStructFieldValue(vv, tt, ex)
		if err != nil {
			return nil, err
		}

		for k, v := range kv {
			res.Set(v).At(k)
		}
	}

	return res, nil
}

func parseNullValue(v reflect.Value, ex ext) (*V, error) {
	if ex.shouldOmitEmpty() {
		return nil, nil
	}
	return NewNull(), nil
}

func parseStructFieldValue(fv reflect.Value, ft reflect.StructField, parentEx ext) (m map[string]*V, err error) {
	m = map[string]*V{}

	if ft.Anonymous {
		numField := fv.NumField()
		for i := 0; i < numField; i++ {
			ffv := fv.Field(i)
			fft := ft.Type.Field(i)

			mm, err := parseStructFieldValue(ffv, fft, parentEx)
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

	fieldName, ex := readFieldTag(ft, "json", parentEx)
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

func readFieldTag(ft reflect.StructField, name string, parentEx ext) (field string, ex ext) {
	tg := ft.Tag.Get(name)

	if tg == "" {
		return ft.Name, ext{
			ignoreOmitempty: parentEx.ignoreOmitempty,
		}
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
	ex.ignoreOmitempty = parentEx.ignoreOmitempty
	return
}
