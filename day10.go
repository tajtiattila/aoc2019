package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

func day10() {
	rc := mustdaydata(10)
	defer rc.Close()

	k, err := loadastchart(rc)
	if err != nil {
		log.Fatal(err)
	}

	c, score := k.findbest()
	log.Println("day10a:", score)

	z, ok := k.zap(c, 200)
	if !ok {
		log.Fatal("zap")
	}
	log.Println("day10b:", z.x*100+z.y)
}

type rat struct {
	num, denom int
}

func findfracs(maxdenom int) []rat {
	var v []rat
	for denom := 2; denom < maxdenom; denom++ {
		for num := 1; num < denom; num++ {
			if gcd(denom, num) == 1 {
				v = append(v, rat{num, denom})
			}
		}
	}
	sort.Slice(v, func(i, j int) bool {
		ir, jr := v[i], v[j]
		return ir.num*jr.denom < jr.num*ir.denom
	})
	return v
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func dirfracs(m int) []point {
	rv := findfracs(m)

	var v []point
	addp := func(q, a, b int) {
		switch q {
		case 0:
			v = append(v, pt(a, -b))
		case 1:
			v = append(v, pt(b, a))
		case 2:
			v = append(v, pt(-a, b))
		case 3:
			v = append(v, pt(-b, -a))
		}
	}

	for q := 0; q < 4; q++ {
		addp(q, 0, 1)
		for _, r := range rv {
			addp(q, r.num, r.denom)
		}
		addp(q, 1, 1)
		for i := len(rv) - 1; i >= 0; i-- {
			r := rv[i]
			addp(q, r.denom, r.num)
		}
	}

	return v
}

type astchart struct {
	dx, dy int

	p []byte // dx*dy bytes

	rays []point
}

func loadastchart(r io.Reader) (astchart, error) {
	scanner := bufio.NewScanner(r)
	var lines [][]byte
	for scanner.Scan() {
		var line []byte
		for _, r := range strings.TrimSpace(scanner.Text()) {
			switch r {
			case '#':
				line = append(line, 1)
			case '.':
				line = append(line, 0)
			default:
				return astchart{}, fmt.Errorf("invalid rune %c", r)
			}
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return astchart{}, err
	}

	dy := len(lines)
	dx := 0
	for _, line := range lines {
		if len(line) > dx {
			dx = len(line)
		}
	}

	p := make([]byte, 0, dx*dy)
	for _, line := range lines {
		p = append(p, line...)
	}
	if len(p) != dx*dy {
		return astchart{}, fmt.Errorf("invalid format")
	}

	m := dx
	if dy > m {
		m = dy
	}

	return astchart{
		dx:   dx,
		dy:   dy,
		p:    p,
		rays: dirfracs(m),
	}, nil
}

func (z astchart) copy() astchart {
	p := make([]byte, len(z.p))
	copy(p, z.p)
	return astchart{
		dx:   z.dx,
		dy:   z.dy,
		p:    p,
		rays: z.rays,
	}
}

func (z astchart) In(p point) bool {
	if p.x < 0 || p.x >= z.dx {
		return false
	}
	if p.y < 0 || p.y >= z.dy {
		return false
	}
	return true
}

func (z astchart) At(p point) byte {
	if !z.In(p) {
		return 0
	}
	return z.p[p.x+p.y*z.dx]
}

func (k *astchart) findbest() (p point, score int) {
	for y := 0; y < k.dy; y++ {
		for x := 0; x < k.dx; x++ {
			if k.At(pt(x, y)) == 0 {
				continue
			}
			if n := k.nvis(pt(x, y)); n > score {
				p, score = pt(x, y), n
			}
		}
	}
	return p, score
}

func (k *astchart) nvis(p point) int {
	n := 0

	for _, r := range k.rays {
		if _, ok := k.ray(p, r); ok {
			n++
		}
	}

	return n
}

func (k *astchart) ray(p, r point) (vis point, ok bool) {
	p.x += r.x
	p.y += r.y
	for k.In(p) {
		if k.At(p) != 0 {
			return p, true
		}
		p.x += r.x
		p.y += r.y
	}
	return p, false
}

func (k *astchart) zap(station point, nbet int) (p point, ok bool) {
	n := 0
	for {
		zapped := false
		for _, r := range k.rays {
			p, ok := k.ray(station, r)
			if ok {
				n++
				if n == nbet {
					return p, true
				}
				zapped = true
			}
		}
		if !zapped {
			return station, false
		}
	}
}
