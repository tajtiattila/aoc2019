package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func day14() {
	rc := mustdaydata(14)
	defer rc.Close()

	m, err := parsefuelrecv(rc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("day14a:", fuelc(m, 1, "FUEL"))

	const trillion = 1000000000000
	log.Println("day14b:", fuelfromore(m, trillion))
}

type fuelctx struct {
	recipe map[string]recfuel

	surplus map[string]int

	ore int
}

func fuelc(m map[string]recfuel, want int, wantu string) int {
	fc := fuelctx{
		recipe:  m,
		surplus: make(map[string]int),
	}
	fc.calc(want, wantu)
	return fc.ore
}

func fuelfromore(m map[string]recfuel, hasore int) int {
	f1 := fuelc(m, 1, "FUEL")

	flo := 0
	fhi := 2 * (hasore + f1 - 1) / f1
	for flo+1 < fhi {
		fm := flo + (fhi-flo)/2
		o := fuelc(m, fm, "FUEL")
		if o <= hasore {
			flo = fm
		} else {
			fhi = fm
		}
	}
	return flo
}

func (fc *fuelctx) calc(want int, chem string) {
	n := fc.surplus[chem]
	if n >= want {
		fc.surplus[chem] -= want
		return
	}

	fc.surplus[chem] = 0
	want -= n

	rc, ok := fc.recipe[chem]
	if !ok {
		if chem == "ORE" {
			fc.ore += want
			return
		} else {
			panic(fmt.Errorf("missing recipe for chem %q", chem))
		}
	}

	batch := rc.o.n
	nb := (want + batch - 1) / batch
	for _, x := range rc.i {
		fc.calc(nb*x.n, x.chem)
	}

	fc.surplus[chem] = nb*batch - want
}

type recelem struct {
	n    int
	chem string
}

type recfuel struct {
	i []recelem
	o recelem
}

func parsefuelrecv(r io.Reader) (map[string]recfuel, error) {
	m := make(map[string]recfuel)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if t == "" {
			continue
		}

		is, os, ok := split2(t, "=>")
		if !ok {
			return nil, fmt.Errorf("parse fuel line %q", t)
		}

		var err error
		var rf recfuel
		rf.o, err = parserecelem(os)
		if err != nil {
			return nil, err
		}

		isv := strings.Split(is, ",")
		for _, x := range isv {
			e, err := parserecelem(x)
			if err != nil {
				return nil, err
			}
			rf.i = append(rf.i, e)
		}

		if _, ok := m[rf.o.chem]; ok {
			return nil, fmt.Errorf("dupe recipe")
		}

		m[rf.o.chem] = rf
	}
	return m, nil
}

func parserecelem(s string) (recelem, error) {
	var e recelem
	_, err := fmt.Sscanf(s, "%d%s", &e.n, &e.chem)
	e.chem = strings.TrimSpace(e.chem)
	return e, err
}

func split2(s, sep string) (before, after string, ok bool) {
	i := strings.Index(s, sep)
	if i < 0 {
		return s, "", false
	}
	return s[:i], s[i+len(sep):], true
}
