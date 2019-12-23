package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/tajtiattila/aoc2019/input"
)

func day8() {
	im := spaceimage(input.MustString(8), 25, 6)

	var n0min, n12 int
	for i, layer := range im {
		s := layer.pixStat()
		if i == 0 || s[0] < n0min {
			n0min, n12 = s[0], s[1]*s[2]
		}
	}
	fmt.Println("8/1:", n12)

	var m imagelayer
	for _, layer := range im {
		m.merge(layer)
	}
	fmt.Printf("8/2:\n%s", m)
}

type imagelayer struct {
	stride int
	pix    []byte
}

func (l imagelayer) pixStat() map[byte]int {
	m := make(map[byte]int)
	for _, d := range l.pix {
		m[d]++
	}
	return m
}

func (l *imagelayer) merge(x imagelayer) {
	if l.pix == nil {
		l.stride = x.stride
		l.pix = make([]byte, len(x.pix))
		copy(l.pix, x.pix)
		return
	}

	for i := range l.pix {
		if l.pix[i] == 2 {
			l.pix[i] = x.pix[i]
		}
	}
}

func (l imagelayer) String() string {
	var sb strings.Builder
	for i, d := range l.pix {
		var r rune
		switch d {
		case 0:
			r = '░'
		case 1:
			r = '█'
		case 2:
			r = '-'
		default:
			r = '?'
		}
		sb.WriteRune(r)
		if i%l.stride == l.stride-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func spaceimage(s string, w, h int) []imagelayer {
	b := []byte(strings.TrimSpace(s))
	for i := range b {
		b[i] = b[i] - '0'
	}
	layerSize := w * h // layer size
	if len(b)%layerSize != 0 {
		log.Fatal("source len mod layer size is nonzero")
	}

	var r []imagelayer
	for p := 0; p < len(b); p += layerSize {
		r = append(r, imagelayer{w, b[p : p+layerSize]})
	}
	return r
}
