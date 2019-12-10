package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func day10() {
	rc := mustdaydata(10)
	defer rc.Close()

	k, err := loadastchart(rc)
	if err != nil {
		log.Fatal(err)
	}

	_, score := k.findbest()
	log.Println("day10a:", score)
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
	return v
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type astchart struct {
	dx, dy int

	p []byte // dx*dy bytes

	rays []rat
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
	rays := findfracs(m)

	return astchart{dx: dx, dy: dy, p: p, rays: rays}, nil
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

	// check the 8 cardinal directions
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x != 0 || y != 0 {
				n += k.nray(p, pt(x, y))
			}
		}
	}

	// check fraction directions
	for _, r := range k.rays {
		n += k.nray(p, pt(r.num, r.denom))
		n += k.nray(p, pt(r.denom, r.num))
		n += k.nray(p, pt(-r.num, r.denom))
		n += k.nray(p, pt(-r.denom, r.num))
		n += k.nray(p, pt(r.num, -r.denom))
		n += k.nray(p, pt(r.denom, -r.num))
		n += k.nray(p, pt(-r.num, -r.denom))
		n += k.nray(p, pt(-r.denom, -r.num))
	}

	return n
}

func (k *astchart) nray(p, r point) int {
	p.x += r.x
	p.y += r.y
	for k.In(p) {
		if k.At(p) != 0 {
			return 1
		}
		p.x += r.x
		p.y += r.y
	}
	return 0
}
