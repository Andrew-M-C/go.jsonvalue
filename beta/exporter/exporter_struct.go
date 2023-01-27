package exporter

import (
	"fmt"
	"reflect"
	"strings"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

type structExporter struct {
	baseExporter

	fields []*structField

	exportersByType map[reflect.Type]structFieldExporter // 这个类型主要用于无锁调用
	typesToAnalyze  map[reflect.Type]struct{}            // 仅在初始化过程中使用, 当遇到了未处理的 struct 类型时, 则继续处理
}

// Export implements Exporter interface
func (e *structExporter) Export(rv any) *jsonvalue.V {
	val, ok := e.checkInputType(rv)
	if !ok {
		return invalidJSON()
	}

	v := jsonvalue.NewObject()

	for _, field := range e.fields {
		if field.exporter == nil {
			field.exporter = e.exportersByType[field.typ]
		}
		if field.exporter == nil {
			continue
		}

		fieldVal := field.getter(val)
		internal.debugf("%v: get field with type %v", val.Type(), fieldVal.Type())

		fieldV, valid := field.exporter.exportField(fieldVal)
		if !valid {
			if field.ext.omitempty {
				continue
			}
		}
		v.Set(fieldV).At(field.key)
	}

	return v
}

// exportField implements structFieldExporter interface
func (e *structExporter) exportField(field reflect.Value) (fieldV *jsonvalue.V, valid bool) {
	fieldV = e.Export(field)
	return fieldV, fieldV.Len() > 0
}

// String implements fmt.Stringer interface
func (e *structExporter) String() string {
	bdr := &strings.Builder{}
	bdr.WriteString(fmt.Sprintf("exporter with type '%v', ", e.typ))
	bdr.WriteString(fmt.Sprintf("with %d field(s)", len(e.fields)))
	return bdr.String()
}

func (e *structExporter) markTypeToParse(typ reflect.Type) {
	if typ == e.typ {
		return
	}

	if _, exist := e.exportersByType[typ]; !exist {
		e.typesToAnalyze[typ] = struct{}{}
	}
}

type structField struct {
	typ      reflect.Type
	key      string
	getter   structFieldGetter
	exporter structFieldExporter
	ext      structFieldExt
}

type structFieldExt struct {
	omitempty bool
	toString  bool
}

type structFieldExporter interface {
	exportField(field reflect.Value) (fieldV *jsonvalue.V, valid bool)
}

type structFieldGetter func(val reflect.Value) reflect.Value

func parseStructExporter(typ reflect.Type) *structExporter {
	e := &structExporter{}
	e.typ = typ
	e.exportersByType = make(map[reflect.Type]structFieldExporter, 1)
	e.typesToAnalyze = map[reflect.Type]struct{}{
		typ: {},
	}

	e.parse()
	e.exportersByType[e.typ] = e
	return e
}

// 解析总入口
func (e *structExporter) parse() {
	for i, t := 0, e.typ; i < t.NumField(); i++ {
		tt := t.Field(i)

		fields := e.parseStructFieldExporter(i, tt)
		e.fields = append(e.fields, fields...)
	}

	// TODO: 处理 typesToAnalyze
}

// 这个函数返回的数据结构可能是 nil, 这个时候要从 exportersByType 中取 exporter
func (e *structExporter) parseStructFieldExporter(
	index int, ft reflect.StructField,
) (fields []*structField) {
	if ft.Anonymous {
		return e.parseStructAnonymousFieldExporter(index, ft)
	}

	if !fieldExported(ft) {
		return
	}

	field := &structField{}

	field.key, field.ext = readFieldTag(ft, "json")
	if field.key == "-" {
		return
	}

	switch ft.Type.Kind() {
	default:
		fallthrough
	case reflect.Invalid, reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return // 上述类型不处理

	case reflect.Bool:
		fallthrough
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		fallthrough
	case reflect.Uintptr, reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		fallthrough
	case reflect.Float32, reflect.Float64, reflect.String:
		field.exporter, _ = validateTypeAndReturnExporter(ft.Type)

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

	case reflect.Struct:
		e.markTypeToParse(ft.Type)
	}

	// 成功到达这里, 那么说明是有有效 field 的
	field.typ = ft.Type
	field.getter = func(val reflect.Value) reflect.Value {
		return val.Field(index)
	}
	fields = append(fields, field)
	return
}

func (e *structExporter) parseStructAnonymousFieldExporter(
	index int, ft reflect.StructField,
) (fields []*structField) {
	// TODO: 参考 parseStructAnonymousFieldValue
	name, ex := readFieldTag(ft, "json")
	if name == "-" {
		return
	}

	t := ft.Type
	switch t.Kind() {
	default:
		fallthrough
	case reflect.Invalid, reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return // 上述类型不处理

	case reflect.Bool:
		fallthrough
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		fallthrough
	case reflect.Uintptr, reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		fallthrough
	case reflect.Float32, reflect.Float64, reflect.String:
		if !fieldExported(ft) {
			// 不可导出的直接返回
			return
		}
		// 可以导出, 非复杂类型, 那么就是单一的一个 field
		exporter, _ := validateTypeAndReturnExporter(ft.Type)
		f := &structField{
			typ: ft.Type,
			key: name,
			getter: func(val reflect.Value) reflect.Value {
				return val.Field(index)
			},
			exporter: exporter,
			ext:      ex,
		}
		fields = append(fields, f)
		return

	case reflect.Array:
		if !fieldExported(ft) {
			// 不可导出的直接返回
			return
		}
		// TODO:

	case reflect.Interface:
		// TODO:

	case reflect.Map:
		// TODO:

	case reflect.Ptr:
		// TODO: 只允许 struct ptr

	case reflect.Slice:
		if !fieldExported(ft) {
			// 不可导出的直接返回
			return
		}
		// TODO:

	case reflect.Struct:
		// TODO:
	}

	return
}

func readFieldTag(ft reflect.StructField, name string) (field string, ex structFieldExt) {
	tg := ft.Tag.Get(name)

	if tg == "" {
		field = ft.Name
		return
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

func fieldExported(ft reflect.StructField) bool {
	n := ft.Name
	if len(n) == 0 {
		return false
	}
	ini := n[0]
	return ini >= 'A' && ini <= 'Z'
}
