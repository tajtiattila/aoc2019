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

	c := intcomp.Comp{
		Mem:    copyrom(rom),
		Input:  intcomp.FixedInput(1),
		Output: intcomp.LogOutput("day5a"),
	}
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
