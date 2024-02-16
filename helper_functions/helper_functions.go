package helper_functions

import (
	"testing"
)

// is it better to pass in the pointer or just the length of the Statements list?
func CheckProgramLength(t *testing.T, programLength int, expectedProgramLength int) {
	if programLength != expectedProgramLength {
		t.Fatalf("program does not have the expected number of statements, got=%d, but expected=%d",
			programLength, expectedProgramLength)
	}
}
