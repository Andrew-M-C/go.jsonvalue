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

	exportersByType map[reflect.Type]structFieldExporter // 内部引用其他类型的 exporter
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
			internal.debugf("%v - type not found: '%v'", e.typ, field.typ)
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
	fieldV = e.Export(field.Interface())
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
		internal.debugf("%v - mark type to analyze '%v'", e.typ, typ)
		e.typesToAnalyze[typ] = struct{}{}
	}
}

func (e *structExporter) storeExporter(t reflect.Type, exp omnipotentExporter) {
	delete(e.typesToAnalyze, t)
	e.exportersByType[t] = exp
	internal.debugf("%v - type analyzed: '%v'", e.typ, t)
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

	e.exportersByType[e.typ] = e
	e.parse()

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
	internal.debugf("%v: remaining struct to analyze: %+v", e.typ, e.typesToAnalyze)

	for len(e.typesToAnalyze) > 0 {
		// get one
		var t reflect.Type
		for typ := range e.typesToAnalyze {
			t = typ
			break
		}

		internal.debugf("%v - now parse type '%v'", e.typ, t)
		if exp, exist := internal.loadExportersByType(t); exist {
			internal.debugf("%v - found type in cache: '%v'", e.typ, t)
			e.storeExporter(t, exp)
			continue
		}

		// 解析类型
		switch t.Kind() {
		default:
			// do nothing

		case reflect.Array:
			// TODO:

		case reflect.Interface:
			// TODO:

		case reflect.Map:
			// TODO:

		case reflect.Slice:
			// TODO:

		case reflect.Struct:
			// TODO: FIXME: 有 bug, 递归了
			exporter := e.parseNotSelfStruct(t)
			e.storeExporter(t, exporter)

		case reflect.Pointer:
			// 检查一下是不是 struct pointer
			elemType := t.Elem()
			if elemType == e.typ { // 如果是自己的话, 那么直接取
				elemE := &pointerExporter{}
				elemE.typ = t
				elemE.elemExporter = e
				e.storeExporter(t, elemE)
				// internal.storeExporterByType(e.typ, elemE)
				internal.debugf("%v - add exporter for type %v", e.typ, t)

			} else { // 如果不是自己的话, 那么就重新取
				var exporter omnipotentExporter
				var err error
				if elemType.Kind() == reflect.Struct {
					exporter = e.parseNotSelfStruct(elemType) // struct 类型的逻辑与前面相同, 避免嵌套
				} else {
					exporter, err = validateTypeAndReturnExporter(elemType)
				}
				if err != nil {
					internal.debugf("%v - ERROR: parse type '%v' error: %v", e.typ, elemType, err)
				} else {
					elemE := &pointerExporter{}
					elemE.typ = t
					elemE.elemExporter = exporter
					e.storeExporter(t, elemE)
					// internal.storeExporterByType(e.typ, elemE)
					internal.debugf("%v - add exporter for type %v", e.typ, t)
				}
			}
		}

		delete(e.typesToAnalyze, t)
	}
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

	// 解析 exporter
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

	case reflect.Pointer:
		elemType := ft.Type.Elem()
		if elemType.Kind() != reflect.Struct {
			internal.debugf("%v: skip field type %v", ft.Type, elemType)
			return
		}
		e.markTypeToParse(ft.Type)

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

	case reflect.Pointer:
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

// 解析 struct 类型成员, 但不是自己。
//
// 针对 struct in struct, 需要再创建一个 *structExporter, 但它的 exportersByType 和
// typesToAnalyze 可以取 parent 的, 以避免循环嵌套
func (e *structExporter) parseNotSelfStruct(t reflect.Type) omnipotentExporter {
	subE := &structExporter{}
	subE.typ = t
	subE.exportersByType = e.exportersByType
	subE.typesToAnalyze = e.typesToAnalyze

	for i := 0; i < t.NumField(); i++ {
		tt := t.Field(i)

		fields := subE.parseStructFieldExporter(i, tt)
		subE.fields = append(subE.fields, fields...)
	}

	return subE
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
