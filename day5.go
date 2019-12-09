package main

import (
	"log"

	"github.com/tajtiattila/aoc2019/intcomp"
)

func day5() {
	rom, err := daydataInts(5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("day5a:", intcomp.MustRun(rom, 1))
	log.Println("day5b:", intcomp.MustRun(rom, 5))
}
