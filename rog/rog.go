// Advent of Code 2019, Day 18
package rog

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"math/bits"
)

const (
	wall  = '#'
	space = '.'
)

func IsKey(x byte) bool {
	return 'a' <= x && x <= 'z'
}

type Point struct {
	X, Y int
}

func Pt(x, y int) Point { return Point{x, y} }

type Map struct {
	Dx, Dy int

	P []byte
}

func Parse(r io.Reader) (Map, error) {
	var m Map

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := bytes.TrimSpace(scanner.Bytes())
		if len(b) == 0 {
			continue
		}
		m.P = append(m.P, b...)
		if m.Dx == 0 {
			m.Dx = len(m.P)
		}
		m.Dy++
	}
	if err := scanner.Err(); err != nil {
		return m, err
	}
	if m.Dx*m.Dy != len(m.P) {
		return m, errors.New("dim")
	}
	return m, nil
}

func (m Map) In(p Point) bool {
	x, y := p.X, p.Y
	if x < 0 || y < 0 || x >= m.Dx || y >= m.Dy {
		return false
	}
	return true
}

func (m Map) ofs(p Point) int {
	return p.X + p.Y*m.Dx
}

func (m Map) At(p Point) byte {
	if !m.In(p) {
		return wall
	}
	return m.P[m.ofs(p)]
}

func (m Map) Find(x byte) (p Point, ok bool) {
	for i, b := range m.P {
		if b == x {
			x := i % m.Dx
			y := i / m.Dx
			return Pt(x, y), true
		}
	}
	return Point{}, false
}

func (m Map) FindKeys() []Point {
	var keys []Point
	for i, b := range m.P {
		if 'a' <= b && b <= 'z' {
			x := i % m.Dx
			y := i / m.Dx
			keys = append(keys, Pt(x, y))
		}
	}
	return keys
}

var dirstep = []Point{
	Pt(0, -1),
	Pt(1, 0),
	Pt(0, 1),
	Pt(-1, 0),
}

func (m Map) Move(a Agent, dir int) (Agent, bool) {
	d := dirstep[dir]
	p := Pt(a.P.X+d.X, a.P.Y+d.Y)
	ok := a.Visit(m.At(p))
	if ok {
		a.P = p
	}
	return a, ok
}

type Agent struct {
	P Point // position

	Keys uint32 // bit 0: 'a', bit 1: 'b'...
}

func (a Agent) NKeys() int {
	return bits.OnesCount32(a.Keys)
}

func (a *Agent) AddKey(x byte) {
	if 'a' <= x && x <= 'z' {
		a.Keys |= uint32(1) << uint(x-'a')
	}
}

func (a Agent) HasKey(x byte) bool {
	if 'a' <= x && x <= 'z' {
		m := a.Keys & (uint32(1) << uint(x-'a'))
		return m != 0
	}
	return false
}

func (a Agent) HasDoorKey(x byte) bool {
	return a.HasKey(x - 'A' + 'a')
}

func (a Agent) CanVisit(x byte) bool {
	switch {

	case x == wall:
		return false

	case 'A' <= x && x <= 'Z':
		return a.HasDoorKey(x)

	default:
		return true
	}
}

func (a *Agent) Visit(x byte) bool {
	a.AddKey(x)
	return a.CanVisit(x)
}
