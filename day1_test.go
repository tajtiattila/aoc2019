package main

import "testing"

func TestFuel(t *testing.T) {
	tt := []struct {
		mod  int
		fuel int
	}{
		{1969, 966},
		{100756, 50346},
	}

	for _, x := range tt {
		f := fuelx(x.mod)
		if x.fuel != f {
			t.Errorf("%d needs fuel %d, got %d", x.mod, x.fuel, f)
		}
	}
}
