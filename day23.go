package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/network"
)

func day23() {
	rom := input.MustInts(23)

	nw := network.New(rom, 50)
	err := nw.Run(func(p network.Packet) bool {
		if p.Addr == 255 {
			fmt.Println("23/1:", p.Y)
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
	fmt.Println("23/2:", lastv)
}
