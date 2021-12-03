package jsonvalue

import (
	"errors"
	"fmt"
	"strconv"
)

// For state machine chart, please refer to ./img/parse_float_state_chart.drawio

func (it *iter) parseNumber(
	offset int,
) (v *V, end int, reachEnd bool, err error) {

	stm := newFloatStateMachine(it, offset)

	for err == nil {
		b, bType := stm.pop()
		if bType == numRuneInvalid {
			err = stm.parseResult()
			break
		}

		switch stm.state {
		default:
		case stateStart:
			err = stm.stateStart(b, bType)
		case stateLeadingZero:
			err = stm.stateLeadingZero(b, bType)
		case stateLeadingNegative:
			err = stm.stateLeadingNegative(b, bType)
		case stateLeadingDigit:
			err = stm.stateLeadingDigit(b, bType)
		case stateFraction:
			err = stm.stateFraction(b, bType)
		case stateIntegerDigit:
			err = stm.stateIntegerDigit(b, bType)
		case stateFractionDigit:
			err = stm.stateFractionDigit(b, bType)
		case stateExponent:
			err = stm.stateExponent(b, bType)
		case stateExponentSign:
			err = stm.stateExponentSign(b, bType)
		case stateExponentDigit:
			err = stm.stateExponentDigit(b, bType)
		}
	}

	return stm.res, stm.progress.offset, stm.progress.remain == 0, err
}

type floatStateMachineState uint

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

type floatStateMachine struct {
	it    *iter
	state floatStateMachineState

	progress struct {
		start  int
		offset int
		remain int
	}

	status struct {
		integer       uint64
		negative      bool
		fractionStart int
		exponent      uint64
		exponentStart int
	}

	res *V
}

func newFloatStateMachine(it *iter, offset int) *floatStateMachine {
	stm := &floatStateMachine{}
	stm.it = it
	stm.progress.start = offset
	stm.progress.offset = offset
	stm.progress.remain = len(it.b) - offset
	stm.state = stateStart
	stm.status.fractionStart = -1
	stm.status.exponentStart = -1

	if it.numRuneTypes == nil {
		it.numRuneTypes = make([]numRuneType, 256)
		for i := '1'; i <= '9'; i++ {
			it.numRuneTypes[i] = numRuneDigitOneNine
		}
		it.numRuneTypes['0'] = numRuneDigitZero
		it.numRuneTypes['E'] = numRuneExponent
		it.numRuneTypes['e'] = numRuneExponent
		it.numRuneTypes['-'] = numRuneNegative
		it.numRuneTypes['+'] = numRunePositive
		it.numRuneTypes['.'] = numRuneFraction
	}

	return stm
}

func (s *floatStateMachine) pop() (b byte, typ numRuneType) {
	if s.progress.remain == 0 {
		return 0, numRuneInvalid
	}

	b = s.it.b[s.progress.offset]

	if typ = s.it.numRuneTypes[int(b)]; typ != numRuneInvalid {
		s.progress.offset++
		s.progress.remain--
		return b, typ
	}
	return 0, numRuneInvalid
}

func (s *floatStateMachine) errorf(f string, a ...interface{}) error {
	a = append([]interface{}{s.progress.offset}, a...)
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
	// a = append([]interface{}{ca, string(s.it.b[s.progress.start:]), s.progress.offset}, a...)
	// return fmt.Errorf("%s - parsing number ('%s') at index %d: "+f, a...)
}

const (
	uintMaxStr    = "18446744073709551615"
	uintMaxDigits = 10000000000000000000
	intMin        = -9223372036854775808
	intMinStr     = "-9223372036854775808"
	intMinAbs     = 9223372036854775808
)

func (s *floatStateMachine) parseResult() error {
	switch s.state {
	case stateLeadingZero, stateLeadingDigit, stateIntegerDigit:
		if s.status.negative {
			return s.parseNegativeIntResult()
		}
		return s.parsePositiveIntResult()
	case stateFractionDigit, stateExponentDigit:
		return s.parseFloatResult()
	default:
		return errors.New("invalid number") // TODO:
	}
}

func (s *floatStateMachine) parseFloatResult() error {
	le := s.progress.offset - s.progress.start
	bytes := s.it.b[s.progress.start:s.progress.offset]

	str := string(bytes)
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return s.errorf("%w", err)
	}

	v := new(Number)
	v.srcByte = bytes
	v.srcOffset, v.srcEnd = 0, le

	v.parsed = true

	v.num.negative = f < 0
	v.num.floated = true
	v.num.i64 = int64(f)
	v.num.u64 = uint64(f)
	v.num.f64 = f

	s.res = v
	return nil
}

