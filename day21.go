package main

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/tajtiattila/aoc2019/input"
	"github.com/tajtiattila/aoc2019/intcomp"
)

func day21() {
	rom := input.MustInts(21)

	// (!A || !B || !C) && D
	ra := r21(rom, "^aj ^bt |tj ^ct |tj &dj", "walk")
	fmt.Println("21/1:", ra)

	/*
		t21(rom, "^aj ^ct |tj &dj", "run")
		t21(rom, "^aj ^ct &ht |tj &dj", "run")
		t21(rom, "^aj ^ct &ht |tj ^bt &ht |tj &dj", "run")
	*/

	//  ABCDEFGHI
	// @##.#.#..##.###        ^aj ^ct |tj &dj
	// @.####.#..###          ^aj ^ct &ht |tj &dj
	// @#.####.#..###
	// @##.####.#..###
	// ok: ^aj ^ct &ht |tj ^bt &ht |tj &dj
	rb := r21(rom, "^aj ^ct &ht |tj ^bt &ht |tj &dj", "run")
	fmt.Println("21/2:", rb)
}

func t21(rom []int, src, fin string) {
	d := springDroid{prg: springScriptComp(src, fin)}
	c := intcomp.New(rom, &d, &d)
	c.Run()
	if d.result > 0 {
		fmt.Println(d.result)
		return
	}
	//fmt.Println(d.output.String())
	fmt.Printf("%-20s   %s\n", d.failpat(), src)
	if d.failpat() == "" {
		fmt.Println(d.output.String())
	}
}

func r21(rom []int, src, fin string) string {
	d := springDroid{prg: springScriptComp(src, fin)}
	c := intcomp.New(rom, &d, &d)
	c.Run()
	if d.result > 0 {
		return fmt.Sprint(d.result)
	}
	return d.failpat()
}

func springScriptComp(src, fin string) string {
	var b strings.Builder
	for _, x := range strings.Fields(src) {
		if len(x) != 3 {
			log.Fatal(x)
		}
		var cmd string
		switch x[0] {
		case '^':
			cmd = "NOT"
		case '&':
			cmd = "AND"
		case '|':
			cmd = "OR"
		default:
			log.Fatal(x)
		}
		b.WriteString(fmt.Sprintf("%s %c %c\n", cmd, x[1], x[2]))
	}
	b.WriteString(fin)
	b.WriteByte('\n')
	return strings.ToUpper(b.String())
}

type springDroid struct {
	prg string
	r   int

	output strings.Builder

	result int
}

func (d *springDroid) ReadInt() (int, error) {
	r, n := utf8.DecodeRuneInString(d.prg[d.r:])
	d.r += n
	return int(r), nil
}

func (d *springDroid) WriteInt(n int) error {
	if d.r == 0 {
		return nil
	}

	if n > 127 {
		d.result = n
	} else {
		d.output.WriteRune(rune(n))
	}

	return nil
}

func (d *springDroid) failpat() string {
	lines := strings.Split(strings.TrimSpace(d.output.String()), "\n")
	n := 0
	var line, ofs int
	for i, l := range lines {
		if l == "" {
			n = i
		}

		j := strings.IndexByte(l, '@')
		if i == n+2 && j >= 0 {
			line, ofs = n+4, j
			break
		}
		if i == n+3 && j >= 0 {
			line, ofs = n+4, j+1
		}
	}
	if line >= len(lines) {
		return ""
	}
	return lines[line][ofs:]
}
