package goja

import (
	"testing"

	"github.com/dop251/goja/unistring"
)

func TestProxyMemUsage(t *testing.T) {
	tests := []struct {
		name        string
		val         *proxyObject
		expected    uint64
		newExpected uint64
		errExpected error
	}{
		{
			name:        "should have a value of SizeEmpty given a nil proxy object",
			val:         nil,
			expected:    SizeEmpty,
			newExpected: SizeEmpty,
			errExpected: nil,
		},
		{
			name: "should account for base object overhead given an empty proxy object",
			val:  &proxyObject{},
			// proxy overhead + baseObject overhead
			expected: SizeEmpty + SizeEmpty,
			// proxy overhead + baseObject overhead
			newExpected: SizeEmpty + SizeEmpty,
			errExpected: nil,
		},
		{
			name: "should account for baseObject and target overhead given a proxy object with empty target",
			val:  &proxyObject{target: &Object{}},
			// proxy overhead + baseObject overhead + target overhead
			expected: SizeEmpty + SizeEmpty + SizeEmpty,
			// proxy overhead + baseObject overhead + target overhead
			newExpected: SizeEmpty + SizeEmpty + SizeEmpty,
			errExpected: nil,
		},
		{
			name: "should account for baseObjet overhead and target given a proxy object with a non-empty target",
			val: &proxyObject{
				target: &Object{
					self: &baseObject{propNames: []unistring.String{"test"}, values: map[unistring.String]Value{"test": valueInt(99)}},
				},
			},
			// proxy overhead + baseObject overhead + target overhead + key/value pair
			expected: SizeEmpty + SizeEmpty + SizeEmpty + (4 + SizeInt),
			// proxy overhead + baseObject overhead + target overhead + key/value pair with string overhead
			newExpected: SizeEmpty + SizeEmpty + SizeEmpty + (4 + SizeString + SizeInt),
			errExpected: nil,
		},
		{
			name: "should account for baseObject overhead given a base dynamic array with an empty handler",
			val:  &proxyObject{handler: &jsProxyHandler{handler: &Object{}}},
			// proxy overhead + baseObject overhead + target overhead
			expected: SizeEmpty + SizeEmpty + SizeEmpty,
			// proxy overhead + baseObject overhead + target overhead
			newExpected: SizeEmpty + SizeEmpty + SizeEmpty,
			errExpected: nil,
		},
		{
			name: "should account for baseObject overhead and handler given a base dynamic array with a non-empty handler",
			val: &proxyObject{
				handler: &jsProxyHandler{
					handler: &Object{
						self: &baseObject{propNames: []unistring.String{"test"}, values: map[unistring.String]Value{"test": valueInt(99)}},
					},
				},
			},
			// proxy overhead + baseObject overhead + target overhead + key/value pair
			expected: SizeEmpty + SizeEmpty + SizeEmpty + (4 + SizeInt),
			// proxy overhead + baseObject overhead + target overhead + key/value pair with string overhead
			newExpected: SizeEmpty + SizeEmpty + SizeEmpty + (4 + SizeString + SizeInt),
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

func TestJSProxyHandlerMemUsage(t *testing.T) {
	tests := []struct {
		name        string
		val         *jsProxyHandler
		expected    uint64
		newExpected uint64
		errExpected error
	}{
		{
			name:        "should have a value of SizeEmpty given a nil proxy handler",
			val:         nil,
			expected:    SizeEmpty,
			newExpected: SizeEmpty,
			errExpected: nil,
		},
		{
			name:        "should have a value of SizeEmpty given an empty proxy handler",
			val:         &jsProxyHandler{},
			expected:    SizeEmpty,
			newExpected: SizeEmpty,
			errExpected: nil,
		},
		{
			name: "should have a value of SizeEmpty given an empty proxy handler",
			val: &jsProxyHandler{
				handler: &Object{
					self: &baseObject{propNames: []unistring.String{"test"}, values: map[unistring.String]Value{"test": valueInt(99)}},
				},
			},
			// baseObject overhead + key/value pair
			expected: SizeEmpty + (4 + SizeInt),
			// baseObject overhead + key/value pair with string overhead
			newExpected: SizeEmpty + (4 + SizeString + SizeInt),
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
