package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day15() {
	rom := input.MustInts(15)

	start := agent15{c: intcomp.New(rom, nil, nil)}

	var dest agent15 // oxigen

	fmt.Println("15/1:", walk15(start, func(a agent15) bool {
		dest = a
		return true
	}))

	fmt.Println("15/2:", walk15(dest, func(a agent15) bool {
		return false
	}))
}

type agent15 struct {
	c *intcomp.Comp
	p point // position
}

func walk15(start agent15, fdest func(a agent15) bool) int {
	cur := []agent15{start}

	type dirent struct {
		input int
		d     point // delta
	}
	dirs := []dirent{
		{1, pt(0, -1)},
		{2, pt(0, 1)},
		{3, pt(-1, 0)},
		{4, pt(1, 0)},
	}

	mdist := map[point]int{
		pt(0, 0): 0,
	}
	var next []agent15
	dist := 0
	for len(cur) != 0 {
		dist++
		for _, c := range cur {
			for _, de := range dirs {
				np := c.p
				np.x += de.d.x
				np.y += de.d.y
				if _, ok := mdist[np]; ok {
					continue
				}
				nc := c.c.Fork(nil, nil)
				r, err := intcomp.Step(nc, de.input)
				if err != nil {
					log.Fatal("step:", err)
				}
				if r == 0 {
					mdist[np] = -1
				} else {
					mdist[np] = dist
					a := agent15{c: nc, p: np}
					next = append(next, a)
					if r != 1 {
						if fdest(a) {
							return dist
						}
					}
				}
			}
		}

		cur, next = next, cur[:0]
	}

	return dist - 1
}
