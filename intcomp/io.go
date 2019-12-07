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
	if i.p < len(i.v) {
		v := i.v[i.p]
		i.p++
		return v, nil
	}
	return 0, io.EOF
}

type multiReader struct {
	v []IntReader
	i int
}

func (r *multiReader) ReadInt() (int, error) {
	for r.i < len(r.v) {
		v, err := r.v[r.i].ReadInt()
		if err != io.EOF {
			return v, err
		}
		r.i++
	}
	return 0, io.EOF
}

func MultiReader(r ...IntReader) IntReader {
	return &multiReader{v: r}
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
