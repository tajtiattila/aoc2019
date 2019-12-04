package main

import "testing"

func TestPW4(t *testing.T) {
	type test struct {
		pw   int
		nrun int
		ok   bool
	}

	tests := []test{
		{111111, 0, true},
		{223450, 0, false},
		{123789, 0, false},
		{123444, 2, false},
		{111122, 2, true},
	}

	for _, tt := range tests {
		if pw4(tt.pw, tt.nrun) != tt.ok {
			t.Errorf("pw(%v, %v) != %v", tt.pw, tt.nrun, tt.ok)
		}
	}
}
