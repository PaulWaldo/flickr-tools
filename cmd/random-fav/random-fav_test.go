package main

import "testing"

func TestDivMod(t *testing.T) {
	cases := []struct {
		numer     int
		denom     int
		quotient  int
		remainder int
	}{
		{5, 12, 0, 5},
		{12, 5, 2, 2},
	}

	for _, c := range cases {
		q, r := divMod(c.numer, c.denom)
		if q != c.quotient {
			t.Errorf("For %d/%d, expecting quotient to be %d but got %d", c.numer, c.denom, c.quotient, q)
		}
		if r != c.remainder {
			t.Errorf("For %d/%d, expecting remainder to be %d but got %d", c.numer, c.denom, c.remainder, r)
		}
	}
}
