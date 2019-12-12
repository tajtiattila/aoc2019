package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func day12() {
	rc, err := daydata(12)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	v, err := scanjovian(rc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("12a:", simjovian(v, 1000))
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
		simjovianstep(v)
	}

	for _, o := range v {
		energy += o.energy()
	}
	return energy
}

func simjovianstep(v []ob3) {
	for i := range v {
		a := &v[i]
		w := v[i+1:]
		for j := range w {
			b := &w[j]
			for k := 0; k < 3; k++ {
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
	}

	for i := range v {
		a := &v[i]
		for k := 0; k < 3; k++ {
			a.p[k] += a.v[k]
		}
	}
}
