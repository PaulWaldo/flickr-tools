package utils

import (
	"os/user"
	"testing"
)

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
		q, r := DivMod(c.numer, c.denom)
		if q != c.quotient {
			t.Errorf("For %d/%d, expecting quotient to be %d but got %d", c.numer, c.denom, c.quotient, q)
		}
		if r != c.remainder {
			t.Errorf("For %d/%d, expecting remainder to be %d but got %d", c.numer, c.denom, c.remainder, r)
		}
	}
}

func TestParseDir(t *testing.T) {
	cases := []struct {
		testPath     string
		failExpected bool
		parsedPath   string
	}{
		{".", false, "."},
		{"/", false, "/"},
		{"non_existant", true, ""},
	}
	for _, c := range cases {
		parsed, err := parseDir(c.testPath)
		if c.failExpected && err == nil {
			t.Errorf("expecting failure for path %s, but got none", c.testPath)
		}
		if !c.failExpected && err != nil {
			t.Errorf("expecting pass for path '%s', but got %s", c.testPath, err)
		}
		if parsed != c.parsedPath {
			t.Errorf("Expecting parsed path to be '%s' but got '%s'", c.parsedPath, parsed)
		}
	}
}

func TestParseTildeDir(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Errorf("Got failure on current user: %s", err)
	}
	homeDir := usr.HomeDir

	parsed, err := parseDir("~")
	if err != nil {
		t.Errorf("Got failure parsing '~': %s", err)
	}

	if parsed != homeDir {
		t.Errorf("Expecting %s, but got %s", homeDir, parsed)
	}
}
