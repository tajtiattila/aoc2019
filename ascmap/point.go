package ascmap

type point struct {
	x, y int
}

func pt(x, y int) point {
	return point{x, y}
}

func addpt(a, b point) point {
	return pt(a.x+b.x, a.y+b.y)
}
