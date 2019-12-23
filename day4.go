package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tajtiattila/aoc2019/input"
)

func day4() {
	var v []int
	mustprocstr(input.MustString(4), "-", 2, func(p string) error {
		n, err := strconv.Atoi(p)
		v = append(v, n)
		return err
	})

	lo, hi := v[0], v[1]
	a, b := 0, 0
	for pw := lo; pw <= hi; pw++ {
		if pw4(pw, 0) {
			a++
		}
		if pw4(pw, 2) {
			b++
		}
	}
	fmt.Println("4/1:", a)
	fmt.Println("4/2:", b)
}

func pw4(pw, wantrun int) bool {
	var ok bool
	var nrun int
	rundone := func() {
		if wantrun == 0 {
			ok = ok || nrun >= 2
		} else {
			ok = ok || nrun == wantrun
		}
		nrun = 0
	}
	var lastr rune
	for _, r := range fmt.Sprint(pw) {
		if r < lastr {
			return false
		}
		if r != lastr {
			rundone()
		}
		nrun++
		lastr = r
	}
	rundone()
	return ok
}

func mustprocstr(s, sep string, nwant int, f func(string) error) {
	parts := strings.Split(strings.TrimSpace(s), sep)
	if len(parts) != nwant {
		log.Fatalf("error processing %.100q (sep %q), want %d parts, got %d",
			s, sep, nwant, len(parts))
	}
	for i, part := range parts {
		if err := f(part); err != nil {
			log.Fatalf("error processing %.100q (sep %q) at %d: %v",
				s, sep, i, err)
		}
	}
}
