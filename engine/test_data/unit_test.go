package test_data

import "testing"

func Size(a int) string {
	switch {
	case a < 0:
		return "negative"
	case a == 0:
		return "zero"
	case a < 10:
		return "small"
	case a < 100:
		return "big"
	case a < 1000:
		return "huge"
	}
	return "enormous"
}

type Test struct {
	in  int
	out string
}

var tests = []Test{
	{-1, "negative"},
	{5, "small"},
}

func TestDemo(t *testing.T) {
	for i, test := range tests {
		size := Size(test.in)
		if size != test.out {
			t.Errorf("#%d: Size(%d)=%s; want %s", i, test.in, size, test.out)
		}
	}
}
