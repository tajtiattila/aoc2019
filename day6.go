package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/tajtiattila/aoc2019/input"
)

func day6() {
	r := input.MustReader(6)

	m, err := parseOrbits(r)
	if err != nil {
		log.Fatal("Parsing orbits:", err)
	}

	na := 0
	for _, b := range m {
		na += b.depth()
	}
	fmt.Println("6/1:", na)

	fmt.Println("6/2:", orbitPathFind(m, "YOU", "SAN")-2)
}

func parseOrbits(r io.Reader) (map[string]*body, error) {
	m := make(map[string]*body)

	o := func(name string) *body {
		if b, ok := m[name]; ok {
			return b
		}
		b := &body{name: name}
		m[name] = b
		return b
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		p, m, err := parseOrbit(scanner.Text())
		if err != nil {
			return nil, err
		}
		o(m).parent = o(p)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func parseOrbit(s string) (planet, moon string, err error) {
	ts := strings.TrimSpace
	s = ts(s)
	i := strings.Index(s, ")")
	if i < 0 {
		return "", "", fmt.Errorf("invalid orbit %q", s)
	}
	planet, moon = ts(s[:i]), ts(s[i+1:])
	return planet, moon, nil
}

func orbitPathFind(m map[string]*body, as, bs string) int {
	a, b := m[as], m[bs]
	if a == nil || b == nil {
		panic("invalid names")
	}

	n := 0
	for a != b {
		switch {
		case a.parentof(b):
			n++
			b = b.parent
		case b.parentof(a):
			n++
			a = a.parent
		default:
			n += 2
			a = a.parent
			b = b.parent
		}

		if a == nil || b == nil {
			panic("impossible")
		}
	}
	return n
}

type body struct {
	name string

	parent *body
}

func (b *body) depth() int {
	n := 0
	for b.parent != nil {
		b = b.parent
		n++
	}
	return n
}

func (b *body) parentof(c *body) bool {
	for ; c != nil; c = c.parent {
		if b == c.parent {
			return true
		}
	}
	return false
}
