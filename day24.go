package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func day24() {
	rc := mustdaydata(24)
	defer rc.Close()

	bs, err := parsebugsim(rc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bs)
	fmt.Println(bs.run1())
}

type bugsim struct {
	dx, dy int

	state uint64
}

func (bs bugsim) String() string {
	var sb strings.Builder
	for y := 0; y < bs.dy; y++ {
		for x := 0; x < bs.dx; x++ {
			if bs.at(pt(x, y)) > 0 {
				sb.WriteRune('#')
			} else {
				sb.WriteRune('.')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (bs bugsim) ofs(p point) int {
	x, y := p.x, p.y
	if x < 0 || y < 0 || x >= bs.dx || y >= bs.dy {
		return -1
	}
	return x + bs.dx*y
}

func (bs bugsim) at(p point) int {
	o := bs.ofs(p)
	if o >= 0 {
		return int((bs.state >> o) & 1)
	}
	return 0
}

var dirdelta = []point{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

func (bs *bugsim) step() {
	var nstate uint64
	bit := 0
	for y := 0; y < bs.dy; y++ {
		for x := 0; x < bs.dx; x++ {
			n := 0
			for _, d := range dirdelta {
				n += bs.at(pt(x+d.x, y+d.y))
			}
			var v uint64
			if bs.at(pt(x, y)) == 1 {
				if n == 1 {
					v = 1
				}
			} else {
				if n == 1 || n == 2 {
					v = 1
				}
			}
			nstate |= v << bit
			bit++
		}
	}
	bs.state = nstate
}

func (bs *bugsim) run1() uint64 {
	seen := make(map[uint64]struct{})
	for {
		//fmt.Printf("%s\n", bs)
		bs.step()
		if _, ok := seen[bs.state]; ok {
			return bs.state
		}
		seen[bs.state] = struct{}{}
	}
}

func parsebugsim(r io.Reader) (*bugsim, error) {
	scanner := bufio.NewScanner(r)
	bit := 0
	var dx, dy int
	var state uint64
	for scanner.Scan() {
		n := len(scanner.Text())
		if n > dx {
			dx = n
		}
		for _, c := range scanner.Text() {
			switch c {
			case '#':
				state |= uint64(1) << bit
			case '.':
				// pass
			default:
				return nil, fmt.Errorf("invalid rune %c", c)
			}
			bit++
		}
		dy++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if dx*dy != bit || bit > 64 {
		return nil, fmt.Errorf("invalid layout: %dx%d/%d", dx, dy, bit)
	}
	return &bugsim{
		dx:    dx,
		dy:    dy,
		state: state,
	}, nil
}
