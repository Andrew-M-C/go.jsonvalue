package jsonvalue

const (
	initialArrayCapacity = 32
)

// Deprecated: Opt is the option of jsonvalue in marshaling. This type is deprecated,
// please use OptXxxx() functions instead.
//
// Opt 表示序列化当前 jsonvalue 类型时的参数。这个类型后续可能不再迭代新字段了，请改用 OptXxxx() 函数进行配置。
type Opt struct {
	// OmitNull tells how to handle null json value. The default value is false.
	// If OmitNull is true, null value will be omitted when marshaling.
	//
	// OmitNull 表示是否忽略 JSON 中的 null 类型值。默认为 false.
	OmitNull bool

	// MarshalLessFunc is used to handle sequences of marshaling. Since object is
	// implemented by hash map, the sequence of keys is unexpectable. For situations
	// those need settled JSON key-value sequence, please use MarshalLessFunc.
	//
	// Note: Elements in an array value would NOT trigger this function as they are
	// already sorted.
	//
	// We provides a example DefaultStringSequence. It is quite useful when calculating
	// idempotence of a JSON text, as key-value sequences should be fixed.
	//
	// MarshalLessFunc 用于处理序列化 JSON 对象类型时，键值对的顺序。由于 object 类型是采用 go 原生的 map 类型，采用哈希算法实现，
	// 因此其键值对的顺序是不可控的。而为了提高效率，jsonvalue 的内部实现中并不会刻意保存键值对的顺序。如果有必要在序列化时固定键值对顺序的话，
	// 可以使用这个函数。
	//
	// 注意：array 类型中键值对的顺序不受这个函数的影响
	//
	// 此外，我们提供了一个例子: DefaultStringSequence。当需要计算 JSON 文本的幂等值时，
	// 由于需要不变的键值对顺序，因此这个函数是非常有用的。
	MarshalLessFunc MarshalLessFunc

	// MarshalKeySequence is used to handle sequance of marshaling. This is much simpler
	// than MarshalLessFunc, just pass a string slice identifying key sequence. For keys
	// those are not in this slice, they would be appended in the end according to result
	// of Go string comparing. Therefore this parameter is useful for ensure idempotence.
	//
	// MarshalKeySequence 也用于处理序列化时的键值对顺序。与 MarshalLessFunc 不同，这个只需要用字符串切片的形式指定键的顺序即可，
	// 实现上更为简易和直观。对于那些不在指定切片中的键，那么将会统一放在结尾，并且按照 go 字符串对比的结果排序。也可以保证幂等。
	MarshalKeySequence []string
	keySequence        map[string]int // generated from MarshalKeySequence

	// FloatNaNHandleType tells what to deal with float NaN.
	//
	// FloatNaNHandleType 表示当处理 float 的时候，如果遇到了 NaN 的话，要如何处理。
	FloatNaNHandleType FloatNaNHandleType
	// FloatNaNToString works with FloatNaNHandleType = FloatNaNConvertToString. It tells what string to replace
	// to with NaN. If not specified, NaN will be set as string "NaN".
	//
	// FloatNaNToString 搭配 FloatNaNHandleType = FloatNaNConvertToString 使用，表示将 NaN 映射为哪个字符串。
	// 这个值如果不指定，则默认会被设置为字符串 "NaN"
	FloatNaNToString string
	// FloatNaNToFloat works with FloatNaNHandleType = FloatNaNConvertToFloat. It tells what float number will
	// be mapped to as for NaN. NaN, +Inf or -Inf are not allowed for this option.
	//
	// FloatNaNToFloat 搭配 FloatNaNHandleType = FloatNaNConvertToFloat 使用，表示将 NaN 映射为哪个 float64 值。
	// 不允许指定为 NaN, +Inf 或 -Inf。如果不指定，则映射为 0
	FloatNaNToFloat float64

	// FloatInfHandleType tells what to deal with float +Inf and -Inf.
	//
	// FloatInfHandleType 表示当处理 float 的时候，如果遇到了 +Inf 和 -Inf 的话，要如何处理。
	FloatInfHandleType FloatInfHandleType
	// FloatInfPositiveToString works with FloatInfHandleType = FloatInfConvertToFloat. It tells what float number will
	// be mapped to as for +Inf. If not specified, +Inf will be set as string "+Inf"
	//
	// FloatInfPositiveToString 搭配 FloatInfHandleType = FloatInfConvertToFloat 使用，表示将 NaN 映射为哪个字符串。
	// 这个值如果不指定，则默认会被设置为字符串 "+Inf"
	FloatInfPositiveToString string
	// FloatInfNegativeToString works with FloatInfHandleType = FloatInfConvertToFloat. It tells what float number will
	// be mapped to as for -Inf. If not specified, -Inf will be set as string "-" + strings.TrimLeft(FloatInfPositiveToString, "+").
	//
	// FloatInfNegativeToString 搭配 FloatInfHandleType = FloatInfConvertToFloat 使用，表示将 NaN 映射为哪个字符串。
	// 这个值如果不指定，则默认会被设置为字符串 "-" + strings.TrimLeft(FloatInfPositiveToString, "+")。
	FloatInfNegativeToString string
	// FloatInfToFloat works with FloatInfHandleType = FloatInfConvertToFloat. It tells what float numbers will be
	// mapped to as for +Inf. And -Inf will be specified as the negative value of this option.
	// +Inf or -Inf are not allowed for this option.
	//
	// FloatInfToFloat 搭配 FloatInfHandleType = FloatInfConvertToFloat 使用，表示将 +Inf 映射为哪个 float64 值。而 -Inf
	// 则会被映射为这个值的负数。
	// 不允许指定为 NaN, +Inf 或 -Inf。如果不指定，则映射为 0
	FloatInfToFloat float64

	// escapeHTML tells what do deal with &, <, > character. Default value is nil, which tells using the default value,
	// which should be 'true'.
	escapeHTML *bool
}

