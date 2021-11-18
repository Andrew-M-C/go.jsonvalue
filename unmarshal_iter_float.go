package jsonvalue

import (
	"fmt"
	"strconv"
)

// For state machine chart, please refer to ./img/parse_float_state_chart.drawio

func (it *iter) parseNumber(
	offset int,
) (v *V, end int, reachEnd bool, err error) {

	stm := newFloatStateMachine(it, offset)
	for stm.res == nil {
		if err = stm.next(); err != nil {
			break
		}
	}

	return stm.res, stm.progress.offset, stm.progress.remain == 0, err
}

type floatStateMachine struct {
	it   *iter
	next func() error

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
	stm := &floatStateMachine{
		it: it,
	}
	stm.progress.start = offset
	stm.progress.offset = offset
	stm.progress.remain = len(it.b) - offset
	stm.next = stm.stateStart
	stm.status.fractionStart = -1
	stm.status.exponentStart = -1
	return stm
}

func (s *floatStateMachine) hasFraction() bool {
	return s.status.fractionStart >= 0
}

func (s *floatStateMachine) hasExponent() bool {
	return s.status.exponentStart >= 0
}

func (s *floatStateMachine) pop() (b byte, ok bool) {
	if s.progress.remain == 0 {
		return 0, false
	}
	b = s.it.b[s.progress.offset]

	good := func() (byte, bool) {
		s.progress.offset++
		s.progress.remain--
		return b, true
	}

	if b >= '0' && b <= '9' {
		return good()
	}
	if b == '-' || b == '+' {
		return good()
	}
	if b == '.' {
		return good()
	}
	if b == 'E' || b == 'e' {
		return good()
	}

	return 0, false
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
	// a = append([]interface{}{ca, s.progress.offset}, a...)
	// return fmt.Errorf("%s - parsing number at index %d: "+f, a...)
}

const (
	uintMaxStr    = "18446744073709551615"
	uintMaxDigits = 10000000000000000000
	intMin        = -9223372036854775808
	intMinStr     = "-9223372036854775808"
	intMinAbs     = 9223372036854775808
)

func (s *floatStateMachine) parseResult() error {
	// parse float
	if s.hasFraction() || s.hasExponent() {
		return s.parseFloatResult()
	}

	// parse negative int
	if s.status.negative {
		return s.parseNegativeIntResult()
	}

	// parse negative int
	return s.parsePositiveIntResult()
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

func (s *floatStateMachine) stateStart() error {
	b, ok := s.pop()
	if !ok {
		return s.errorf("zero string")
	}

	switch b {
	case '0':
		s.next = s.stateLeadingZero
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.status.integer = uint64(b) - '0'
		s.next = s.stateLeadingDigit
	case '-':
		s.status.negative = true
		s.next = s.stateLeadingNegative
	default:
		return s.errorf("illegal character 0x%02x", b)
	}

	return nil
}

func (s *floatStateMachine) stateLeadingZero() error {
	b, ok := s.pop()
	if !ok {
		s.res = NewInt(0)
		return nil
	}

	switch b {
	case 'E', 'e':
		s.status.exponentStart = s.progress.offset
		s.next = s.stateExponent
	case '.':
		s.status.fractionStart = s.progress.offset
		s.next = s.stateFraction
	default:
		return s.errorf("illegal character 0x%02x", b)
	}

	return nil
}

func (s *floatStateMachine) stateLeadingDigit() error {
	b, ok := s.pop()
	if !ok {
		return s.parseResult()
	}

	if b >= '0' && b <= '9' {
		s.status.integer = s.status.integer*10 + uint64(b-'0')
		s.next = s.stateIntegerDigit
	} else if b == '.' {
		s.status.fractionStart = s.progress.offset
		s.next = s.stateFraction
	} else if b == 'E' || b == 'e' {
		s.status.exponentStart = s.progress.offset
		s.next = s.stateExponent
	} else {
		return s.errorf("illegal character 0x%02x", b)
	}

	return nil
}

func (s *floatStateMachine) stateLeadingNegative() error {
	b, ok := s.pop()
	if !ok {
		return s.errorf("expect digit after negative symbol")
	}

	if b >= '1' && b <= '9' {
		s.status.integer = s.status.integer*10 + uint64(b-'0')
		s.next = s.stateLeadingDigit

	} else if b == '0' {
		s.next = s.stateLeadingZero

	} else {
		return s.errorf("illegal character 0x%02x", b)
	}

	return nil
}

func (s *floatStateMachine) stateIntegerDigit() error {
	b, ok := s.pop()
	if !ok {
		return s.parseResult()
	}

	if b >= '0' && b <= '9' {
		s.status.integer = s.status.integer*10 + uint64(b-'0')
	} else if b == 'E' || b == 'e' {
		s.status.exponentStart = s.progress.offset
		s.next = s.stateExponent
	} else if b == '.' {
		s.status.fractionStart = s.progress.offset
		s.next = s.stateFraction
	} else {
		return s.errorf("illegal character 0x%02x", b)
	}

	return nil
}

func (s *floatStateMachine) stateFraction() error {
	b, ok := s.pop()
	if !ok {
		return s.errorf("expect digit after fraction symbol")
	}

	if b >= '0' && b <= '9' {
		s.status.fractionStart = s.progress.offset
		s.next = s.stateFractionDigit
	} else {
		return s.errorf("illegal character 0x%02x", b)
	}

	return nil
}

func (s *floatStateMachine) stateExponent() error {
	b, ok := s.pop()
	if !ok {
		return s.errorf("expect digit after exponent symbol")
	}

	if b >= '0' && b <= '9' {
		s.status.exponent = uint64(b - '0')
		s.next = s.stateExponentDigit
	} else if b == '+' || b == '-' {
		s.next = s.stateExponentSign
	} else {
		return s.errorf("illegal character 0x%02x after exponent", b)
	}

	return nil
}

func (s *floatStateMachine) stateExponentSign() error {
	b, ok := s.pop()
	if !ok {
		return s.errorf("expect digit after signed symbol")
	}

	if b >= '0' && b <= '9' {
		s.status.exponent = uint64(b - '0')
		s.next = s.stateExponentDigit
	} else {
		return s.errorf("illegal character 0x%02x after exponent", b)
	}

	return nil
}

func (s *floatStateMachine) stateFractionDigit() error {
	b, ok := s.pop()
	if !ok {
		return s.parseResult()
	}

	if b >= '0' && b <= '9' {
		// continue
	} else if b == 'E' || b == 'e' {
		s.status.exponentStart = s.progress.offset
		s.next = s.stateExponent
	} else {
		return s.errorf("illegal character 0x%02x", b)
	}

	return nil
}

func (s *floatStateMachine) stateExponentDigit() error {
	b, ok := s.pop()
	if !ok {
		return s.parseResult()
	}

	if b >= '0' && b <= '9' {
		// continue
	} else {
		return s.errorf("illegal character 0x%02x", b)
	}

	return nil
}
