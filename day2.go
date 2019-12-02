package main

import (
	"fmt"
	"log"
)

func day2() {
	mem, err := daydataInts(2)
	if err != nil {
		log.Fatal(err)
	}

	mem[1] = 12
	mem[2] = 2
	_, _, err = Intcomp(mem, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("day2:", mem[0])
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
