package main

import (
	"log"
	"strings"
)

func day8() {
	im := spaceimage(mustdaydatastr(8), 25, 6)

	var n0min, n12 int
	for i, layer := range im {
		s := layer.pixStat()
		if i == 0 || s[0] < n0min {
			n0min, n12 = s[0], s[1]*s[2]
		}
	}
	log.Println("day8a:", n12)
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
