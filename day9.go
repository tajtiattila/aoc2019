package main

import (
	"fmt"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day9() {
	rom := input.MustInts(9)

	fmt.Println("9/1:", intcomp.MustRun(rom, 1))
	fmt.Println("9/2:", intcomp.MustRun(rom, 2))
}
