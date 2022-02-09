package jsonvalue

import (
	"fmt"
	"strconv"
)

// For state machine chart, please refer to ./img/parse_float_state_chart.drawio

func (it iter) parseNumber(
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

	// len(it)-idx means remain bytes

	for ; len(it)-idx > 0 && !edgeFound; idx++ {
		b := it[idx]

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
					err = it.numErrorf(idx, "unexpected zero")
					return
				}
			} else if integer == 0 {
				err = it.numErrorf(idx, "unexpected zero")
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
				err = it.numErrorf(idx, "unexpected exponent symbol")
				return
			}
			exponentGot = true
			floated = true

		case '+':
			// Codes below not needed because this error is caught in outer logic
			// if !floated {
			// 	err = it.numErrorf(idx, "unexpected positive symbol")
			// 	return
			// }

		case '-':
			if !floated {
				if idx != offset {
					err = it.numErrorf(idx, "unexpected negative symbol")
					return
				}
				negative = true
			}

		case '.':
			if idx == offset || floated || exponentGot || dotGot {
				err = it.numErrorf(idx, "unexpected dot symbol")
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
			err = it.numErrorf(offset, "integer after dot missing")
			return
		}
		v, err = it.parseFloatResult(offset, idx)
	} else {
		if integer > 0 && it[offset] == '0' {
			err = it.numErrorf(offset, "non-zero integer should not start with zero")
			return
		}

		firstB := it[offset]
		if idx-offset == 1 {
			if firstB >= '0' && firstB <= '9' {
				// OK
			} else {
				err = it.numErrorf(offset, "invalid number format")
				return
			}
		}

		if negative {
			v, err = it.parseNegativeIntResult(offset, idx, integer)
		} else {
			v, err = it.parsePositiveIntResult(offset, idx, integer)
		}
	}

	return v, idx, len(it)-idx == 0, err
}

func (it iter) numErrorf(offset int, f string, a ...interface{}) error {
	a = append([]interface{}{offset}, a...)
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

	// a = append([]interface{}{ca, string(it), offset}, a...)
	// return fmt.Errorf("%s - parsing number \"%s\" at index %d: "+f, a...)
}

const (
	uintMaxStr    = "18446744073709551615"
	uintMaxDigits = 10000000000000000000
	intMin        = -9223372036854775808
	intMinStr     = "-9223372036854775808"
	intMinAbs     = 9223372036854775808
)

func (it iter) parseFloatResult(start, end int) (*V, error) {
	f, err := strconv.ParseFloat(unsafeBtoS(it[start:end]), 64)
	if err != nil {
		return nil, it.numErrorf(start, "%w", err)
	}

	return newFloat64ByRaw(f, it[start:end]), nil
}

func (it iter) parsePositiveIntResult(start, end int, integer uint64) (*V, error) {
	le := end - start

	if le > len(uintMaxStr) {
		return nil, it.numErrorf(start, "value too large")
	} else if le == len(uintMaxStr) {
		if integer < uintMaxDigits {
			return nil, it.numErrorf(start, "value too large")
		}
	}

	return NewUint64(integer), nil
}

func (it iter) parseNegativeIntResult(start, end int, integer uint64) (*V, error) {
	le := end - start

	if le > len(intMinStr) {
		return nil, it.numErrorf(start, "absolute value too large")
	} else if le == len(intMinStr) {
		if integer > intMinAbs {
			return nil, it.numErrorf(start, "absolute value too large")
		}
	}

	negative := -int64(integer)
	return NewInt64(negative), nil
}
