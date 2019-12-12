package main

import (
	"strings"
	"testing"
)

func TestJovian(t *testing.T) {
	type test struct {
		src    string
		nchk   int
		steps  int
		energy int
	}

	tests := []test{
		{
			src: `<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>
`,
			nchk:   1,
			steps:  10,
			energy: 179,
		},
		{
			src: `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>
`,
			nchk:   10,
			steps:  100,
			energy: 1940,
		},
	}
	for _, tt := range tests {
		v, err := scanjovian(strings.NewReader(tt.src))
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range v {
			t.Log(o)
		}
		simjovian(v, tt.nchk)
		t.Logf("after %v steps", tt.nchk)
		for _, o := range v {
			t.Log(o)
		}

		e := simjovian(v, tt.steps-tt.nchk)
		t.Logf("after %v steps", tt.steps)
		for _, o := range v {
			t.Log(o)
		}
		t.Log(e, tt.energy)
		if e != tt.energy {
			t.Fail()
		}
	}
}
