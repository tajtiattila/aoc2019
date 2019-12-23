package main

import (
	"fmt"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day5() {
	rom := input.MustInts(5)

	fmt.Println("5/1:", intcomp.MustRun(rom, 1))
	fmt.Println("5/2:", intcomp.MustRun(rom, 5))
}
