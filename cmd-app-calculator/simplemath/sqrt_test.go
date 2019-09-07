package simplemath

import "testing"

func TestSqrt(t *testing.T) {
	r := Sqrt(8)
	if r != 2 {
		t.Errorf("Sqrt(8) failed. Got %v, expected 4.", r)
	}
}
