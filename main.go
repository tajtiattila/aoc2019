package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/tajtiattila/aoc2019/input"
)

var cli struct {
	verbose bool
}

func main() {
	flag.BoolVar(&cli.verbose, "v", false, "verbose mode")
	clearCache := flag.Bool("cc", false, "clear cache")
	flag.BoolVar(&input.IgnoreCache, "ic", false, "ignore cache")
	flag.Parse()

	input.ReportDownload = cli.verbose

	if *clearCache {
		if err := input.ClearCache(); err != nil {
			log.Fatalln("clear cache:", err)
		}
	}

	want := make(map[int]struct{})
	for _, a := range flag.Args() {
		if ok := parsearg(want, a); !ok {
			log.Fatalf("invalid argument %q", a)
		}
	}

	for i, f := range regfn {
		if f == nil {
			continue
		}
		if _, ok := want[i]; ok || len(want) == 0 {
			f()
		}
	}
}

func init() {
	regd(1, day1)
	regd(2, day2)
	regd(3, day3)
	regd(4, day4)
	regd(5, day5)
	regd(6, day6)
	regd(7, day7)
	regd(8, day8)
	regd(9, day9)
	regd(10, day10)
	regd(11, day11)
	regd(12, day12)
	regd(13, day13)
	regd(14, day14)
	regd(15, day15)
	regd(16, day16)
	regd(17, day17)
	regd(18, day18)
	regd(19, day19)
	regd(20, day20)
	regd(21, day21)
	regd(22, day22)
	regd(23, day23)
	regd(24, day24)
	regd(25, day25)
}

var regfn []func()

func regd(d int, f func()) {
	if d >= len(regfn) {
		x := make([]func(), d+1)
		copy(x, regfn)
		regfn = x
	}

	regfn[d] = f
}

func parsearg(m map[int]struct{}, arg string) bool {
	var a, b int
	var aerr, berr error
	if strings.HasSuffix(arg, "+") {
		a, aerr = strconv.Atoi(strings.TrimSuffix(arg, "+"))
		b = len(regfn)
	} else if i := strings.Index(arg, ".."); i >= 0 {
		a, aerr = strconv.Atoi(arg[:i])
		b, berr = strconv.Atoi(arg[i+2:])
	} else {
		a, aerr = strconv.Atoi(arg)
		b = a
	}
	if aerr != nil || berr != nil {
		return false
	}
	for i := a; i <= b; i++ {
		m[i] = struct{}{}
	}
	return true
}
