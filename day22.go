package main

import (
	"log"

	"github.com/tajtiattila/aoc2019/spacecard"
)

func day22() {
	rc := mustdaydata(22)
	defer rc.Close()

	ops, err := spacecard.ParseOps(rc)
	if err != nil {
		log.Fatal(err)
	}

	const deck1 = 10007
	calc := spacecard.Calc(deck1, ops)
	log.Println("day22a:", calc.Index(2019))

	const deck2 = 119315717514047
	const nrept = 101741582076661
	calc = spacecard.Calc(deck2, ops)
	log.Println("day22b:", calc.Repeat(nrept).Inv().Index(2020))
}
