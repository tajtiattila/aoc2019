package astar

import "container/heap"

// Point represents a position in the problem space,
// such as location and additional attributes.
//
// Concrete types of Point must be equality comparable.
type Point interface{}

// StepFunc is used in FindPath to find available steps at position p.
// It appends those steps them to dst, and returns the new slice.
type StepFunc func(p Point, dst []Step) []Step

// Step is returned by StepFunc to report
// next positions, costs and estimates of steps.
type Step struct {
	// P is the position after this step.
	P Point

	// Cost is the cost of this step. It must be positive.
	Cost int

	// EstimateLeft is the estimate of the
	// cost to reach the goal from Point.
	//
	// It must be a minimum estimate; i.e.
	// any possible path should cost at least this amount.
	//
	// A zero estimate means that the goal is reached.
	EstimateLeft int
}

type pointStat struct {
	point Point

	totalCost int // total cost to reach this point
	estimate  int // totalCost + EstimateLeft
}

type statEntry struct {
	from Point // predecessor in path

	totalCost int // total cost to reach this point
}

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

func FindPath(start Point,
	stepf func(p Point, dst []Step) []Step) (path []Point, cost int) {

	m := make(map[Point]statEntry)
	m[start] = statEntry{}

	active := activeHeap{
		pointStat{
			point: start,
		},
	}

	// bestCost is the cost of the best path
	// (smallest cost) found so far
	bestCost := maxInt
	var bestGoal Point

	var steps []Step
	for len(active) > 0 {
		a := heap.Pop(&active).(pointStat)

		steps = stepf(a.point, steps[:0])

		for _, s := range steps {
			if s.Cost <= 0 {
				panic("cost must be positive")
			}

			tc := a.totalCost + s.Cost
			est := tc + s.EstimateLeft
			if est > bestCost {
				// this way we can't get better
				continue
			}

			if s.EstimateLeft == 0 {
				// goal reached
				if tc < bestCost {
					bestCost, bestGoal = tc, s.P
				}
			}

			x, ok := m[s.P]
			if !ok || tc < x.totalCost {
				m[s.P] = statEntry{
					from:      a.point,
					totalCost: tc,
				}

				heap.Push(&active, pointStat{
					point:     s.P,
					totalCost: tc,
					estimate:  est,
				})
			}
		}
	}

	for p := bestGoal; p != nil; p = m[p].from {
		path = append(path, p)
	}

	// reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, bestCost
}

type activeHeap []pointStat

func (h activeHeap) Len() int      { return len(h) }
func (h activeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h activeHeap) Less(i, j int) bool {
	return h[i].estimate < h[j].estimate
}

func (h *activeHeap) Push(x interface{}) {
	*h = append(*h, x.(pointStat))
}

func (h *activeHeap) Pop() interface{} {
	n := len(*h) - 1
	e := (*h)[n]
	*h = (*h)[:n]
	return e
}
