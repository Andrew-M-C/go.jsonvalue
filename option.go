package jsonvalue

import "github.com/Andrew-M-C/go.jsonvalue/internal/buffer"

const (
	initialArrayCapacity = 32
	asciiSize            = 128
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

	// ignoreJsonOmitempty ignore json tag "omitempty", which means that every data
	// would be parsed into *jsonvalue.V
	ignoreJsonOmitempty bool

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

	// marshalBySetSequence enables object key sequence by when it is set.
	//
	// 按照 key 被设置的顺序处理序列化时的 marshal 顺序
	marshalBySetSequence bool

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

	// unicodeEscapingFunc defines how to escaping a unicode greater than 0x7F to buffer.
	unicodeEscapingFunc func(r rune, buf buffer.Buffer)

	// asciiCharEscapingFunc defines how to marshal bytes lower than 0x80.
	asciiCharEscapingFunc [asciiSize]func(b byte, buf buffer.Buffer)

	// escProperties
	escProperties escapingProperties

	// indent
	indent struct {
		enabled bool
		prefix  string
		indent  string
		cnt     int
	}
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

// SetDefaultMarshalOptions set default option for marshaling. It is quite
// useful to invoke this function once in certern init funciton. Or you can
// invoke it after main entry. It is goroutine-safe.
//
// Please keep in mind that it takes effect globally and affects ALL marshaling
// behaviors in the future until the process ends. Please ensure that these
// options are acceptable for ALL future marshaling.
//
// However, you can still add additional options in later marshaling.
//
// SetDefaultMarshalOptions 设置序列化时的默认参数。使用方可以在 init 函数阶段，或者是
// main 函数启动后立刻调用该函数，以调整序列化时的默认行为。这个函数是协程安全的。
//
// 请记住，这个函数影响的是后续所有的序列化行为，请确保这个配置对后续的其他操作是可行的。
//
// 当然，你也可以在后续的操作中，基于原先配置的默认选项基础上，添加其他附加选项。
func SetDefaultMarshalOptions(opts ...Option) {
	opt := emptyOptions()
	opt.combineOptionsFrom(opts)
	internal.defaultMarshalOption = opt
}

// ResetDefaultMarshalOptions reset default marshaling options to system default.
//
// ResetDefaultMarshalOptions 重设序列化时的默认选项为系统最原始的版本。
func ResetDefaultMarshalOptions() {
	internal.defaultMarshalOption = emptyOptions()
}

func emptyOptions() *Opt {
	return &Opt{}
}

func getDefaultOptions() *Opt {
	res := *internal.defaultMarshalOption
	return &res
}

func combineOptions(opts []Option) *Opt {
	opt := getDefaultOptions()
	opt.combineOptionsFrom(opts)
	return opt
}

func (opt *Opt) combineOptionsFrom(opts []Option) {
	for _, o := range opts {
		o.mergeTo(opt)
	}
	opt.parseEscapingFuncs()
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

// ==== IgnoreOmitempty ====

// OptIgnoreOmitempty is used in Import() and New() function. This option tells
// jsonvalue to ignore json tag "omitempty", which means that every field would
// be parsed into *jsonvalue.V.
//
// OptIgnoreOmitempty 用在 Import 和 New() 函数中。这个选项将会忽略 json 标签中的
// "omitempty" 参数。换句话说, 所有的字段都会被解析并包装到 *jsonvalue.V 值中。
func OptIgnoreOmitempty() Option {
	return optIgnoreOmitempty{}
}

type optIgnoreOmitempty struct{}

func (optIgnoreOmitempty) mergeTo(opt *Opt) {
	opt.ignoreJsonOmitempty = true
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
		opt.MarshalKeySequence = nil
		opt.keySequence = nil
		opt.marshalBySetSequence = false
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
	opt.MarshalLessFunc = nil
	opt.MarshalKeySequence = o.seq
	opt.keySequence = nil
	opt.marshalBySetSequence = false
}

// ==== marshalBySetSequence ====

// OptSetSequence tells that when marshaling an object, the key will be sorted by
// the time they are added into or refreshed in its parent object. The later a key
//
//	is set or updated, the later it and its value will be marshaled.
//
// OptSetSequence 指定在序列化 object 时，按照一个 key 被设置时的顺序进行序列化。如果一个
// key 越晚添加到 object 类型，则在序列化的时候越靠后。
func OptSetSequence() Option {
	return optSetSequence{}
}

type optSetSequence struct{}

func (optSetSequence) mergeTo(opt *Opt) {
	opt.MarshalLessFunc = nil
	opt.MarshalKeySequence = nil
	opt.keySequence = nil
	opt.marshalBySetSequence = true
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
	return optFloatNaNNull{}
}

type optFloatNaNNull struct{}

func (optFloatNaNNull) mergeTo(opt *Opt) {
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
	return optFloatInfNull{}
}

type optFloatInfNull struct{}

func (optFloatInfNull) mergeTo(opt *Opt) {
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
// can arise when embedding JSON in HTML. If not specified, HTML symbols above will be escaped by default.
//
// OptEscapeHTML 指定部分 HTML 符号是否会被转义。相关的 HTML 符号为 &, <, > 三个。如无指定，则默认会被转义
func OptEscapeHTML(on bool) Option {
	return optEscapeHTML(on)
}

type optEscapeHTML bool

func (o optEscapeHTML) mergeTo(opt *Opt) {
	if o {
		opt.escProperties = opt.escProperties.clear(escapeWithoutHTML)
	} else {
		opt.escProperties = opt.escProperties.set(escapeWithoutHTML)
	}
}

// ==== do or do not not use ASCII escaping ====

// OptUTF8 specifies that all unicodes greater than 0x7F, will NOT be escaped by \uXXXX format but UTF-8.
//
// OptUTF8 指定使用 UTF-8 编码。也就是说针对大于 0x7F 的 unicode 字符，将不会使用默认的 \uXXXX 格式进行编码，而是直接使用
// UTF-8。
func OptUTF8() Option {
	return optUTF8(true)
}

type optUTF8 bool

func (o optUTF8) mergeTo(opt *Opt) {
	opt.escProperties = opt.escProperties.set(escapeUTF8)
}

// ==== ignore slash ====

// OptEscapeSlash specifies whether we should escape slash (/) symbol. In JSON standard, this character
// should be escaped as '\/'. But non-escaping will not affect anything. If not specfied, slash will be
// escaped by default.
//
//	OptEscapeSlash 指定是否需要转移斜杠 (/) 符号。在 JSON 标准中这个符号是需要被转移为 '\/' 的,
//
// 但是不转义这个符号也不会带来什么问题。如无明确指定，如无指定，默认情况下，斜杠是会被转义的。
func OptEscapeSlash(on bool) Option {
	return optEscSlash(on)
}

type optEscSlash bool

func (o optEscSlash) mergeTo(opt *Opt) {
	if o {
		opt.escProperties = opt.escProperties.clear(escapeIgnoreSlash)
	} else {
		opt.escProperties = opt.escProperties.set(escapeIgnoreSlash)
	}
}

// escapingProperties is a bit mask, showing the option for escaping
// characters.
//
// escapingProperties 是一个位掩码，表明转义特殊字符的方法
type escapingProperties uint8

const (
	escapeUTF8        = 0
	escapeWithoutHTML = 1
	escapeIgnoreSlash = 2
)

func (esc escapingProperties) set(mask escapingProperties) escapingProperties {
	return esc | (1 << mask)
}

func (esc escapingProperties) clear(mask escapingProperties) escapingProperties {
	return esc & ^(1 << mask)
}

func (esc escapingProperties) has(mask escapingProperties) bool {
	return esc == esc.set(mask)
}

// parseEscapingFuncs parse escaping functions by escapingProperties.
func (o *Opt) parseEscapingFuncs() {
	// init bytes lower than 0x80
	for i := range o.asciiCharEscapingFunc {
		o.asciiCharEscapingFunc[i] = escapeNothing
	}

	// ASCII control bytes should always escaped
	o.asciiCharEscapingFunc[0x00] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x01] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x02] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x03] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x04] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x05] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x06] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x07] = escAsciiControlChar
	// 0x08 is \b, encoding/json marshal as \u0008, but according to JSON standard, it should be "\b"
	o.asciiCharEscapingFunc[0x0E] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x0F] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x10] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x11] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x12] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x13] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x14] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x15] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x16] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x17] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x18] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x19] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x1A] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x1B] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x1C] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x1D] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x1E] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x1F] = escAsciiControlChar
	o.asciiCharEscapingFunc[0x7F] = escAsciiControlChar // encoding/json does not escape DEL

	// ASCII characters always to be escaped
	o.asciiCharEscapingFunc['"'] = escDoubleQuote
	o.asciiCharEscapingFunc['/'] = escSlash
	o.asciiCharEscapingFunc['\\'] = escBackslash
	o.asciiCharEscapingFunc['\b'] = escBackspace
	o.asciiCharEscapingFunc['\f'] = escVertTab
	o.asciiCharEscapingFunc['\t'] = escTab
	o.asciiCharEscapingFunc['\n'] = escNewLine
	o.asciiCharEscapingFunc['\r'] = escReturn
	o.asciiCharEscapingFunc['<'] = escLeftAngle
	o.asciiCharEscapingFunc['>'] = escRightAngle
	o.asciiCharEscapingFunc['&'] = escAnd
	// o.asciiCharEscapingFunc['%'] = escPercent

	// unicodes >= 0x80
	if o.escProperties.has(escapeUTF8) {
		o.unicodeEscapingFunc = escapeGreaterUnicodeToBuffByUTF8
	} else {
		o.unicodeEscapingFunc = escapeGreaterUnicodeToBuffByUTF16
	}

	// ignore slash?
	if o.escProperties.has(escapeIgnoreSlash) {
		o.asciiCharEscapingFunc['/'] = escapeNothing
	}

	// without HTML?
	if o.escProperties.has(escapeWithoutHTML) {
		o.asciiCharEscapingFunc['<'] = escapeNothing
		o.asciiCharEscapingFunc['>'] = escapeNothing
		o.asciiCharEscapingFunc['&'] = escapeNothing
	}
}

// ==== indent ====

// OptIndent appliesiIndent to format the output.
//
// OptIndent 指定序列化时的缩进。
func OptIndent(prefix, indent string) Option {
	return optionIndent{prefix, indent}
}

type optionIndent [2]string

func (o optionIndent) mergeTo(opt *Opt) {
	opt.indent.enabled = true
	opt.indent.prefix = o[0]
	opt.indent.indent = o[1]
}
