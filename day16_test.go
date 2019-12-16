package main

import "testing"

func TestSigPhase(t *testing.T) {
	type test struct {
		src  string
		n    int
		want string
	}
	tests := []test{
		{"12345678", 4, "01029498"},
		{"80871224585914546619083218645595", 100, "24176176"},
		{"19617804207202209144916044189917", 100, "73745418"},
		{"69317163492948606335995924319873", 100, "52432133"},
	}

	for _, tt := range tests {
		got := phasesigN(tt.src, tt.n)[:8]
		if got != tt.want {
			t.Fatalf("%q on phase %v is %q, want %q", tt.src, tt.n, got, tt.want)
		}
	}
}
