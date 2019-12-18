package main

import (
	"errors"
	"log"

	"github.com/tajtiattila/aoc2019/astar"
	"github.com/tajtiattila/aoc2019/rog"
)

func day18() {
	rc := mustdaydata(18)
	defer rc.Close()

	m, err := rog.Parse(rc)
	if err != nil {
		log.Fatal(err)
	}

	cost, err := day18a(m)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("day18a", cost)
}

func day18a(m rog.Map) (int, error) {
	nkeys := len(m.FindKeys())

	p, ok := m.Find('@')
	if !ok {
		return 0, errors.New("no entry")
	}
	start := rog.Agent{P: p}

	goalf := func(p astar.Point, d []astar.Step) []astar.Step {
		a := p.(rog.Agent)
		for _, stop := range rog.NextStops(a, m, nil) {
			g := stop.Agent
			d = append(d, astar.Step{
				P:            g,
				Cost:         stop.Cost,
				EstimateLeft: nkeys - g.NKeys(),
			})
		}
		return d
	}

	_, cost := astar.FindPath(start, goalf)
	return cost, nil
}
