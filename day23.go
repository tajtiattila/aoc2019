package main

import (
	"log"

	"github.com/tajtiattila/aoc2019/network"
)

func day23() {
	rom := mustdaydataInts(23)

	nw := network.New(rom, 50)
	err := nw.Run(func(p network.Packet) bool {
		if p.Addr == 255 {
			log.Println("day23a:", p.Y)
			return false
		}
		return true
	})
	if err != nil {
		log.Fatal(err)
	}
}
