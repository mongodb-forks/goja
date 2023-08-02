package goja

import (
	"testing"
)

func TestStringAsciiMemUsage(t *testing.T) {
	tests := []struct {
		name        string
		val         asciiString
		expected    uint64
		newExpected uint64
		errExpected error
	}{
		{
			name:        "should have a value given the length of the string",
			val:         asciiString("yo"),
			expected:    2,
			newExpected: 2 + SizeString, // length with string overhead
			errExpected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			total, newTotal, err := tc.val.MemUsage(NewMemUsageContext(New(), 100, 100, 100, 100, nil))
			if err != tc.errExpected {
				t.Fatalf("Unexpected error. Actual: %v Expected: %v", err, tc.errExpected)
			}
			if err != nil && tc.errExpected != nil && err.Error() != tc.errExpected.Error() {
				t.Fatalf("Errors do not match. Actual: %v Expected: %v", err, tc.errExpected)
			}
			if total != tc.expected {
				t.Fatalf("Unexpected memory return. Actual: %v Expected: %v", total, tc.expected)
			}
			if newTotal != tc.newExpected {
				t.Fatalf("Unexpected new memory return. Actual: %v Expected: %v", newTotal, tc.newExpected)
			}
		})
	}
}
