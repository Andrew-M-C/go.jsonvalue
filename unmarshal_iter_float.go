package jsonvalue

import (
	"fmt"
	"strconv"
)

// For state machine chart, please refer to ./img/parse_float_state_chart.drawio

var (
	numRuneTypes []numRuneType
)

func init() {
	numRuneTypes = make([]numRuneType, 256)
	for i := '1'; i <= '9'; i++ {
		numRuneTypes[i] = numRuneDigitOneNine
	}
	numRuneTypes['0'] = numRuneDigitZero
	numRuneTypes['E'] = numRuneExponent
	numRuneTypes['e'] = numRuneExponent
	numRuneTypes['-'] = numRuneNegative
	numRuneTypes['+'] = numRunePositive
	numRuneTypes['.'] = numRuneFraction
}

func (it iter) parseNumber(
	offset int,
) (v *V, end int, reachEnd bool, err error) {

	stm := newFloatStateMachine(it, offset)
	var bType numRuneType
	var b byte
	var integer uint64

	for err == nil {
		stm, b, bType = stm.pop(it)
		if bType == numRuneInvalid {
			v, err = stm.parseResult(offset, stm.offset(), it, integer)
			break
		}

		switch stm.state() {
		default:
		case stateStart:
			stm, integer, err = stm.stateStart(b, bType)
		case stateLeadingZero:
			stm, err = stm.stateLeadingZero(b, bType)
		case stateLeadingNegative:
			stm, integer, err = stm.stateLeadingNegative(b, bType, integer)
		case stateLeadingDigit:
			stm, integer, err = stm.stateLeadingDigit(b, bType, integer)
		case stateFraction:
			stm, err = stm.stateFraction(b, bType)
		case stateIntegerDigit:
			stm, integer, err = stm.stateIntegerDigit(b, bType, integer)
		case stateFractionDigit:
			stm, err = stm.stateFractionDigit(b, bType)
		case stateExponent:
			stm, err = stm.stateExponent(b, bType)
		case stateExponentSign:
			stm, err = stm.stateExponentSign(b, bType)
		case stateExponentDigit:
			err = stm.stateExponentDigit(b, bType)
		}
	}

	return v, stm.offset(), stm.remain() == 0, err
}

type floatStateMachineState uint8

const (
	stateStart floatStateMachineState = iota
	stateLeadingZero
	stateLeadingNegative
	stateLeadingDigit
	stateFraction
	stateIntegerDigit
	stateFractionDigit
	stateExponent
	stateExponentSign
	stateExponentDigit
)

type numRuneType uint8

const (
	numRuneInvalid numRuneType = iota
	numRuneDigitOneNine
	numRuneDigitZero
	numRuneNegative
	numRunePositive
	numRuneFraction
	numRuneExponent
)

type floatStateMachine uint64

func (stm floatStateMachine) state() floatStateMachineState {
	return floatStateMachineState(stm & 0xFF)
}

// 最低一个字节用于 state
func (stm floatStateMachine) withState(state floatStateMachineState) floatStateMachine {
	stm = stm & 0xFFFFFFFFFFFFFF00
	stm = stm | floatStateMachine(state)
	return stm
}

// remain 高 uint31 用于 remain
func (stm floatStateMachine) remain() int {
	return int(stm & 0x7FFFFFFF00000000 >> 32)
}

func (stm floatStateMachine) withRemain(remain int) floatStateMachine {
	return stm | (floatStateMachine(remain) << 32)
}

func (stm floatStateMachine) withRemainMinusOne() floatStateMachine {
	return stm - 0x100000000
}

// offset 低 uint32 - uint8 用于 offset
func (stm floatStateMachine) offset() int {
	return int(stm & 0xFFFFFF00 >> 8)
}

func (stm floatStateMachine) withOffset(offset int) floatStateMachine {
	return stm | (floatStateMachine(offset) << 8)
}

func (stm floatStateMachine) withOffsetAddOne() floatStateMachine {
	return stm + 0x100
}

// negative 最高位用于 negative
func (stm floatStateMachine) negative() bool {
	return stm&(0x8000000000000000) != 0
}

