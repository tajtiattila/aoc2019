package intcomp

import (
	"io"
	"log"
)

type IntReader interface {
	ReadInt() (int, error)
}

type IntWriter interface {
	WriteInt(int) error
}

func FixedInput(values ...int) IntReader {
	return &fixedInput{v: values}
}

type fixedInput struct {
	v []int // input values
	p int   // position
}

func (i *fixedInput) ReadInt() (int, error) {
	if i.p <= len(i.v) {
		v := i.v[i.p]
		i.p++
		return v, nil
	}
	return 0, io.EOF
}

func LogOutput(prefix string) IntWriter {
	return &logOutput{prefix: prefix}
}

type logOutput struct {
	prefix string
}

func (o *logOutput) WriteInt(v int) error {
	log.Printf("%s: %d", o.prefix, v)
	return nil
}