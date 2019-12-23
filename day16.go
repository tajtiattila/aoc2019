package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tajtiattila/aoc2019/input"
)

func day16() {
	src := strings.TrimSpace(input.MustString(16))

	//phasemat(20)

	fmt.Println("16/1:", phasesigN(src, 100)[:8])
	fmt.Println("16/2:", phase100tenKmsg(src))
}

func phasesigN(sig string, n int) string {
	sigb := []byte(sig)
	for i := 0; i < n; i++ {
		phasesigb(sigb)
	}
	return string(sigb)
}

func phasesigb(sig []byte) {
	src := make([]int, 0, len(sig))
	for _, b := range sig {
		src = append(src, int(b-'0'))
	}

	for i := range sig {
		w := 0
		for j, v := range src {
			w += v * repat(i, j)
		}
		if w < 0 {
			w = -w
		}
		sig[i] = '0' + byte(w%10)
	}
}

func repat(rpt int, p int) int {
	values := []int{0, 1, 0, -1}
	i := ((p + 1) / (rpt + 1)) % len(values)
	return values[i]
}

func phase100tenKmsg(src string) string {
	const (
		tenK         = 10000
		offsetDigits = 7
		msgLen       = 8
	)
	b := bytes.Repeat([]byte(src), tenK)
	offset := 0
	for _, b := range b[:offsetDigits] {
		offset = offset*10 + int(b-'0')
	}
	if offset < len(b)/2 {
		panic("can't calculate")
	}
	phaseTopHalf100(b)
	m := b[offset:]
	return string(m[:msgLen])
}

func phaseTopHalf100(v []byte) {
	phaseTopHalfN(v, 100)
}

func phaseTopHalfN(v []byte, count int) {
	w := make([]byte, len(v))
	for i := range w {
		w[i] = v[i] - '0'
	}

	h := len(w) / 2
	for n := 0; n < count; n++ {
		sum := 0
		for i := len(w) - 1; i >= h; i-- {
			sum += int(w[i])
			w[i] = byte(sum % 10)
		}
	}

	for i := range w {
		v[i] = w[i] + '0'
	}
}

func phasemat(sz int) {
	m := make([][]int, sz)
	for i := range m {
		v := make([]int, sz)
		for j := range v {
			v[j] = repat(i, j)
		}
		m[i] = v
	}

	for i, v := range m {
		fmt.Printf("%2d %2d\n", i, v)
	}
}