func (stm floatStateMachine) withNegative() floatStateMachine {
	return stm | 0x8000000000000000
}

func newFloatStateMachine(it iter, offset int) floatStateMachine {
	remain := len(it) - offset
	stm := floatStateMachine(0)
	stm = stm.withRemain(remain)
	stm = stm.withOffset(offset)

	return stm
}

func (s floatStateMachine) pop(it iter) (_ floatStateMachine, b byte, typ numRuneType) {
	remain := s.remain()
	if remain == 0 {
		return s, 0, numRuneInvalid
	}

	b = it[s.offset()]

	if typ = numRuneTypes[int(b)]; typ != numRuneInvalid {
		s = s.withOffsetAddOne()
		s = s.withRemainMinusOne()
		return s, b, typ
	}
	return s, 0, numRuneInvalid
}

// var (
// 	lastBytes []byte
// )

func (s floatStateMachine) errorf(f string, a ...interface{}) error {
	a = append([]interface{}{s.offset()}, a...)
	return fmt.Errorf("parsing number at index %d: "+f, a...)

	// debug ONLY
	// getCaller := func(skip int) string {
	// 	pc, _, _, ok := runtime.Caller(skip + 1)
	// 	if !ok {
	// 		return "<caller N/A>"
	// 	}
	// 	ca := runtime.CallersFrames([]uintptr{pc})
	// 	fr, _ := ca.Next()

	// 	fu := filepath.Ext(fr.Function)
	// 	fu = strings.TrimLeft(fu, ".")
	// 	li := fr.Line

	// 	return fmt.Sprintf("%s(), Line %d", fu, li)
	// }
	// ca := getCaller(1)

	// a = append([]interface{}{ca, string(lastBytes), s.offset()}, a...)
	// return fmt.Errorf("%s - parsing number ('%s') at index %d: "+f, a...)
}

const (
	uintMaxStr    = "18446744073709551615"
	uintMaxDigits = 10000000000000000000
	intMin        = -9223372036854775808
	intMinStr     = "-9223372036854775808"
	intMinAbs     = 9223372036854775808
)

func (s floatStateMachine) parseResult(start, end int, b []byte, integer uint64) (*V, error) {
	switch s.state() {
	case stateLeadingZero, stateLeadingDigit, stateIntegerDigit:
		if s.negative() {
			return s.parseNegativeIntResult(start, end, b, integer)
		}
		return s.parsePositiveIntResult(start, end, b, integer)
	case stateFractionDigit, stateExponentDigit:
		return s.parseFloatResult(start, end, b)
	default:
		return nil, s.errorf("invalid state: %v", s.state()) // TODO:
	}
}

func (s floatStateMachine) parseFloatResult(start, end int, b []byte) (*V, error) {
	str := string(b[start:end])
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, s.errorf("%w", err)
	}

	v := new(Number)
	v.srcByte = b
	v.srcOffset, v.srcEnd = start, end

	v.parsed = true

	v.num.negative = f < 0
	v.num.floated = true
	v.num.i64 = int64(f)
	v.num.u64 = uint64(f)
	v.num.f64 = f

	return v, nil
}

func (s floatStateMachine) parsePositiveIntResult(start, end int, b []byte, integer uint64) (*V, error) {
	le := end - start

	if le > len(uintMaxStr) {
		return nil, s.errorf("value too large")
	} else if le == len(uintMaxStr) {
		if integer < uintMaxDigits {
			return nil, s.errorf("value too large")
		}
	}

	v := new(Number)
	v.srcByte = b
	v.srcOffset, v.srcEnd = start, end

	v.parsed = true

	v.num.negative = false
	v.num.floated = false
	v.num.i64 = int64(integer)
	v.num.u64 = uint64(integer)
	v.num.f64 = float64(integer)

	return v, nil
}