func (o *Opt) shouldEscapeHTML() bool {
	return o.escapeHTML == nil || *o.escapeHTML
}

type FloatNaNHandleType uint8

const (
	// FloatNaNTreatAsError indicates that error will be returned when a float number is NaN when marshaling.
	//
	// FloatNaNTreatAsError 表示当 marshal 遇到 NaN 时，返回错误。这是默认选项。
	FloatNaNTreatAsError FloatNaNHandleType = 0
	// FloatNaNConvertToFloat indicates that NaN will be replaced as another float number when marshaling. This option
	// works with option FloatNaNToFloat.
	//
	// FloatNaNConvertToFloat 表示当 marshal 遇到 NaN 时，将值置为另一个数。搭配 FloatNaNToFloat 选项使用。
	FloatNaNConvertToFloat FloatNaNHandleType = 1
	// FloatNaNNull indicates that NaN key-value pair will be set as null when marshaling.
	//
	// FloatNaNNull 表示当 marshal 遇到 NaN 时，则将值设置为 null
	FloatNaNNull FloatNaNHandleType = 2
	// FloatNaNConvertToString indicates that NaN will be replaced as a string when marshaling. This option
	// works with option FloatNaNToString.
	//
	// FloatNaNConvertToString 表示当 marshal 遇到 NaN 时，将值设置为一个字符串。搭配 FloatNaNToString 选项使用。
	FloatNaNConvertToString FloatNaNHandleType = 3
)

type FloatInfHandleType uint8

const (
	// FloatInfTreatAsError indicates that error will be returned when a float number is Inf or -Inf when marshaling.
	//
	// FloatInfTreatAsError 表示当 marshal 遇到 Inf 或 -Inf 时，返回错误。这是默认选项。
	FloatInfTreatAsError FloatInfHandleType = 0
	// FloatInfConvertToFloat indicates that Inf and -Inf will be replaced as another float number when marshaling.
	// This option works with option FloatInfToFloat.
	//
	// FloatInfConvertToFloat 表示当 marshal 遇到 Inf 或 -Inf 时，将值置为另一个数。搭配 FloatInfToFloat 选项使用。
	FloatInfConvertToFloat FloatInfHandleType = 1
	// FloatInfNull indicates that Inf or -Inf key-value pair will be set as null when marshaling.
	//
	// FloatInfNull 表示当 marshal 遇到 Inf 和 -Inf 时，则将值设置为 null
	FloatInfNull FloatInfHandleType = 2
	// FloatInfConvertToString indicates that Inf anf -Inf will be replaced as a string when marshaling. This option
	// works with option FloatInfPositiveToString and FloatInfNegativeToString.
	//
	// FloatInfConvertToString 表示当 marshal 遇到 Inf 和 -Inf 时，将值设置为一个字符串。搭配 FloatInfPositiveToString
	// FloatInfNegativeToString 选项使用。
	FloatInfConvertToString FloatInfHandleType = 3
)

