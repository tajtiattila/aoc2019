package ascmap

import (
	"bytes"
	"fmt"
	"strings"
)

type Ascmap struct {
	dx, dy int
	m      []byte
}

func FromInts(v []int) Ascmap {
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
	return Ascmap{dx, dy, m}
}

func (am Ascmap) atp(p point) byte {
	return am.At(p.x, p.y)
}

func (am Ascmap) At(x, y int) byte {
	if x < 0 || x >= am.dx || y < 0 || y >= am.dy {
		return '.'
	}
	return am.m[y*am.dx+x]
}

func (am Ascmap) intersectv(x, y int) int {
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

func (am Ascmap) AlignParam() int {
	p := 0
	for y := 0; y < am.dy; y++ {
		for x := 0; x < am.dx; x++ {
			p += am.intersectv(x, y)
		}
	}
	return p
}

func (am Ascmap) String() string {
	var sb strings.Builder
	for y := 0; y < am.dy; y++ {
		o := y * am.dx
		sb.Write(am.m[o : o+am.dx])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func ascdir(b byte) int {
	switch b {
	case '^':
		return 0
	case '>':
		return 1
	case 'v':
		return 2
	case '<':
		return 3
	default:
		return -1
	}
}

func (am Ascmap) Scafprog() []int {
	return tricmd(am.scafpimp()).ints()
}

func (am Ascmap) scafpimp() []string {
	for y := 0; y < am.dy; y++ {
		for x := 0; x < am.dx; x++ {
			d := ascdir(am.At(x, y))
			if d >= 0 {
				bot := ascbot{
					am: am,
					p:  pt(x, y),
					d:  d,
				}
				bot.run()
				return bot.cmds
			}
		}
	}
	panic("impossible")
}

type ascbot struct {
	am Ascmap

	p point
	d int

	cmds []string
}

func (b *ascbot) run() {
	start := true

	for {
		ok := b.turn(start)
		start = false

		if !ok {
			return
		}

		b.fwd()
	}
}

var dirdelta = []point{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

func (b *ascbot) turn(start bool) bool {
	ignore := -1
	if !start {
		ignore = (b.d + 2) % 4 // 180°
	}
	nd := b.contdir(ignore)
	if nd == -1 {
		// can't continue
		return false
	}
	right := (nd + 4 - b.d) % 4
	switch right {
	case 1:
		b.cmd("R")
	case 2:
		b.cmd("R", "R")
	case 3:
		b.cmd("L")
	}
	b.d = nd
	return true
}

func (b *ascbot) fwd() {
	d := dirdelta[b.d]
	p := b.p
	i := 0
	for b.am.atp(addpt(p, d)) == '#' {
		i++
		p = addpt(p, d)
	}
	b.p = p
	b.cmd(fmt.Sprint(i))
}

func (b *ascbot) contdir(ignore int) int {
	for i, d := range dirdelta {
		if i != ignore && b.am.atp(addpt(b.p, d)) == '#' {
			return i
		}
	}
	return -1
}

func (b *ascbot) cmd(s ...string) {
	b.cmds = append(b.cmds, s...)
}

type tripletcmd struct {
	Main    string
	A, B, C string
}

func (c tripletcmd) ints() []int {
	s := fmt.Sprintf("%s\n%s\n%s\n%s\nn\n", c.Main, c.A, c.B, c.C)
	var v []int
	for _, r := range s {
		v = append(v, int(r))
	}
	return v
}

func tricmd(src []string) tripletcmd {
	jx := func(x []string) string {
		return strings.Join(x, ",")
	}
	for i := 2; i < len(src); i++ {
		a := src[:i]
		if b, c, ok := trysplit(src, a); ok {
			m, _ := mainprg(src, a, b, c)
			return tripletcmd{
				Main: jx(m),
				A:    jx(a),
				B:    jx(b),
				C:    jx(c),
			}
		}
	}
	return tripletcmd{
		Main: "A",
		A:    jx(src),
	}
}

func prglenok(v []string) bool {
	return len(strings.Join(v, ",")) <= 20
}

func mainprg(src, a, b, c []string) ([]string, bool) {
	if !prglenok(a) || !prglenok(b) || !prglenok(c) {
		return nil, false
	}

	var m []string
	i := 0
	for i < len(src) {
		if strvpfx(src[i:], a) {
			i += len(a)
			m = append(m, "A")
		} else if strvpfx(src[i:], b) {
			i += len(b)
			m = append(m, "B")
		} else if strvpfx(src[i:], c) {
			i += len(c)
			m = append(m, "C")
		} else {
			return nil, false
		}
	}

	return m, prglenok(m)
}

func trysplit(src, a []string) (b, c []string, ok bool) {
	i := 0
	for strvpfx(src[i:], a) {
		i += len(a)
	}

	for k := i + 2; k < len(src); k++ {
		b = src[i:k]
		if c, ok = trysplitc(src, a, b); ok {
			return b, c, ok
		}
	}

	return nil, nil, false
}

func trysplitc(src, a, b []string) (c []string, ok bool) {
	i := 0
	for {
		if strvpfx(src[i:], a) {
			i += len(a)
		} else if strvpfx(src[i:], b) {
			i += len(b)
		} else {
			break
		}
	}

	j := i + strvfind(src[i:], a)
	if jb := i + strvfind(src[i:], b); jb < j {
		j = jb
	}

	for k := i + 2; k <= j; k++ {
		c = src[i:k]
		if _, ok := mainprg(src, a, b, c); ok {
			return c, ok
		}
	}

	return nil, false
}

func strvfind(v, search []string) int {
	end := len(v) - len(search)
	for i := 0; i <= end; i++ {
		if strvpfx(v[i:], search) {
			return i
		}
	}
	return len(v)
}

func strvpfx(v, prefix []string) bool {
	if len(prefix) > len(v) {
		return false
	}
	for i := range prefix {
		if v[i] != prefix[i] {
			return false
		}
	}
	return true
}

func eqints(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func SplitOutput(v []int) (msg string, dust int) {
	n := len(v) - 1
	var sb strings.Builder
	for _, r := range v[:n] {
		sb.WriteByte(byte(r))
	}
	return sb.String(), v[n]
}
