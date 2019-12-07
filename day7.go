package main

import (
	"io"
	"log"

	"github.com/tajtiattila/aoc2019/intcomp"
)

func day7() {
	rom, err := daydataInts(7)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("day7a:", maxsig(rom))
}

func maxsig(rom []int) int {
	var rmax int
	hasrmax := false

	v := []int{0, 1, 2, 3, 4}
	for {
		r := multiamp(rom, v)
		if !hasrmax || r > rmax {
			rmax = r
			hasrmax = true
		}
		if !perm(v) {
			break
		}
	}

	return rmax
}

func multiamp(rom []int, phase []int) int {
	if len(phase) == 0 {
		return 0
	}

	var ch <-chan int
	for i, ph := range phase {
		var in intcomp.IntReader
		if i == 0 {
			in = intcomp.FixedInput(ph, 0)
		} else {
			in = intcomp.MultiReader(intcomp.FixedInput(ph), &chanInput{ch})
		}
		ch = amp(rom, in)
	}

	return <-ch
}

func amp(rom []int, in intcomp.IntReader) <-chan int {
	chout := make(chan int)
	go func() {
		defer close(chout)
		c := intcomp.New(rom, in, &chanOutput{chout})
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	return chout
}

type chanInput struct {
	ch <-chan int
}

func (i *chanInput) ReadInt() (int, error) {
	v, ok := <-i.ch
	if !ok {
		return 0, io.EOF
	}
	return v, nil
}

type chanOutput struct {
	ch chan<- int
}

func (o *chanOutput) WriteInt(n int) error {
	o.ch <- n
	return nil
}

func perm(v []int) bool {
	if len(v) < 2 {
		return false // no more permutations
	}

	// 1. find i := max(k) where v[k] < v[k+1]
	// 2. find j := max(l) where v[k] < v[l]
	// 3. swap v[i], v[j]
	// 4. reverse v[i+1:]

	// 1. find max(i) where v[i] < v[i+1]
	i := len(v) - 1
	for i > 0 && v[i-1] >= v[i] {
		i--
	}
	if i == 0 {
		return false
	}
	i--

	// 2. find max(j) where v[i] < v[j]
	j := len(v) - 1
	for v[j] < v[i] {
		j--
	}

	// 3. swap v[i], v[j]
	v[i], v[j] = v[j], v[i]

	// 4. reverse v[i+1:]
	for m, n := i+1, len(v)-1; m < n; m, n = m+1, n-1 {
		v[m], v[n] = v[n], v[m]
	}

	return true
}
