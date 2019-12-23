package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day11() {
	rom := input.MustInts(11)

	bot := new(paintbot)
	c := intcomp.New(rom, bot, bot)
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("11/1:", len(bot.hull))

	bot = new(paintbot)
	bot.paint(1)
	c = intcomp.New(rom, bot, bot)
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("11/2:\n%s", bot.render())
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
	return render(b.hull, func(c int) rune {
		if c == 0 {
			return '░'
		} else {
			return '█'
		}
	})
}
