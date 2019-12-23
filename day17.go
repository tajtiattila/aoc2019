package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/ascmap"
	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day17() {
	rom := input.MustInts(17)

	img, err := intcomp.Run(rom)
	if err != nil {
		log.Fatal(err)
	}

	am := ascmap.FromInts(img)

	fmt.Println("17/1:", am.AlignParam())

	rom[0] = 2
	r, err := intcomp.Run(rom, am.Scafprog()...)
	if err != nil {
		log.Fatal(err)
	}
	_, dust := ascmap.SplitOutput(r)

	fmt.Println("17/2:", dust)
}
