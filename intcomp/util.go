package intcomp

import "log"

func MustRun(rom []int, input ...int) []int {
	v, err := Run(rom, input...)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func Run(rom []int, input ...int) ([]int, error) {
	var o SliceOutput
	c := New(rom, FixedInput(input...), &o)
	err := c.Run()
	return o.O, err
}
