package jsonvalue

import "testing"

func TestUnmarshalRejectsInvalidSeparators(t *testing.T) {
	invalid := []string{
		`[,1]`,
		`[1,,2]`,
		`[1 2]`,
		`{"a":1,,"b":2}`,
		`{"a":1 "b":2}`,
	}

	for _, input := range invalid {
		t.Run(input, func(t *testing.T) {
			if _, err := UnmarshalString(input); err == nil {
				t.Fatalf("UnmarshalString(%q) unexpectedly succeeded", input)
			}
		})
	}
}

func TestUnmarshalAcceptsCommaSeparatedValues(t *testing.T) {
	valid := []string{
		`[]`,
		`[1,2]`,
		`[1, {"nested": true}]`,
		`{}`,
		`{"a":1,"b":2}`,
		`{"a":[1,2],"b":{"c":3}}`,
	}

	for _, input := range valid {
		t.Run(input, func(t *testing.T) {
			if _, err := UnmarshalString(input); err != nil {
				t.Fatalf("UnmarshalString(%q): %v", input, err)
			}
		})
	}
}

func TestUnmarshalPreservesTrailingCommaCompatibility(t *testing.T) {
	// Existing package tests use trailing commas in object input.
	for _, input := range []string{`[1,]`, `{"a":1,}`} {
		if _, err := UnmarshalString(input); err != nil {
			t.Fatalf("UnmarshalString(%q): %v", input, err)
		}
	}
}
