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

		a := phasesigN(tt.src, 100)
		b := phasesigN(tt.src, 100)
		h := len(tt.src) / 2
		if a[h:] != b[h:] {
			t.Fatal("hopp")
		}

		if len(a) < 30 {
			continue
		}
		x := []byte(tt.src)
		phaseTopHalf100(x)
		c := string(x)
		if a[h:] != c[h:] {
			t.Log(a)
			t.Log(c)
			t.Error("a!=c")
		}
	}
}

func TestSigPhaseB(t *testing.T) {
	type test struct {
		src string
		msg string
	}
	tests := []test{
		{"03036732577212944063491565474664", "84462026"},
		{"02935109699940807407585447034323", "78725270"},
		{"03081770884921959731165446850517", "53553731"},
	}

	for _, tt := range tests {
		got := phase100tenKmsg(tt.src)
		if got != tt.msg {
			t.Fatalf("%q got msg %q, want %q", tt.src, got, tt.msg)
		}
	}
}
