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

	first := true
	var lastv int
	nw = network.New(rom, 50)
	for {
		v, err := nw.NAT()
		if err != nil {
			log.Fatal(err)
		}
		if !first && v == lastv {
			break
		}
		first, lastv = false, v
	}
	log.Println("day23b:", lastv)
}
