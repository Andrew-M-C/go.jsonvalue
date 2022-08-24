package jsonvalue

import (
	"fmt"
	"strconv"
)

// For state machine chart, please refer to ./img/parse_float_state_chart.drawio

func (u *unmarshaler) parseNumber(
	offset int,
) (v *V, end int, reachEnd bool, err error) {

	idx := offset
	negative := false
	floated := false
	exponentGot := false
	dotGot := false
	intAfterDotGot := false
	integer := uint64(0)
	edgeFound := false

	// len(u.b)-idx means remain bytes

	for ; len(u.b)-idx > 0 && !edgeFound; idx++ {
		b := u.b[idx]

		switch b {
		default:
			edgeFound = true

		case '0':
			if idx == offset {
				// OK
			} else if exponentGot {
				// OK
			} else if dotGot {
				intAfterDotGot = true
			} else if negative {
				if integer == 0 && idx != offset+1 {
					err = u.numErrorf(idx, "unexpected zero")
					return
				}
			} else if integer == 0 {
				err = u.numErrorf(idx, "unexpected zero")
				return
			}
			integer *= 10

		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if !floated {
				integer = integer*10 + uint64(b) - '0'
			} else if !exponentGot {
				intAfterDotGot = true
			}

		case 'e', 'E':
			if exponentGot {
				err = u.numErrorf(idx, "unexpected exponent symbol")
				return
			}
			exponentGot = true
			floated = true

		case '+':
			// Codes below not needed because this error is caught in outer logic
			// if !floated {
			// 	err = u.numErrorf(idx, "unexpected positive symbol")
			// 	return
			// }

		case '-':
			if !floated {
				if idx != offset {
					err = u.numErrorf(idx, "unexpected negative symbol")
					return
				}
				negative = true
			}

		case '.':
			if idx == offset || floated || exponentGot || dotGot {
				err = u.numErrorf(idx, "unexpected dot symbol")
				return
			}
			dotGot = true
			floated = true
		}
	}

	if edgeFound {
		idx--
	}

	if floated {
		if dotGot && !intAfterDotGot {
			err = u.numErrorf(offset, "integer after dot missing")
			return
		}
		v, err = u.parseFloatResult(offset, idx)
	} else {
		if integer > 0 && u.b[offset] == '0' {
			err = u.numErrorf(offset, "non-zero integer should not start with zero")
			return
		}

		firstB := u.b[offset]
		if idx-offset == 1 {
			if firstB >= '0' && firstB <= '9' {
				// OK
			} else {
				err = u.numErrorf(offset, "invalid number format")
				return
			}
		}

		if negative {
			v, err = u.parseNegativeIntResult(offset, idx, integer)
		} else {
			v, err = u.parsePositiveIntResult(offset, idx, integer)
		}
	}

	return v, idx, len(u.b)-idx == 0, err
}

func (u *unmarshaler) numErrorf(offset int, f string, a ...any) error {
	a = append([]any{offset}, a...)
	return fmt.Errorf("parsing number at index %d: "+f, a...)

	// debug ONLY below

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

	// a = append([]any{ca, string(u.b), offset}, a...)
	// return fmt.Errorf("%s - parsing number \"%s\" at index %d: "+f, a...)
}

const (
	uintMaxStr    = "18446744073709551615"
	uintMaxDigits = 10000000000000000000
	intMin        = -9223372036854775808
	intMinStr     = "-9223372036854775808"
	intMinAbs     = 9223372036854775808
)

func (u *unmarshaler) parseFloatResult(start, end int) (*V, error) {
	f, err := strconv.ParseFloat(unsafeBtoS(u.b[start:end]), 64)
	if err != nil {
		return nil, u.numErrorf(start, "%w", err)
	}

	v := u.new(Number)
	v.srcByte = u.b[start:end]

	v.num.negative = f < 0
	v.num.floated = true
	v.num.i64 = int64(f)
	v.num.u64 = uint64(f)
	v.num.f64 = f

	return v, nil
}

func (u *unmarshaler) parsePositiveIntResult(start, end int, integer uint64) (*V, error) {
	le := end - start

	if le > len(uintMaxStr) {
		return nil, u.numErrorf(start, "value too large")
	} else if le == len(uintMaxStr) {
		if integer < uintMaxDigits {
			return nil, u.numErrorf(start, "value too large")
		}
	}

	v := u.new(Number)
	v.srcByte = u.b[start:end]

	v.num.negative = false
	v.num.floated = false
	v.num.i64 = int64(integer)
	v.num.u64 = uint64(integer)
	v.num.f64 = float64(integer)

	return v, nil
}

func (u *unmarshaler) parseNegativeIntResult(start, end int, integer uint64) (*V, error) {
	le := end - start

	if le > len(intMinStr) {
		return nil, u.numErrorf(start, "absolute value too large")
	} else if le == len(intMinStr) {
		if integer > intMinAbs {
			return nil, u.numErrorf(start, "absolute value too large")
		}
	}

	v := u.new(Number)
	v.srcByte = u.b[start:end]

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
