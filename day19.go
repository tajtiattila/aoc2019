package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day19() {
	var beam tracbeam

	n := 0
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			if beam.at(x, y) {
				n++
			}
		}
	}

	fmt.Println("19/1:", n)
	x, y := beam.fitbox(100)
	fmt.Println("19/2:", x*10000+y)
}

type span struct {
	lo, hi int
}

type tracbeam struct {
	rom []int

	ln []span
}

func (b *tracbeam) at(x, y int) bool {
	ln := b.scanline(y)
	return ln.lo <= x && x < ln.hi
}

func (b *tracbeam) scanline(y int) span {
	for len(b.ln) <= y {
		b.addline()
	}

	return b.ln[y]
}

func (b *tracbeam) addline() {
	y := len(b.ln)

	last := span{0, 0}

	if y > 0 {
		last = b.ln[y-1]
	}

	x0 := last.lo
	stop := last.hi + 100 + 10*y
	for ; !b.atprog(x0, y); x0++ {
		if x0 > stop {
			noscan := span{last.hi, last.hi}
			b.ln = append(b.ln, noscan)
			return
		}
	}

	x1 := x0
	if last.hi > x1 {
		x1 = last.hi
	}
	for ; b.atprog(x1, y); x1++ {
	}

	b.ln = append(b.ln, span{x0, x1})
}

func (b *tracbeam) atprog(x, y int) bool {
	if b.rom == nil {
		b.rom = input.MustInts(19)
	}

	v, err := intcomp.Run(b.rom, x, y)
	if err != nil {
		log.Fatal(err)
	}

	if len(v) != 1 {
		log.Fatal("i/o")
	}

	return v[0] > 0
}

func (b *tracbeam) hit(p, delta point, maxd int, want bool) point {
	for i := 0; i < maxd; i++ {
		if b.at(p.x, p.y) == want {
			return p
		}
		p = addpt(p, delta)
	}
	log.Fatal("no hit")
	return p
}

func (b *tracbeam) dir(dist int) point {
	m := dist * 100

	y0 := b.hit(pt(dist, 0), pt(0, 1), m, true).y
	y1 := b.hit(pt(dist, y0), pt(0, 1), m, false).y

	x0 := b.hit(pt(0, dist), pt(1, 0), m, true).x
	x1 := b.hit(pt(x0, dist), pt(1, 0), m, false).x

	return pt((x0+x1)/2, (y0+y1)/2)
}

func (b *tracbeam) fitbox(dim int) (x, y int) {
	fby := func(y int) bool {
		_, ok := b.fitboxy(y, dim)
		return ok
	}

	y = 1
	for !fby(y) {
		y *= 2
	}

	lo, hi := y/2, y
	for lo < hi {
		m := lo + (hi-lo)/2
		if fby(m) {
			hi = m
		} else {
			lo = m + 1
		}
	}

	// binary search might have skipped valid positions
	ymin := lo - 10*dim
	if ymin < 0 {
		ymin = 0
	}
	for y = ymin; !fby(y); y++ {
	}

	x, _ = b.fitboxy(y, dim)
	return x, y
}

func (b *tracbeam) fitboxy(y, dim int) (minx int, ok bool) {
	ln := b.scanline(y)
	xmax := ln.hi - dim + 1
	for x := ln.lo; x < xmax; x++ {
		if b.fitboxat(x, y, dim) {
			return x, true
		}
	}
	return ln.lo, false
}

func (b *tracbeam) fitboxat(x, y, dim int) bool {
	for ymax := y + dim; y < ymax; y++ {
		ln := b.scanline(y)
		if x < ln.lo || x+dim > ln.hi {
			return false
		}
	}
	return true
}
