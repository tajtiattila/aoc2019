package spacecard

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Deck(n int) []int {
	v := make([]int, n)
	for i := range v {
		v[i] = i
	}
	return v
}

func Index(deck []int, n int) int {
	for i, x := range deck {
		if x == n {
			return i
		}
	}
	return -1
}

func Apply(deck []int, ops []Op) {
	var d, t []int
	d = deck
	for _, o := range ops {
		d, t = o.apply(d, t)
	}
	copy(deck, d)
}

func Track(n, ndeck int, ops []Op) int {
	for _, o := range ops {
		n = o.track(n, ndeck)
	}
	return n
}

func ParseOps(r io.Reader) ([]Op, error) {
	var v []Op
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if i, err := instr(scanner.Text()); err != nil {
			return nil, err
		} else {
			v = append(v, i)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return v, nil
}

func instr(s string) (Op, error) {
	if s == "deal into new stack" {
		return &dealIntoNew{}, nil
	}

	const pdi = "deal with increment "
	if strings.HasPrefix(s, pdi) {
		i, err := strconv.Atoi(strings.TrimPrefix(s, pdi))
		return &dealWithInc{i}, err
	}

	const pci = "cut "
	if strings.HasPrefix(s, pci) {
		i, err := strconv.Atoi(strings.TrimPrefix(s, pci))
		return &cutop{i}, err
	}

	return nil, fmt.Errorf("unknown command %q", s)
}

type Op interface {
	apply(deck, tmp []int) (ndeck, ntmp []int)
	track(n, ndeck int) int
	comp(c *Calculation)
}

type dealIntoNew struct{}

func (*dealIntoNew) apply(deck, tmp []int) (ndeck, ntmp []int) {
	l := len(deck)
	h := l / 2
	for i := 0; i < h; i++ {
		j := l - (i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck, tmp
}

func (*dealIntoNew) track(n, ndeck int) int {
	return ndeck - (n + 1)
}

type cutop struct {
	n int
}

func (op *cutop) apply(deck, tmp []int) (ndeck, ntmp []int) {
	deck, tmp = decktmp(deck, tmp)
	n := op.n
	if n < 0 {
		n += len(tmp)
	}
	m := copy(deck, tmp[n:])
	copy(deck[m:], tmp[:n])
	return deck, tmp
}

func (op *cutop) track(n, ndeck int) int {
	n -= op.n
	for n < 0 {
		n += ndeck
	}
	for n >= ndeck {
		n -= ndeck
	}
	return n
}

type dealWithInc struct {
	n int
}

func (op *dealWithInc) apply(deck, tmp []int) (ndeck, ntmp []int) {
	n := op.n
	if n <= 0 || len(deck)%n == 0 {
		panic("can't deal inc")
	}

	deck, tmp = decktmp(deck, tmp)
	for i, x := range tmp {
		j := (i * n) % len(deck)
		deck[j] = x
	}
	return deck, tmp
}

func (op *dealWithInc) track(n, ndeck int) int {
	return (n * op.n) % ndeck
}

func decktmp(deck, tmp []int) (ndeck, ntmp []int) {
	if cap(tmp) >= len(deck) {
		tmp = tmp[:len(deck)]
	} else {
		tmp = make([]int, len(deck))
	}
	return tmp, deck // swap
}
