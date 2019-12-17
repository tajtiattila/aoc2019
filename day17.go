package main

import (
	"log"

	"github.com/tajtiattila/aoc2019/ascmap"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day17() {
	rom := mustdaydataInts(17)

	img, err := intcomp.Run(rom)
	if err != nil {
		log.Fatal(err)
	}

	am := ascmap.FromInts(img)

	log.Println("day17a:", am.AlignParam())

	rom[0] = 2
	r, err := intcomp.Run(rom, am.Scafprog()...)
	if err != nil {
		log.Fatal(err)
	}
	_, dust := ascmap.SplitOutput(r)

	log.Println("day17b:", dust)
}
