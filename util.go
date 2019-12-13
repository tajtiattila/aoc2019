package main

import "strings"

func render(m map[point]int, tilef func(int) rune) string {
	var x0, x1, y0, y1 int
	for p := range m {
		if p.x < x0 {
			x0 = p.x
		}
		if x1 < p.x {
			x1 = p.x
		}
		if p.y < y0 {
			y0 = p.y
		}
		if y1 < p.y {
			y1 = p.y
		}
	}
	dx := x1 - x0 + 1
	dy := y1 - y0 + 1
	v := make([]int, dx*dy)
	for p, col := range m {
		o := (p.x - x0) + (p.y-y0)*dx
		v[o] = col
	}

	var sb strings.Builder
	for i, col := range v {
		sb.WriteRune(tilef(col))
		if (i+1)%dx == 0 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
