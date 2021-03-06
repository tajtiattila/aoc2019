package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/tajtiattila/aoc2019/input"
)

func day12() {
	r := input.MustReader(12)

	v, err := scanjovian(r)
	if err != nil {
		log.Fatal(err)
	}

	w := make([]ob3, len(v))
	copy(w, v)
	fmt.Println("12/1:", simjovian(w, 1000))

	copy(w, v)
	fmt.Println("12/2:", repjovian(w))
}

type pt3 [3]int

func parsept3(s string) (pt3, error) {
	var p pt3
	_, err := fmt.Sscanf(s, "<x=%d, y=%d, z=%d>", &p[0], &p[1], &p[2])
	return p, err
}

func (p pt3) mag() int {
	var v int
	for k := 0; k < 3; k++ {
		a := p[k]
		if a < 0 {
			a = -a
		}
		v += a
	}
	return v
}

func packi(i int) uint64 {
	if i >= 0 {
		return 2 * uint64(i)
	} else {
		return 1 + 2*uint64(-i)
	}
}

func (p pt3) pack(shift uint) uint64 {
	return (packi(p[0]) << (2 * shift)) |
		(packi(p[1]) << shift) |
		packi(p[2])
}

func (p pt3) String() string {
	return fmt.Sprintf("<x=%d, y=%d, z=%d>", p[0], p[1], p[2])
}

type ob3 struct {
	p, v pt3
}

func (o ob3) energy() int {
	return o.p.mag() * o.v.mag()
}

func (o ob3) String() string {
	return fmt.Sprintf("pos=%v, vel=%v", o.p, o.v)
}

func scanjovian(r io.Reader) ([]ob3, error) {
	var v []ob3
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		p, err := parsept3(scanner.Text())
		if err != nil {
			return nil, err
		}
		v = append(v, ob3{p: p})
	}
	return v, nil
}

func simjovian(v []ob3, nstep int) (energy int) {
	for i := 0; i < nstep; i++ {
		stepjovian(v)
	}

	for _, o := range v {
		energy += o.energy()
	}
	return energy
}

func stepjovian(v []ob3) {
	for k := 0; k < 3; k++ {
		stepjovianc(v, k)
	}
}

func stepjovianc(v []ob3, k int) {
	for i := range v {
		a := &v[i]
		w := v[i+1:]
		for j := range w {
			b := &w[j]
			if a.p[k] != b.p[k] {
				if a.p[k] < b.p[k] {
					a.v[k]++
					b.v[k]--
				} else {
					a.v[k]--
					b.v[k]++
				}
			}
		}
	}

	for i := range v {
		a := &v[i]
		a.p[k] += a.v[k]
	}
}

func repjovian(v []ob3) int64 {
	if len(v) <= 4 {
		return repjovian4(v)
	}
	panic("not implemented")
}

type jov4 [8]int16

func jov4c(k int, v []ob3) jov4 {
	var j jov4
	for i := range v {
		j[2*i] = int16(v[i].p[k])
		j[2*i+1] = int16(v[i].v[k])
	}
	return j
}

func repjovian4(v []ob3) int64 {
	seen := make([]map[jov4]int, 3)
	for i := range seen {
		seen[i] = make(map[jov4]int)
	}

	m := make([]int, 3)
	done := 0

	w := make([]ob3, len(v))
	copy(w, v)
	n := 0
	for done != 3 {
		for k := 0; k < 3; k++ {
			if m[k] != 0 {
				continue
			}
			state := jov4c(k, w)
			if _, ok := seen[k][state]; ok {
				m[k] = n
				done++
			}
			seen[k][state] = n
		}
		stepjovian(w)
		n++
	}

	r := int64(1)
	for _, v := range m {
		r = lcm64(r, int64(v))
	}
	return r
}

func gcd64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm64(a, b int64) int64 {
	return a / gcd64(a, b) * b
}