// Option is used for additional options when marshaling. Can be either a Opt{} (not pointer to it) or other
// options generated by jsonvalue.OptXxxx() functions.
//
// Option 表示用于序列化的额外选项。可以是一个 Opt{} 结构体值（而不是它的指针），或者是使用 jsonvalue.OptXxxx() 函数生成的选项。
type Option interface {
	mergeTo(*Opt)
}

func (o Opt) mergeTo(tgt *Opt) {
	*tgt = o
}

// CombineOptions is a function for internal use, which combine severial Options together. Please
// do not use this.
//
// CombineOptions 用于 jsonvalue 内部使用，合并入参的多个额外选项。
func CombineOptions(opts []Option) *Opt {
	return combineOptions(opts)
}

func combineOptions(opts []Option) *Opt {
	opt := &Opt{}
	for _, o := range opts {
		o.mergeTo(opt)
	}
	return opt
}

// ==== OmitNull ====

// OptOmitNull configures OmitNull field in Opt{}, identifying whether null values should be omitted when marshaling.
//
// OptOmitNull 配置 Opt{} 中的 OmitNull 字段，表示是否忽略 null 值。
func OptOmitNull(b bool) Option {
	return &optOmitNull{b: b}
}

type optOmitNull struct {
	b bool
}

func (o *optOmitNull) mergeTo(opt *Opt) {
	opt.OmitNull = o.b
}

// ==== MarshalLessFunc ===

// OptKeySequenceWithLessFunc configures MarshalLessFunc field in Opt{}, which defines key sequence when marshaling.
//
// OptKeySequenceWithLessFunc 配置 Opt{} 中的 MarshalLessFunc 字段，配置序列化时的键顺序。
func OptKeySequenceWithLessFunc(f MarshalLessFunc) Option {
	return &optMarshalLessFunc{f: f}
}

// OptDefaultStringSequence configures MarshalLessFunc field in Opt{} as jsonvalue.DefaultStringSequence, which
// is dictionary sequence.
//
// OptDefaultStringSequence 配置 Opt{} 中的 MarshalLessFunc 字段为 jsonvalue.DefaultStringSequence，也就是字典序。
func OptDefaultStringSequence() Option {
	return &optMarshalLessFunc{f: DefaultStringSequence}
}

type optMarshalLessFunc struct {
	f MarshalLessFunc
}

func (o *optMarshalLessFunc) mergeTo(opt *Opt) {
	if o.f != nil {
		opt.MarshalLessFunc = o.f
	}
}

// ==== MarshalKeySequence ====

// OptKeySequence configures MarshalKeySequence field in Opt{}.
//
// OptKeySequence 配置 Opt{} 中的 MarshalKeySequence 字段。
func OptKeySequence(seq []string) Option {
	return &optMarshalKeySequence{seq: seq}
}

type optMarshalKeySequence struct {
	seq []string
}

func (o *optMarshalKeySequence) mergeTo(opt *Opt) {
	opt.MarshalKeySequence = o.seq
}

// ==== FloatNaNConvertToFloat ====

// OptFloatNaNToFloat tells that when marshaling float NaN, replace it as another valid float number.
//
// OptFloatNaNToFloat 指定当遇到 NaN 时，将值替换成一个有效的 float 值。
func OptFloatNaNToFloat(f float64) Option {
	return &optFloatNaNConvertToFloat{f: f}
}

type optFloatNaNConvertToFloat struct {
	f float64
}

func (o *optFloatNaNConvertToFloat) mergeTo(opt *Opt) {
	opt.FloatNaNHandleType = FloatNaNConvertToFloat
	opt.FloatNaNToFloat = o.f
}

// ==== FloatNaNNull ====

