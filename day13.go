package main

import (
	"log"

	"github.com/tajtiattila/aoc2019/intcomp"
)

func day13() {
	rom := mustdaydataInts(13)

	var ac arcabi
	ca := intcomp.New(rom, nil, intcomp.CallFuncOutput(ac.draw))
	if err := ca.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("13a", ac.countTile(2))
}

type arcabi struct {
	tile map[point]int
}

func (ac *arcabi) draw(v []int) (int, error) {
	if len(v) < 3 {
		return 0, nil
	}
	x, y, tile := v[0], v[1], v[2]
	if ac.tile == nil {
		ac.tile = make(map[point]int)
	}
	ac.tile[pt(x, y)] = tile
	return 3, nil
}

func (ac *arcabi) countTile(t int) int {
	n := 0
	for _, x := range ac.tile {
		if x == t {
			n++
		}
	}
	return n
}
