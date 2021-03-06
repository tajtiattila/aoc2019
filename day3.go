package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tajtiattila/aoc2019/input"
)

func day3() {
	r := input.MustReader(3)

	var wires []wire
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		w, err := parsewire(0, 0, scanner.Text())
		if err != nil {
			log.Fatal("bad wire", scanner.Text(), err)
		}
		wires = append(wires, w)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("scan error", err)
	}

	w0 := wires[0]
	w1 := wires[1]
	var im, id int
	for i := range get_intersections(w0.s, w1.s) {
		m := iabs(i.x) + iabs(i.y)
		if m != 0 && (im == 0 || m < im) {
			im = m
		}
		d := w0.d[i] + w1.d[i]
		if d != 0 && (id == 0 || d < id) {
			id = d
		}
	}
	fmt.Println("3/1:", im)
	fmt.Println("3/2:", id)
}

func iabs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}

func get_intersections(w1, w2 []wiresegment) <-chan point {
	ch := make(chan point)
	go func() {
		defer close(ch)
		for _, s1 := range w1 {
			for _, s2 := range w2 {
				s1.intersects(s2, ch)
			}
		}
	}()
	return ch
}

type point struct {
	x, y int
}

func pt(x, y int) point {
	return point{x, y}
}

func addpt(a, b point) point {
	return pt(a.x+b.x, a.y+b.y)
}

type wiresegment struct {
	vert   bool
	u      int // y for horizontal, x for vertical
	v0, v1 int // x for horizontal, y for vertical
}

func wireseg(x0, y0, x1, y1 int) wiresegment {
	if y0 == y1 {
		if x1 < x0 {
			x0, x1 = x1, x0
		}
		return wiresegment{false, y0, x0, x1}
	} else if x0 == x1 {
		if y1 < y0 {
			y0, y1 = y1, y0
		}
		return wiresegment{true, x0, y0, y1}
	}

	panic("impossible")
}

func (a wiresegment) intersects(b wiresegment, ch chan<- point) {
	if a.vert == b.vert {
		if a.u != b.u || b.v1 < a.v0 || a.v1 < b.v0 {
			return
		}
		// overlapping parallel segments
		v0 := a.v0
		if v0 < b.v0 {
			v0 = b.v0
		}
		v1 := a.v1
		if b.v1 < v1 {
			v1 = b.v1
		}
		for v := v0; v <= v1; v++ {
			if a.vert {
				ch <- pt(a.u, v)
			} else {
				ch <- pt(v, a.u)
			}
		}
	}

	if a.v0 <= b.u && b.u <= a.v1 &&
		b.v0 <= a.u && a.u <= b.v1 {

		if a.vert {
			ch <- pt(a.u, b.u)
		} else {
			ch <- pt(b.u, a.u)
		}
	}
}

type wire struct {
	s []wiresegment // segments
	d map[point]int // distance from source
}

func parsewire(x, y int, s string) (wire, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")

	var w []wiresegment
	m := make(map[point]int)
	d := 0
	for i, p := range parts {
		if len(p) < 2 {
			return wire{}, fmt.Errorf("invalid part %s at %d", p, i)
		}
		var dx, dy int
		switch p[0] {
		case 'U':
			dy = 1
		case 'D':
			dy = -1
		case 'R':
			dx = 1
		case 'L':
			dx = -1
		}
		l, err := strconv.Atoi(p[1:])
		if err != nil {
			return wire{}, fmt.Errorf("invalid part %s at %d", p, i)
		}

		sx, sy := x, y
		for i := 0; i < l; i++ {
			x += dx
			y += dy
			d++
			if _, ok := m[pt(x, y)]; !ok {
				m[pt(x, y)] = d
			}
		}
		w = append(w, wireseg(sx, sy, x, y))
	}

	return wire{w, m}, nil
}
