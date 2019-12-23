package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day13() {
	rom := input.MustInts(13)

	var ac arcabi
	ca := intcomp.New(rom, ac.input(), intcomp.CallFuncOutput(ac.draw))
	if err := ca.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("13/1:", ac.countTile(2))

	rom[0] = 2
	ac.tile = nil
	cb := intcomp.New(rom, ac.input(), intcomp.CallFuncOutput(ac.draw))
	if err := cb.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("13/2:", ac.score)
}

type arcabi struct {
	tile map[point]int

	ball point
	padx int

	score int
}

func (ac *arcabi) input() intcomp.IntReader {
	return intcomp.IntReaderFunc(func() (int, error) {
		d := 0
		if ac.ball.x != ac.padx {
			if ac.ball.x > ac.padx {
				d = 1
			} else {
				d = -1
			}
		}
		return d, nil
	})
}

func (ac *arcabi) draw(v []int) (int, error) {
	if len(v) < 3 {
		return 0, nil
	}
	x, y, tile := v[0], v[1], v[2]
	if x < 0 {
		if x != -1 || y != 0 {
			return 3, fmt.Errorf("unexpected")
		}
		ac.score = tile
	} else {
		if ac.tile == nil {
			ac.tile = make(map[point]int)
		}
		ac.tile[pt(x, y)] = tile
		switch tile {
		case 3: // pad
			ac.padx = x
		case 4: // ball
			ac.ball = pt(x, y)
		}
	}
	return 3, nil
}

func (ac *arcabi) countTile(t int) int {
	n := 0
	for _, x := range ac.tile {
		if x == t {
			n++
		}
	}
	return n
}

func (ac *arcabi) render() string {
	s := fmt.Sprintln("Score:", ac.score)
	return s + render(ac.tile, func(t int) rune {
		switch t {
		case 0:
			return ' '
		case 1: // wall
			return '█'
		case 2: // block
			return '░'
		case 3: // horzpad
			return '-'
		case 4: // ball
			return '●'
		default:
			return '?'
		}
	})
}
