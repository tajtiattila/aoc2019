package main

import (
	"strings"
	"testing"
)

func TestMultBugSim(t *testing.T) {
	type test struct {
		nstep int
		want  int // n. of bugs
		src   string
	}

	tests := []test{
		{10, 99, `
....#
#..#.
#.?##
..#..
#....`},
	}

	for _, tt := range tests {
		r := strings.NewReader(strings.TrimSpace(tt.src))
		bs, err := parsebugsim(r)
		if err != nil {
			t.Fatal(err)
		}
		got := multibugsim(bs.state, tt.nstep)
		if got != tt.want {
			t.Fatalf("multibugsim got %v, want %v", got, tt.want)
		}
	}
}
