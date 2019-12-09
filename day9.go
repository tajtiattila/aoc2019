package main

import (
	"log"

	"github.com/tajtiattila/aoc2019/intcomp"
)

func day9() {
	rom := mustdaydataInts(9)

	log.Println("day9a:", intcomp.MustRun(rom, 1))
}
