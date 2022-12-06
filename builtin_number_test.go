package goja

import "testing"

func TestIsSafeInteger(t *testing.T) {
	const SCRIPT = `
	assert.sameValue(Number.isSafeInteger(1.0), true);
	assert.sameValue(Number.isSafeInteger(1), true);
	assert.sameValue(Number.isSafeInteger(0), true);
	assert.sameValue(Number.isSafeInteger(-1), true);
	assert.sameValue(Number.isSafeInteger('1'), false);
	assert.sameValue(Number.isSafeInteger(1.1), false);
	`
	testScript1(TESTLIB+SCRIPT, _undefined, t)
}
