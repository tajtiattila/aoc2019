package pluto

import (
	"github.com/tajtiattila/aoc2019/astar"
)

func ShortestPathLen(m *Map) int {
	_, l := ShortestPath(m)
	return l
}
func ShortestPath(m *Map) ([]Point, int) {
	goalf := func(p astar.Point, d []astar.Step) []astar.Step {
		steps := m.Steps(p.(Point))
		for _, p := range steps {
			l := 1
			if p == m.Goal() {
				l = 0
			}
			d = append(d, astar.Step{
				P:    p,
				Cost: 1,

				EstimateLeft: l,
			})
		}
		return d
	}

	xpath, cost := astar.FindPath(m.Start(), goalf)
	path := make([]Point, len(xpath))
	for i := range xpath {
		path[i] = xpath[i].(Point)
	}
	return path, cost
}

func manh(a, b Point) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}
