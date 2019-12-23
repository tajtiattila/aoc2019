package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/pluto"
)

func day20() {
	r := input.MustReader(20)

	m, err := pluto.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("20/1:", pluto.ShortestPathLen(m))
	fmt.Println("20/2:", pluto.RecShortestPathLen(m))
}
