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

	costb, err := day18b(m)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("day18b", costb)
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

type tetrAgent [4]rog.Agent

func (ta tetrAgent) withMove(i int, a rog.Agent) tetrAgent {
	ta[i] = a
	for j := range ta {
		ta[j].Keys |= a.Keys
	}
	return ta
}

func day18b(m rog.Map) (int, error) {
	m = m.FixStart()

	nkeys := len(m.FindKeys())

	entries := m.FindAll('@')
	if len(entries) != 4 {
		return 0, errors.New("invalid start")
	}

	var start tetrAgent
	for i, p := range entries {
		start[i].P = p
	}

	goalf := func(p astar.Point, d []astar.Step) []astar.Step {
		tag := p.(tetrAgent)
		for i, a := range tag {
			for _, stop := range rog.NextStops(a, m, nil) {
				g := stop.Agent
				d = append(d, astar.Step{
					P:            tag.withMove(i, g),
					Cost:         stop.Cost,
					EstimateLeft: nkeys - g.NKeys(),
				})
			}
		}
		return d
	}

	_, cost := astar.FindPath(start, goalf)
	return cost, nil
}
