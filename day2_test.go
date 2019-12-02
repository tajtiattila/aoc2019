package main

import "testing"

func TestIntcomp(t *testing.T) {
	type test struct {
		prog []int

		expect int
		at     int
	}

	p := func(expect, at int, pgm ...int) test {
		return test{
			prog:   pgm,
			expect: expect,
			at:     at,
		}
	}
	tests := []test{
		p(3500, 0, 1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50),
		p(2, 0, 1, 0, 0, 0, 99),
		p(6, 3, 2, 3, 0, 3, 99),
		p(9801, 5, 2, 4, 4, 5, 99, 0),
		p(30, 0, 1, 1, 1, 4, 99, 5, 6, 0, 99),
	}

	for i, tt := range tests {
		v, _, err := Intcomp(tt.prog, 0)
		if err != nil {
			t.Fatal(err)
		}
		if v[tt.at] != tt.expect {
			t.Errorf("Program %d error, at %d got %d want %d %v", i, tt.at, v[tt.at], tt.expect, v)
		}
	}
}
