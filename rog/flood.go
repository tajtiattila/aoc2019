package rog

type Stop struct {
	Agent     // new state (position and keys)
	Cost  int // distance traveled
}

// NextStops returns next stops for the Agent a,
// where new keys may be picked up.
func NextStops(a Agent, m Map, stops []Stop) []Stop {
	if !m.In(a.P) {
		panic("impossible")
	}

	seen := make([]bool, len(m.P))
	seen[m.ofs(a.P)] = true

	cur := []Point{a.P}
	var next []Point
	dist := 0
	for len(cur) != 0 {
		dist++
		for _, c := range cur {
			for _, d := range dirstep {
				n := Pt(c.X+d.X, c.Y+d.Y)
				if o := m.ofs(n); m.In(n) && !seen[o] {
					seen[o] = true
					x := m.At(n)
					if a.CanVisit(x) {
						if IsKey(x) && !a.HasKey(x) {
							s := Agent{n, a.Keys}
							s.AddKey(x)
							stops = append(stops, Stop{s, dist})
						} else {
							next = append(next, n)
						}
					}
				}
			}
		}

		cur, next = next, cur[:0]
	}

	return stops
}