func (s *floatStateMachine) parsePositiveIntResult() error {
	le := s.progress.offset - s.progress.start

	if le > len(uintMaxStr) {
		return s.errorf("value too large")
	} else if le == len(uintMaxStr) {
		if s.status.integer < uintMaxDigits {
			return s.errorf("value too large")
		}
	}

	v := new(Number)
	v.srcByte = s.it.b
	v.srcOffset, v.srcEnd = s.progress.offset, s.progress.offset+le

	v.parsed = true

	v.num.negative = false
	v.num.floated = false
	v.num.i64 = int64(s.status.integer)
	v.num.u64 = uint64(s.status.integer)
	v.num.f64 = float64(s.status.integer)

	s.res = v
	return nil
}

func (s *floatStateMachine) parseNegativeIntResult() error {
	le := s.progress.offset - s.progress.start

	if le > len(intMinStr) {
		return s.errorf("absolute value too large")
	} else if le == len(intMinStr) {
		if s.status.integer > intMinAbs {
			return s.errorf("absolute value too large")
		}
	}

	v := new(Number)
	v.srcByte = s.it.b
	v.srcOffset, v.srcEnd = s.progress.offset, s.progress.offset+le

	v.parsed = true
	v.num.negative = true
	v.num.floated = false

	if s.status.integer == intMinAbs {
		v.num.i64 = intMin
	} else {
		v.num.i64 = -int64(s.status.integer)
	}

	v.num.u64 = uint64(v.num.i64)
	v.num.f64 = float64(s.status.integer)

	s.res = v
	return nil
}

func (s *floatStateMachine) stateStart(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitZero:
		s.state = stateLeadingZero
	case numRuneDigitOneNine:
		s.status.integer = uint64(b) - '0'
		s.state = stateLeadingDigit
	case numRuneNegative:
		s.status.negative = true
		s.state = stateLeadingNegative
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}

func (s *floatStateMachine) stateLeadingZero(b byte, typ numRuneType) error {
	switch typ {
	case numRuneExponent:
		s.status.exponentStart = s.progress.offset
		s.state = stateExponent
	case numRuneFraction:
		s.status.fractionStart = s.progress.offset
		s.state = stateFraction
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}

func (s *floatStateMachine) stateLeadingDigit(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitZero, numRuneDigitOneNine:
		s.status.integer = s.status.integer*10 + uint64(b-'0')
		s.state = stateIntegerDigit
	case numRuneFraction:
		s.status.fractionStart = s.progress.offset
		s.state = stateFraction
	case numRuneExponent:
		s.status.exponentStart = s.progress.offset
		s.state = stateExponent
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}

func (s *floatStateMachine) stateLeadingNegative(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitOneNine:
		s.status.integer = s.status.integer*10 + uint64(b-'0')
		s.state = stateLeadingDigit
	case numRuneDigitZero:
		s.state = stateLeadingZero
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}

func (s *floatStateMachine) stateIntegerDigit(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		s.status.integer = s.status.integer*10 + uint64(b-'0')
	case numRuneFraction:
		s.status.fractionStart = s.progress.offset
		s.state = stateFraction
	case numRuneExponent:
		s.status.exponentStart = s.progress.offset
		s.state = stateExponent
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}

func (s *floatStateMachine) stateFraction(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		s.status.fractionStart = s.progress.offset
		s.state = stateFractionDigit
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}

func (s *floatStateMachine) stateExponent(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		s.status.exponent = uint64(b - '0')
		s.state = stateExponentDigit
	case numRunePositive, numRuneNegative:
		s.state = stateExponentSign
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}

func (s *floatStateMachine) stateExponentSign(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		s.status.exponent = uint64(b - '0')
		s.state = stateExponentDigit
	default:
		return s.errorf("illegal character 0x%02x after exponent", b)
	}
	return nil
}

func (s *floatStateMachine) stateFractionDigit(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		// OK
	case numRuneExponent:
		s.status.exponentStart = s.progress.offset
		s.state = stateExponent
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}

func (s *floatStateMachine) stateExponentDigit(b byte, typ numRuneType) error {
	switch typ {
	case numRuneDigitOneNine, numRuneDigitZero:
		// OK
	default:
		return s.errorf("illegal character 0x%02x", b)
	}
	return nil
}
