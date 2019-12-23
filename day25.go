package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day25() {
	rom := input.MustInts(25)

	t := newtermio(rom)
	t.ignoreitems = []string{"giant electromagnet", "infinite loop", "molten lava", "photons", "escape pod"}
	t.track(nil)
	//fmt.Println(t.inv)
	t.run(t.secpath)
	ok := t.shuffle(t.secdoor, t.inv...)
	fmt.Printf("25/1:\n%s\n", t.result(ok))
}

type termio struct {
	c *intcomp.Comp

	room  string
	state int // 1: items, 2: doors
	doors []string
	items []string

	ignoreitems []string

	inv []string

	secpath string // path to Security Checkpoint
	secdoor string // door to enter Pressure-Sensitive Floor

	linebuf []byte

	input []int
	n     int // read index

	output bytes.Buffer
}

func newtermio(rom []int) *termio {
	t := new(termio)
	t.c = intcomp.New(rom, t, t)
	t.c.Run()
	return t
}

func (t *termio) result(truncate bool) string {
	s := t.output.String()
	if truncate {
		i := strings.LastIndex(s, "\n\n")
		if i > 0 {
			s = s[i+2:]
		}
	}
	return s
}

func (t *termio) run(args ...interface{}) {
	cmd := fmt.Sprintln(args...)
	fmt.Fprint(&t.output, "Command> ", cmd)
	t.input = t.input[:0]
	t.n = 0
	for _, r := range cmd {
		t.input = append(t.input, int(r))
	}
	t.c.Run()
}

func (t *termio) shuffle(move string, items ...string) bool {
	startroom := t.room
	imax := 1 << len(items)
	for i := 0; i < imax; i++ {
		for j, item := range items {
			if (i & (1 << j)) == 0 {
				t.run("take", item)
			} else {
				t.run("drop", item)
			}
		}
		//t.run("inv")
		t.run(move)
		if t.room != startroom {
			return true
		}
	}
	return false
}

func backdir(a string) string {
	switch a {
	case "north":
		return "south"
	case "south":
		return "north"
	case "east":
		return "west"
	case "west":
		return "east"
	}
	return a
}

func (t *termio) track(path []string) {

Items:
	for _, item := range t.items {
		for _, x := range t.ignoreitems {
			if x == item {
				continue Items
			}
		}
		t.run("take", item)
		t.inv = append(t.inv, item)
	}

	var back string
	if len(path) > 0 {
		back = backdir(path[len(path)-1])
	}

	for _, dir := range []string{"north", "east", "south", "west"} {
		if dir != back && t.hasdoor(dir) {
			t.run(dir)
			if t.room == "Security Checkpoint" {
				t.secpath = strings.Join(path, "\n") + fmt.Sprintf("\n%s\n", dir)
				for _, x := range t.doors {
					if x != backdir(dir) {
						t.secdoor = x
					}
				}
			} else {
				t.track(append(path, dir))
			}
			t.run(backdir(dir))
		}
	}
}

func (t *termio) hasdoor(s string) bool {
	for _, x := range t.doors {
		if x == s {
			return true
		}
	}
	return false
}

func (t *termio) WriteInt(n int) error {
	if n == '\n' {
		t.line(string(t.linebuf))
		t.linebuf = t.linebuf[:0]
		return nil
	}
	t.linebuf = append(t.linebuf, byte(n))
	return nil
}

func (t *termio) line(l string) {
	s := l
	if strings.HasPrefix(s, "== ") {
		s = strings.TrimPrefix(s, "== ")
		s = strings.TrimSuffix(s, " ==")
		t.room = s
		t.state = 0
		t.items = t.items[:0]
		t.doors = t.doors[:0]
	} else if strings.HasPrefix(s, "- ") {
		what := strings.TrimPrefix(s, "- ")
		switch t.state {
		case 1:
			t.items = append(t.items, what)
		case 2:
			t.doors = append(t.doors, what)
		}
	} else {
		switch s {
		case "Items here:":
			t.state = 1
		case "Doors here lead:":
			t.state = 2
		default:
			t.state = 0
		}
	}

	if l == "Command?" {
		return
	}

	fmt.Fprintln(&t.output, l)
}

func (t *termio) ReadInt() (int, error) {
	if t.n >= len(t.input) {
		return -1, io.EOF
	}

	n := t.n
	t.n++
	return t.input[n], nil
}
