package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/bits"
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

	state0 := bs.state

	fmt.Println(bs.run1())

	fmt.Println(multibugsim(state0, 200))
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
			case '.', '?':
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

func steprange(state []uint64) (min, max int) {
	for i, v := range state {
		if v != 0 {
			min = i - 1
			break
		}
	}

	for i := min + 1; i < len(state); i++ {
		v := state[i]
		if v != 0 {
			max = i + 1
		}
	}

	return min, max
}

func multibugsim(state0 uint64, nsteps int) int {
	sz := nsteps * 3
	state := make([]uint64, sz)
	tmp := make([]uint64, sz)
	state[sz/2] = state0

	for i := 0; i < nsteps; i++ {
		multibugstep(state, tmp)
		state, tmp = tmp, state
	}

	n := 0
	for _, v := range state {
		n += bits.OnesCount64(v)
	}
	return n
}

func multibugstep(state, nstate []uint64) {
	imin, imax := steprange(state)
	if imin == imax {
		copy(nstate, state)
		return
	}

	if imin < 1 || imax > len(state)-1 {
		panic("overflow")
	}

	for i := range nstate {
		if i < imin || i > imax {
			nstate[i] = 0
			continue
		}

		var s uint64
		for bit, deltas := range bsdelta {
			n := 0
			for _, d := range deltas {
				n += int((state[i+d.dlevel] >> d.bit) & 1)
			}
			if n == 1 || (n == 2 &&
				((state[i]>>bit)&1) == 0) {
				s |= (uint64(1) << bit)
			}
		}
		nstate[i] = s
	}
}

const bsdim = 5

type bsdeltainf struct {
	dlevel int
	bit    int
}

var bsdelta [][]bsdeltainf

func init() {
	const mid = int(bsdim) / 2

	addi := func(p, dlevel, pcond int) {
		if len(bsdelta) != bsdim*bsdim {
			bsdelta = make([][]bsdeltainf, bsdim*bsdim)
		}
		bsdelta[p] = append(bsdelta[p], bsdeltainf{dlevel: dlevel, bit: pcond})
	}

	add := func(x, y, dx, dy int) {
		if x == mid && y == mid {
			return
		}

		src := x + y*bsdim

		nx, ny := x+dx, y+dy
		if nx == mid && ny == mid {
			// inner neighbor
			var ofs, stride int
			if dx == 0 {
				ofs, stride = (mid-dy*2)*bsdim, 1
			} else { // dy == 0
				ofs, stride = mid-dx*2, bsdim
			}
			for i := 0; i < bsdim; i++ {
				addi(src, 1, ofs)
				ofs += stride
			}
			return
		}

		if nx < 0 || ny < 0 || nx >= bsdim || ny >= bsdim {
			// outer neighbor
			ox := mid + dx
			oy := mid + dy
			addi(src, -1, ox+oy*bsdim)
			return
		}

		addi(src, 0, nx+ny*bsdim)
	}

	for y := 0; y < bsdim; y++ {
		for x := 0; x < bsdim; x++ {
			for _, d := range dirdelta {
				add(x, y, d.x, d.y)
			}
		}
	}
}