func (s floatStateMachine) parseNegativeIntResult(start, end int, b []byte, integer uint64) (*V, error) {
	le := end - start

	if le > len(intMinStr) {
		return nil, s.errorf("absolute value too large")
	} else if le == len(intMinStr) {
		if integer > intMinAbs {
			return nil, s.errorf("absolute value too large")
		}
	}

	v := new(Number)
	v.srcByte = b
	v.srcOffset, v.srcEnd = start, end

	v.parsed = true
	v.num.negative = true
	v.num.floated = false

	if integer == intMinAbs {
		v.num.i64 = intMin
	} else {
		v.num.i64 = -int64(integer)
	}

	v.num.u64 = uint64(v.num.i64)
	v.num.f64 = float64(integer)

	return v, nil
}

func (s floatStateMachine) stateStart(b byte, typ numRuneType) (floatStateMachine, uint64, error) {
	integer := uint64(0)
	switch typ {
	case numRuneDigitZero:
		s = s.withState(stateLeadingZero)
	case numRuneDigitOneNine:
		integer = uint64(b) - '0'
		s = s.withState(stateLeadingDigit)
	case numRuneNegative:
		s = s.withNegative()
		s = s.withState(stateLeadingNegative)
	default:
		return 0, 0, s.errorf("illegal character 0x%02x", b)
	}
	return s, integer, nil
}

func (s floatStateMachine) stateLeadingZero(b byte, typ numRuneType) (floatStateMachine, error) {
	switch typ {
	case numRuneExponent:
		s = s.withState(stateExponent)
	case numRuneFraction:
		s = s.withState(stateFraction)
	default:
		return s, s.errorf("illegal character 0x%02x", b)
	}
	return s, nil
}

func (s floatStateMachine) stateLeadingDigit(
	b byte, typ numRuneType, integer uint64,
) (floatStateMachine, uint64, error) {
	switch typ {
	case numRuneDigitZero, numRuneDigitOneNine:
		integer = integer*10 + uint64(b-'0')
		s = s.withState(stateIntegerDigit)
	case numRuneFraction:
		s = s.withState(stateFraction)
	case numRuneExponent:
		s = s.withState(stateExponent)
	default:
		return s, 0, s.errorf("illegal character 0x%02x", b)
	}
	return s, integer, nil
}

func (s floatStateMachine) stateLeadingNegative(
	b byte, typ numRuneType, integer uint64,
) (floatStateMachine, uint64, error) {
	switch typ {
	case numRuneDigitOneNine:
		integer = integer*10 + uint64(b-'0')
		s = s.withState(stateLeadingDigit)
	case numRuneDigitZero:
		s = s.withState(stateLeadingZero)
	default:
		return s, integer, s.errorf("illegal character 0x%02x", b)
	}
	return s, integer, nil
}

func (s floatStateMachine) stateIntegerDigit(
	b byte, typ numRuneType, integer uint64,
) (floatStateMachine, uint64, error) {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		integer = integer*10 + uint64(b-'0')
	case numRuneFraction:
		s = s.withState(stateFraction)
	case numRuneExponent:
		s = s.withState(stateExponent)
	default:
		return s, integer, s.errorf("illegal character 0x%02x", b)
	}
	return s, integer, nil
}

func (s floatStateMachine) stateFraction(b byte, typ numRuneType) (floatStateMachine, error) {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		s = s.withState(stateFractionDigit)
	default:
		return s, s.errorf("illegal character 0x%02x", b)
	}
	return s, nil
}

func (s floatStateMachine) stateExponent(b byte, typ numRuneType) (floatStateMachine, error) {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		s = s.withState(stateExponentDigit)
	case numRunePositive, numRuneNegative:
		s = s.withState(stateExponentSign)
	default:
		return s, s.errorf("illegal character 0x%02x", b)
	}
	return s, nil
}

func (s floatStateMachine) stateExponentSign(b byte, typ numRuneType) (floatStateMachine, error) {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		s = s.withState(stateExponentDigit)
	default:
		return s, s.errorf("illegal character 0x%02x after exponent", b)
	}
	return s, nil
}

func (s floatStateMachine) stateFractionDigit(b byte, typ numRuneType) (floatStateMachine, error) {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		// OK
	case numRuneExponent:
		s = s.withState(stateExponent)
	default:
		return s, s.errorf("illegal character 0x%02x", b)
	}
	return s, nil
}

func (s floatStateMachine) stateExponentDigit(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		// OK
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}
