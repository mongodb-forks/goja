package goja

import (
	"testing"
)

func TestDestructMemUsage(t *testing.T) {
	tests := []struct {
		name        string
		val         *destructKeyedSource
		expected    uint64
		newExpected uint64
		errExpected error
	}{
		{
			name:        "should have a value given by the wrapped value",
			val:         &destructKeyedSource{wrapped: valueInt(99)},
			expected:    SizeInt, // wrapped value mem
			newExpected: SizeInt, // wrapped value mem
			errExpected: nil,
		},
		{
			name:        "should have a value of 0 given a nil wrapped value",
			val:         &destructKeyedSource{},
			expected:    0,
			newExpected: 0,
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
