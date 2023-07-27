package goja

import "testing"

func TestMemUsageLimitExceeded(t *testing.T) {
	tests := []struct {
		name        string
		memUsage    uint64
		mu          *MemUsageContext
		expected    bool
		errExpected error
	}{
		{
			name:     "did not exceed returns false",
			memUsage: 12,
			mu: &MemUsageContext{
				MemUsageExceedsLimit: func(memUsage uint64) bool {
					return memUsage > 50
				},
			},
			expected:    false,
			errExpected: nil,
		},
		{
			name:        "undefined function returns error",
			memUsage:    12,
			mu:          &MemUsageContext{},
			expected:    false,
			errExpected: errMemUsageExceedsLimitNil,
		},
		{
			name:     "memory exceeds threshold returns true",
			memUsage: 700,
			mu: &MemUsageContext{
				MemUsageExceedsLimit: func(memUsage uint64) bool {
					return memUsage > 50
				},
			},
			expected:    true,
			errExpected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.mu.MemUsageLimitExceeded(tc.memUsage)
			if actual != tc.expected {
				t.Fatalf("ACTUAL: %v EXPECTED: %v", actual, tc.expected)
			}
			if err == nil && tc.errExpected != nil || err != nil && tc.errExpected == nil {
				t.Fatalf("Unexpected error. Actual: %v Expected; %v", err, tc.errExpected)
			}
			if err != nil && tc.errExpected != nil && err.Error() != tc.errExpected.Error() {
				t.Fatalf("Errors do not match. Actual: %v Expected: %v", err, tc.errExpected)
			}
		})
	}
}
