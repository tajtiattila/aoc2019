package main

import (
	"fmt"
	"log"
	"strconv"
)

func day4() {
	var v []int
	mustprocstr(mustdaydatastr(4), "-", 2, func(p string) error {
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
	log.Println("day4a:", a)
	log.Println("day4b:", b)
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