// OptFloatNaNToNull will replace a float value to null if it is NaN.
//
// OptFloatNaNToNull 表示当遇到 NaN 时，将值替换成 null
func OptFloatNaNToNull() Option {
	return globalOptFloatNaNNull
}

type optFloatNaNNull struct{}

var (
	globalOptFloatNaNNull = &optFloatNaNNull{}
)

func (o *optFloatNaNNull) mergeTo(opt *Opt) {
	opt.FloatNaNHandleType = FloatNaNNull
}

// ==== FloatNaNConvertToString ====

// OptFloatNaNToStringNaN will replace a float value to string "NaN" if it is NaN.
//
// OptFloatNaNToStringNaN 表示遇到 NaN 时，将其替换成字符串 "NaN"
func OptFloatNaNToStringNaN() Option {
	return &optFloatNaNConvertToString{s: "NaN"}
}

// OptFloatNaNToString will replace a float value to specified string if it is NaN. If empty string is given, will
// replace as "NaN".
//
// OptFloatNaNToString 表示当遇到 NaN 时，将其替换成指定的字符串。如果指定空字符串，则替换成 "NaN"
func OptFloatNaNToString(s string) Option {
	return &optFloatNaNConvertToString{s: s}
}

type optFloatNaNConvertToString struct {
	s string
}

func (o *optFloatNaNConvertToString) mergeTo(opt *Opt) {
	opt.FloatNaNHandleType = FloatNaNConvertToString
	opt.FloatNaNToString = o.s
}

// ==== FloatInfConvertToFloat ====

// OptFloatInfToFloat will replace a +Inf float value to specified f, while -f if the value is -Inf.
//
// OptFloatInfToFloat 表示当遇到 +Inf 时，将其替换成另一个 float 值；如果是 -Inf，则会替换成其取负数。
func OptFloatInfToFloat(f float64) Option {
	return &optFloatInfConvertToFloat{f: f}
}

type optFloatInfConvertToFloat struct {
	f float64
}

func (o *optFloatInfConvertToFloat) mergeTo(opt *Opt) {
	opt.FloatInfHandleType = FloatInfConvertToFloat
	opt.FloatInfToFloat = o.f
}

// ==== FloatInfNull ====

// OptFloatInfToNull will replace a float value to null if it is +/-Inf.
//
// OptFloatInfToNull 表示当遇到 +/-Inf 时，将值替换成 null
func OptFloatInfToNull() Option {
	return globalOptFloatInfNull
}

type optFloatInfNull struct{}

var (
	globalOptFloatInfNull = &optFloatInfNull{}
)

func (o *optFloatInfNull) mergeTo(opt *Opt) {
	opt.FloatInfHandleType = FloatInfNull
}

// ==== FloatInfConvertToString ====

// OptFloatInfToStringInf will replace a +Inf value to string "+Inf", while -Inf to "-Inf"
//
// OptFloatInfToStringInf 表示遇到 +/-Inf 时，相应地将其替换成字符串 "+Inf" 和 "-Inf"
func OptFloatInfToStringInf() Option {
	return &optFloatInfConvertToString{}
}

// OptFloatInfToString tells what string to replace when marshaling +Inf and -Inf numbers.
//
// OptFloatInfToString 表示遇到 +/-Inf 时，将其替换成什么字符串。
func OptFloatInfToString(positiveInf, negativeInf string) Option {
	return &optFloatInfConvertToString{
		positive: positiveInf,
		negative: negativeInf,
	}
}

type optFloatInfConvertToString struct {
	positive string
	negative string
}

func (o *optFloatInfConvertToString) mergeTo(opt *Opt) {
	opt.FloatInfHandleType = FloatInfConvertToString
	opt.FloatInfPositiveToString = o.positive
	opt.FloatInfNegativeToString = o.negative
}

// ==== escapeHTML ====

// OptEscapeHTML specifies whether problematic HTML characters should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e to avoid certain safety problems that
// can arise when embedding JSON in HTML.
func OptEscapeHTML(on bool) Option {
	return &optEscapeHTML{
		escapeHTML: on,
	}
}

type optEscapeHTML struct {
	escapeHTML bool
}

func (o *optEscapeHTML) mergeTo(opt *Opt) {
	opt.escapeHTML = &o.escapeHTML
}
