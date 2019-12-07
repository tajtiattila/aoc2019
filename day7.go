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

	log.Println("day7a:", maxsig(rom, false))
	log.Println("day7b:", maxsig(rom, true))
}

func maxsig(rom []int, feedback bool) int {
	var rmax int
	hasrmax := false

	const nphase = 5
	v := make([]int, nphase)
	for i := range v {
		if feedback {
			v[i] = nphase + i
		} else {
			v[i] = i
		}
	}

	for {
		var r int
		if feedback {
			r = run_feedbackloop_amp(rom, v)
		} else {
			r = runmultiamp(rom, v)
		}
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

func runmultiamp(rom []int, phase []int) int {
	if len(phase) == 0 {
		return 0
	}

	chin, chout := multiamp(rom, phase, nil)

	chin <- 0

	return <-chout
}

func run_feedbackloop_amp(rom []int, phase []int) int {
	if len(phase) == 0 {
		return 0
	}

	chdone := make(chan struct{}, len(phase))
	chin, chout := multiamp(rom, phase, chdone)

	v := 0

Loop:
	for {
		select {
		case chin <- v:
			v = <-chout
		case <-chdone:
			break Loop
		}
	}

	return v
}

func multiamp(rom []int, phase []int, chdone chan<- struct{}) (chin chan<- int, chout <-chan int) {
	chstart := make(chan int)
	var ch <-chan int
	ch = chstart
	for _, ph := range phase {
		ch = amp(rom,
			intcomp.MultiReader(
				intcomp.FixedInput(ph),
				&chanInput{ch},
			),
			chdone)
	}

	return chstart, ch
}

func amp(rom []int, in intcomp.IntReader, chdone chan<- struct{}) <-chan int {
	chout := make(chan int)
	go func() {
		defer close(chout)
		if chdone != nil {
			defer func() {
				chdone <- struct{}{}
			}()
		}
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
