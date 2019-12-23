package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/spacecard"
)

func day22() {
	r := input.MustReader(22)

	ops, err := spacecard.ParseOps(r)
	if err != nil {
		log.Fatal(err)
	}

	const deck1 = 10007
	calc := spacecard.Calc(deck1, ops)
	fmt.Println("22/1:", calc.Index(2019))

	const deck2 = 119315717514047
	const nrept = 101741582076661
	calc = spacecard.Calc(deck2, ops)
	fmt.Println("22/2:", calc.Repeat(nrept).Inv().Index(2020))
}
