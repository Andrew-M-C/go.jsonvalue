package exporter

import (
	"fmt"
	"reflect"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

type any = interface{}

// Exporter export specified type to jsonvalue
type Exporter interface {
	Type() reflect.Type
	Export(any) *jsonvalue.V
}

// omnipotentExporter includes all type of exporter interfaces
type omnipotentExporter interface {
	fmt.Stringer
	Exporter
	structFieldExporter
}

// ParseExporter parse an prototype of type and return an exporter interface if
// its type is OK.
func ParseExporter(prototype interface{}) (Exporter, error) {
	t := reflect.TypeOf(prototype)

	if e, exist := internal.loadExportersByType(t); exist {
		return e, nil
	}

	e, err := validateTypeAndReturnExporter(t)
	if err != nil {
		return e, err
	}

	internal.storeExporterByType(t, e)
	return e, nil
}

// 检查入参合法性并返回相应的处理函数
func validateTypeAndReturnExporter(t reflect.Type) (e omnipotentExporter, err error) {
	base := baseExporter{
		typ: t,
	}

	switch t.Kind() {
	default:
		fallthrough
	case reflect.Invalid, reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedType, t)

	case reflect.Bool:
		e = &boolExporter{
			baseExporter: base,
		}

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		e = &intExporter{
			baseExporter: base,
		}

	case reflect.Uintptr, reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		e = &uintExporter{
			baseExporter: base,
		}

	case reflect.Float32:
		e = &float32Exporter{
			baseExporter: base,
		}

	case reflect.Float64:
		e = &float64Exporter{
			baseExporter: base,
		}

	case reflect.Array:
		// TODO:

	case reflect.Interface:
		// TODO:

	case reflect.Map:
		// TODO:

	case reflect.Ptr:
		// TODO: 只允许 struct ptr

	case reflect.Slice:
		// TODO:

	case reflect.String:
		e = &stringExporter{
			baseExporter: base,
		}

	case reflect.Struct:
		e = parseStructExporter(t)
	}

	return e, nil
}

func invalidJSON() *jsonvalue.V {
	internal.debugf("invalid JSON will be returned")
	return &jsonvalue.V{}
}

// -------- base of all exporter --------

type baseExporter struct {
	typ reflect.Type
}

// Type implements Exporter interface
func (e *baseExporter) Type() reflect.Type {
	return e.typ
}

func (e *baseExporter) checkInputType(v any) (reflect.Value, bool) {
	typ := reflect.TypeOf(v)
	ok := e.typ == typ
	if !ok {
		internal.debugf("exporter type '%v', but given '%v'", e.typ, typ)
		return reflect.Value{}, false
	}
	return reflect.ValueOf(v), ok
}

func (e *baseExporter) String() string {
	return fmt.Sprintf("exporter with type %v''", e.typ)
}

// -------- bool --------

type boolExporter struct {
	baseExporter
}

func (e *boolExporter) Export(rv any) *jsonvalue.V {
	val, ok := e.checkInputType(rv)
	if !ok {
		return invalidJSON()
	}
	return jsonvalue.NewBool(val.Bool())
}

// exportField 用于实现 structFieldExporter
func (e *boolExporter) exportField(field reflect.Value) (*jsonvalue.V, bool) {
	v := e.Export(field.Interface())
	return v, v.Bool()
}

// -------- int --------

type intExporter struct {
	baseExporter
}

func (e *intExporter) Export(rv any) *jsonvalue.V {
	val, ok := e.checkInputType(rv)
	if !ok {
		return invalidJSON()
	}
	return jsonvalue.NewInt64(val.Int())
}

// exportField 用于实现 structFieldExporter
func (e *intExporter) exportField(field reflect.Value) (*jsonvalue.V, bool) {
	v := e.Export(field.Interface())
	return v, v.Int64() != 0
}

// -------- uint --------

type uintExporter struct {
	baseExporter
}

func (e *uintExporter) Export(rv any) *jsonvalue.V {
	val, ok := e.checkInputType(rv)
	if !ok {
		return invalidJSON()
	}
	return jsonvalue.NewUint64(val.Uint())
}

// exportField 用于实现 structFieldExporter
func (e *uintExporter) exportField(field reflect.Value) (*jsonvalue.V, bool) {
	v := e.Export(field.Interface())
	return v, v.Uint64() != 0
}

// -------- float32 --------

type float32Exporter struct {
	baseExporter
}

func (e *float32Exporter) Export(rv any) *jsonvalue.V {
	val, ok := e.checkInputType(rv)
	if !ok {
		return invalidJSON()
	}
	return jsonvalue.NewFloat32(float32(val.Float()))
}

// exportField 用于实现 structFieldExporter
func (e *float32Exporter) exportField(field reflect.Value) (*jsonvalue.V, bool) {
	v := e.Export(field.Interface())
	return v, v.Float32() != 0
}

// -------- float64 --------

type float64Exporter struct {
	baseExporter
}

func (e *float64Exporter) Export(rv any) *jsonvalue.V {
	val, ok := e.checkInputType(rv)
	if !ok {
		return invalidJSON()
	}
	return jsonvalue.NewFloat64(val.Float())
}

// exportField 用于实现 structFieldExporter
func (e *float64Exporter) exportField(field reflect.Value) (*jsonvalue.V, bool) {
	v := e.Export(field.Interface())
	return v, v.Float64() != 0
}

// -------- string --------

type stringExporter struct {
	baseExporter
}

func (e *stringExporter) Export(rv any) *jsonvalue.V {
	val, ok := e.checkInputType(rv)
	if !ok {
		return invalidJSON()
	}
	return jsonvalue.NewString(val.String())
}

// exportField 用于实现 structFieldExporter
func (e *stringExporter) exportField(field reflect.Value) (*jsonvalue.V, bool) {
	v := e.Export(field.Interface())
	return v, field.String() != ""
}
