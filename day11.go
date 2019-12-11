package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/tajtiattila/aoc2019/intcomp"
)

func day11() {
	rom := mustdaydataInts(11)

	bot := new(paintbot)
	c := intcomp.New(rom, bot, bot)
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("day11a:", len(bot.hull))

	bot = new(paintbot)
	bot.paint(1)
	c = intcomp.New(rom, bot, bot)
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
	log.Printf("day11b:\n%s", bot.render())
}

type paintbot struct {
	hull map[point]int

	pos point
	dir int // 0: north, 1: east...

	noutput int // output calls
}

func (b *paintbot) ReadInt() (n int, err error) {
	return b.hull[b.pos], nil
}

func (b *paintbot) WriteInt(n int) error {
	paint := b.noutput%2 == 0
	b.noutput++
	if paint {
		return b.paint(n)
	} else {
		return b.turn(n)
	}
}

func (b *paintbot) paint(n int) error {
	if b.hull == nil {
		b.hull = make(map[point]int)
	}
	b.hull[b.pos] = n
	return nil
}

func (b *paintbot) turn(n int) error {
	switch n {
	case 0:
		b.dir = (b.dir + 3) % 4
	case 1:
		b.dir = (b.dir + 1) % 4
	default:
		return fmt.Errorf("Invalid input: %v", n)
	}
	switch b.dir {
	case 0:
		b.pos.y -= 1
	case 1:
		b.pos.x += 1
	case 2:
		b.pos.y += 1
	case 3:
		b.pos.x -= 1
	default:
		panic("impossible")
	}
	return nil
}

func (b *paintbot) render() string {
	var x0, x1, y0, y1 int
	for p := range b.hull {
		if p.x < x0 {
			x0 = p.x
		}
		if x1 < p.x {
			x1 = p.x
		}
		if p.y < y0 {
			y0 = p.y
		}
		if y1 < p.y {
			y1 = p.y
		}
	}
	dx := x1 - x0 + 1
	dy := y1 - y0 + 1
	v := make([]int, dx*dy)
	for p, col := range b.hull {
		o := (p.x - x0) + (p.y-y0)*dx
		v[o] = col
	}

	var sb strings.Builder
	for i, col := range v {
		var r rune
		if col == 0 {
			r = '░'
		} else {
			r = '█'
		}
		sb.WriteRune(r)
		if (i+1)%dx == 0 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
