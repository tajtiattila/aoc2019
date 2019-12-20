package pluto

import (
	"bufio"
	"fmt"
	"io"
)

type Point struct {
	X, Y int
}

func pt(x, y int) Point {
	return Point{x, y}
}

func add(a, b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
}

type Map struct {
	Dx, Dy int

	p []byte

	start, goal Point

	// portal entry/exit positions
	portal map[[2]byte][]warp

	// warps tunnels
	warp map[warp]Point
}

type warp struct {
	p   Point
	dir int // 0:north, 1:east, 2:south, 3:west
}

func (m *Map) Start() Point {
	return m.start
}

func (m *Map) Goal() Point {
	return m.goal
}

func (m *Map) Steps(at Point) []Point {
	if m.at(at.X, at.Y) != '.' {
		return nil
	}

	var dst []Point
	for i, d := range dirs {
		if next := add(at, d); m.at(next.X, next.Y) == '.' {
			dst = append(dst, next)
		} else if pt, ok := m.warp[warp{at, i}]; ok {
			dst = append(dst, pt)
		}

	}
	return dst
}

type RecPoint struct {
	X, Y, Z int
}

func (r RecPoint) Pt() Point {
	return Point{r.X, r.Y}
}

func addr(a RecPoint, b Point) RecPoint {
	return RecPoint{a.X + b.X, a.Y + b.Y, a.Z}
}

func (m *Map) RecStart() RecPoint {
	return RecPoint{m.start.X, m.start.Y, 0}
}

func (m *Map) RecGoal() RecPoint {
	return RecPoint{m.goal.X, m.goal.Y, 0}
}

func (m *Map) warpDelta(p RecPoint) int {
	outer := p.X == 2 || p.Y == 2 || p.X == m.Dx-3 || p.Y == m.Dy-3
	if outer {
		if p.Z == 0 {
			return 0 // wall
		}
		return -1
	} else {
		return 1
	}
}

func (m *Map) RecSteps(p RecPoint) []RecPoint {
	if m.at(p.X, p.Y) != '.' {
		return nil
	}

	var dst []RecPoint
	for i, d := range dirs {
		if next := addr(p, d); m.at(next.X, next.Y) == '.' {
			dst = append(dst, next)
		} else if pt, ok := m.warp[warp{p.Pt(), i}]; ok {
			if d := m.warpDelta(p); d != 0 {
				next := RecPoint{pt.X, pt.Y, p.Z + d}
				dst = append(dst, next)
			}
		}
	}
	return dst
}

func Parse(r io.Reader) (*Map, error) {
	var dx int
	var lines [][]byte
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := make([]byte, len(scanner.Bytes()))
		copy(l, scanner.Bytes())
		lines = append(lines, l)
		if len(l) > dx {
			dx = len(l)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	dy := len(lines)
	p := make([]byte, dx*dy)
	for i := range p {
		p[i] = ' '
	}
	for y, l := range lines {
		copy(p[y*dx:], l)
	}

	m := &Map{
		Dx: dx,
		Dy: dy,
		p:  p,
	}

	// find portals
	for y := 2; y < m.Dy-2; y++ {
		for x := 2; x < m.Dx-2; x++ {
			m.addPortals(x, y)
		}
	}

	// check start&exit
	z := Point{}
	if m.start == z || m.goal == z {
		return nil, fmt.Errorf("start/exit missing")
	}

	// warp tunnels
	for k, v := range m.portal {
		if len(v) != 2 {
			return nil, fmt.Errorf("invalid portal %q: %v", k, v)
		}

		if m.warp == nil {
			m.warp = make(map[warp]Point)
		}

		m.warp[v[0]] = v[1].p
		m.warp[v[1]] = v[0].p
	}

	return m, nil
}

var dirs = []Point{
	pt(0, -1),
	pt(1, 0),
	pt(0, 1),
	pt(-1, 0),
}

func (m *Map) addPortals(x, y int) {
	if m.at(x, y) != '.' {
		return
	}

	for i, d := range dirs {
		p0 := pt(x, y)
		p1 := add(p0, d)
		p2 := add(p1, d)
		c1 := m.at(p1.X, p1.Y)
		c2 := m.at(p2.X, p2.Y)
		if d.X < 0 || d.Y < 0 {
			c1, c2 = c2, c1
		}
		if 'A' <= c1 && c1 <= 'Z' && 'A' <= c2 && c2 <= 'Z' {
			if c1 == c2 {
				switch c1 {
				case 'A':
					m.start = p0
					continue
				case 'Z':
					m.goal = p0
					continue
				}
			}
			pk := [2]byte{c1, c2}
			if m.portal == nil {
				m.portal = make(map[[2]byte][]warp)
			}
			w := m.portal[pk]
			w = append(w, warp{p0, i})
			m.portal[pk] = w
		}
	}
}

func (m *Map) In(x, y int) bool {
	return 0 <= x && x < m.Dx &&
		0 <= y && y < m.Dy
}

func (m *Map) at(x, y int) byte {
	if !m.In(x, y) {
		return ' '
	}
	return m.p[x+y*m.Dx]
}
