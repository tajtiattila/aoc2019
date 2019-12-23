package main

import (
	"fmt"

	"github.com/tajtiattila/aoc2019/input"
)

func day1() {
	ints := input.MustInts(1)

	sum1 := 0
	sum2 := 0
	for _, n := range ints {
		sum1 += fuel(n)
		sum2 += fuelx(n)
	}

	fmt.Println("1/1:", sum1)
	fmt.Println("1/2:", sum2)
}

func fuel(n int) int {
	f := (n / 3) - 2
	if f > 0 {
		return f
	}
	return 0
}

func fuelx(n int) int {
	x := 0
	for n > 0 {
		m := fuel(n)
		x += m
		n = m
	}
	return x
}
