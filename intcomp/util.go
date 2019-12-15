package intcomp

import (
	"errors"
	"log"
)

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

// Step runs c feeding it the inputs, and pausing on output.
// It returns the first output of c.
func Step(c *Comp, inputs ...int) (int, error) {
	in := stepInput{
		v: inputs,
	}
	out := pausingOutput{
		v: make([]int, 1),
	}

	ci, co := c.Input, c.Output
	defer func() {
		c.Input = ci
		c.Output = co
	}()

	c.Input = &in
	c.Output = &out

	err := c.Run()
	switch {
	case err == nil:
		if in.p != len(in.v) {
			err = errors.New("inputs weren't consumed")
		} else {
			err = NoOutput
		}
	case err == PauseOutput:
		err = nil
	}

	return out.v[0], err
}

// InsufficientInput is returned from Step when input wasn't sufficient
var InsufficientInput error = errors.New("insufficient input")

// NoOutput is returned from Step when no output was received
var NoOutput error = errors.New("no output")

type stepInput struct {
	v []int // input values
	p int   // position
}

func (i *stepInput) ReadInt() (int, error) {
	if i.p < len(i.v) {
		v := i.v[i.p]
		i.p++
		return v, nil
	}
	return 0, InsufficientInput
}

type pausingOutput struct {
	v []int // output values
	p int   // position
}

func (o *pausingOutput) WriteInt(n int) error {
	if o.p >= len(o.v) {
		return errors.New("pausing output logic")
	}
	o.v[o.p] = n
	o.p++
	if o.p == len(o.v) {
		return PauseOutput
	}
	return nil
}
