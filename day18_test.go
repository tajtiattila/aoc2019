package main

import (
	"strings"
	"testing"

	"github.com/tajtiattila/aoc2019/rog"
)

func TestDay18a(t *testing.T) {
	type test struct {
		cost int
		src  string
	}

	var tests = []test{
		{8, `
#########
#b.A.@.a#
#########`},

		{86, `
########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################`},

		{132, `
########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################`},

		{136, `
#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`},

		{81, `
########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`},
	}

	for i, tt := range tests {
		r := strings.NewReader(strings.TrimSpace(tt.src))
		m, err := rog.Parse(r)
		if err != nil {
			t.Fatal(err)
		}

		got, err := day18a(m)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.cost {
			t.Errorf("test %d got cost %d, want %d", i, got, tt.cost)
		}
	}
}

func TestDay18b(t *testing.T) {
	type test struct {
		cost int
		src  string
	}

	var tests = []test{
		{8, `
#######
#a.#Cd#
##...##
##.@.##
##...##
#cB#Ab#
#######`},

		{24, `
###############
#d.ABC.#.....a#
######@#@######
###############
######@#@######
#b.....#.....c#
###############`},

		{32, `
#############
#DcBa.#.GhKl#
#.###@#@#I###
#e#d#####j#k#
###C#@#@###J#
#fEbA.#.FgHi#
#############`},

		{72, `
#############
#g#f.D#..h#l#
#F###e#E###.#
#dCba@#@BcIJ#
#############
#nK.L@#@G...#
#M###N#H###.#
#o#m..#i#jk.#
#############`},
	}

	for i, tt := range tests {
		r := strings.NewReader(strings.TrimSpace(tt.src))
		m, err := rog.Parse(r)
		if err != nil {
			t.Fatal(err)
		}

		got, err := day18b(m)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.cost {
			t.Errorf("test %d got cost %d, want %d", i, got, tt.cost)
		}
	}
}
