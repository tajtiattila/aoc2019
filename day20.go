package main

import (
	"log"

	"github.com/tajtiattila/aoc2019/pluto"
)

func day20() {
	rc := mustdaydata(20)
	defer rc.Close()

	m, err := pluto.Parse(rc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("day20a:", pluto.ShortestPathLen(m))
}
