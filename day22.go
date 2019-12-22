package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

func day22() {
	rc := mustdaydata(22)
	defer rc.Close()

	deck := factoryDeck(10007)
	err := deck.apply(rc)
	if err != nil {
		log.Fatal(err)
	}

	// 6625 too low
	log.Println("day22a:", deck.index(2019))
}

type shuffler struct {
	v, tmp []int
}

func factoryDeck(n int) shuffler {
	v := make([]int, n)
	for i := range v {
		v[i] = i
	}
	return shuffler{v: v}
}

func (s *shuffler) apply(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := s.instr(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (r *shuffler) instr(s string) error {
	if s == "deal into new stack" {
		r.deal()
		return nil
	}

	const pdi = "deal with increment "
	if strings.HasPrefix(s, pdi) {
		i, err := strconv.Atoi(strings.TrimPrefix(s, pdi))
		if err == nil {
			r.inc(i)
		}
		return err
	}

	const pci = "cut "
	if strings.HasPrefix(s, pci) {
		i, err := strconv.Atoi(strings.TrimPrefix(s, pci))
		if err == nil {
			r.cut(i)
		}
		return err
	}

	return fmt.Errorf("unknown command %q", s)
}

func (r *shuffler) index(n int) int {
	for i, x := range r.v {
		if x == n {
			return i
		}
	}
	return -1
}

func (r *shuffler) deal() {
	l := len(r.v)
	h := l / 2
	for i := 0; i < h; i++ {
		j := l - (i + 1)
		r.v[i], r.v[j] = r.v[j], r.v[i]
	}
}

func (r *shuffler) cut(n int) {
	s, d := r.itmp()
	if n < 0 {
		n += len(s)
	}
	m := copy(d, s[n:])
	copy(d[m:], s[:n])
}

func (r *shuffler) inc(n int) {
	if n <= 0 || len(r.v)%n == 0 {
		panic("can't deal inc")
	}

	s, d := r.itmp()
	for i, x := range s {
		j := (i * n) % len(s)
		d[j] = x
	}
}

func (r *shuffler) itmp() (src, dst []int) {
	src = r.v
	if cap(r.tmp) >= len(r.v) {
		dst = r.tmp[:len(r.v)]
	} else {
		dst = make([]int, len(r.v))
	}

	r.v, r.tmp = dst, src
	return src, dst
}
