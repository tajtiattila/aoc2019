package main

import (
	"bytes"
	"log"
	"strings"

	"github.com/tajtiattila/aoc2019/intcomp"
)

func day17() {
	rom := mustdaydataInts(17)

	img, err := intcomp.Run(rom)
	if err != nil {
		log.Fatal(err)
	}

	am := ascmapInts(img)

	log.Println("day17a:", am.alignparam())
}

type ascmap struct {
	dx, dy int
	m      []byte
}

func ascmapInts(v []int) ascmap {
	b := make([]byte, len(v))
	for i, n := range v {
		b[i] = byte(n)
	}

	lines := bytes.Split(b, []byte("\n"))
	dy := len(lines)
	var dx int
	for _, l := range lines {
		if len(l) > dx {
			dx = len(l)
		}
	}

	m := make([]byte, dx*dy)
	for y, l := range lines {
		o0 := y * dx
		copy(m[o0:], l)
		o1 := o0 + dx
		o0 += len(l)
		for ; o0 != o1; o0++ {
			m[o0] = '.'
		}
	}
	return ascmap{dx, dy, m}
}

func (am ascmap) At(x, y int) byte {
	if x < 0 || x >= am.dx || y < 0 || y >= am.dy {
		return '.'
	}
	return am.m[y*am.dx+x]
}

func (am ascmap) intersectv(x, y int) int {
	dv := []point{
		pt(0, 0),
		pt(-1, 0),
		pt(1, 0),
		pt(0, -1),
		pt(0, 1),
	}

	for _, d := range dv {
		if am.At(x+d.x, y+d.y) != '#' {
			return 0
		}
	}

	return x * y
}

func (am ascmap) alignparam() int {
	p := 0
	for y := 0; y < am.dy; y++ {
		for x := 0; x < am.dx; x++ {
			p += am.intersectv(x, y)
		}
	}
	return p
}

func (am ascmap) String() string {
	var sb strings.Builder
	for y := 0; y < am.dy; y++ {
		o := y * am.dx
		sb.Write(am.m[o : o+am.dx])
		sb.WriteByte('\n')
	}
	return sb.String()
}
