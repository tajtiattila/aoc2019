package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/input"
)

func copyrom(v []int) []int {
	mem := make([]int, len(v))
	copy(mem, v)
	return mem
}

func day2() {
	rom := input.MustInts(2)

	mem := copyrom(rom)
	mem[1] = 12
	mem[2] = 2
	_, _, err := Intcomp(mem, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("2/1:", mem[0])

	const expect = 19690720
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			mem := copyrom(rom)
			mem[1] = noun
			mem[2] = verb
			_, _, err = Intcomp(mem, 0)
			if err != nil {
				log.Fatal(err)
			}
			if expect == mem[0] {
				fmt.Println("2/2:", noun*100+verb)
			}
		}
	}
}

func Intcomp(v []int, pc int) ([]int, int, error) {
	if pc >= len(v) {
		return v, pc, fmt.Errorf("pc after slice")
	}

	var missing int

	var addrerr error
	addr := func(i int) *int {
		if i < len(v) {
			return &v[i]
		}
		if addrerr == nil {
			addrerr = fmt.Errorf("error accessing addr %d > max %d", i, len(v))
		}
		return &missing
	}
	rel := func(i int) *int {
		j := addr(pc + i)
		return addr(*j)
	}

	for {
		switch v[pc] {
		case 1:
			*rel(3) = *rel(1) + *rel(2)
		case 2:
			*rel(3) = *rel(1) * *rel(2)
		case 99:
			return v, pc, nil
		default:
			return v, pc, fmt.Errorf("invalid opcode %v at %d", v[pc], pc)
		}

		if addrerr != nil {
			return v, pc, addrerr
		}

		pc += 4
	}
}
