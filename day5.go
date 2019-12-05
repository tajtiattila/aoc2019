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

	run := func(inp int, pfx string) {
		c := intcomp.Comp{
			Mem:    copyrom(rom),
			Input:  intcomp.FixedInput(inp),
			Output: intcomp.LogOutput(pfx),
		}
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}
	}

	run(1, "day5a")
	run(5, "day5b")
}
