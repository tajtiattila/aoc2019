package intcomp

import "fmt"

type Comp struct {
	Mem []int // computer memory

	Input  IntReader
	Output IntWriter

	IC int // instruction counter
}

// Opcodes
const (
	opAdd  = 1
	opMult = 2

	opInput  = 3
	opOutput = 4

	opJumpIfTrue  = 5
	opJumpIfFalse = 6
	opLessThan    = 7
	opEquals      = 8

	opHalt = 99
)

// Parameter modes
const (
	modePosition  = 0
	modeImmediate = 1
)

type opFunc func(c *Comp, inst int) error

var opm []opFunc

func init() {
	opm = make([]opFunc, 100)
	for i := range opm {
		opm[i] = fInvalid
	}

	for inst, f := range map[int]opFunc{
		opAdd:    fAdd,
		opMult:   fMult,
		opInput:  fInput,
		opOutput: fOutput,

		opJumpIfTrue:  fJumpIfTrue,
		opJumpIfFalse: fJumpIfFalse,
		opLessThan:    fLessThan,
		opEquals:      fEquals,

		opHalt: fHalt,
	} {
		opm[inst] = f
	}
}

func (c *Comp) Run() error {
	for !c.Done() {
		if err := c.Step(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Comp) Step() error {
	if c.IC < 0 || c.IC >= len(c.Mem) {
		return fmt.Errorf("Invalid IC")
	}

	inst := c.Mem[c.IC]
	f := opm[inst%100]
	return f(c, inst)
}

func (c *Comp) Done() bool {
	if c.IC >= len(c.Mem) {
		return false
	}

	return c.Mem[c.IC] == opHalt
}

type oper struct {
	c *Comp

	inst int // current instruction

	err error // first error during execution

	invalid int // write destionation for invalid write
}

func (h *oper) fail(f string, args ...interface{}) {
	if h.err != nil {
		return
	}
	what := fmt.Sprintf(f, args...)
	h.err = fmt.Errorf("%s for instruction %d at %d", what, h.c.IC, h.inst)
}

func (h *oper) ptr(i int) *int {
	if i < len(h.c.Mem) {
		return &h.c.Mem[i]
	}

	h.fail("Invalid memory access")
	return &h.invalid
}

func (h *oper) rarg(n int) int {
	v := *h.ptr(h.c.IC + n)

	mode := digit(h.inst, n+1)
	switch mode {

	case modePosition:
		return *h.ptr(v)

	case modeImmediate:
		return v
	}

	h.fail("Invalid opcode mode %d", mode)
	return 0
}

func (h *oper) warg(n int, val int) {
	if h.err != nil {
		return
	}

	v := *h.ptr(h.c.IC + n)
	mode := digit(h.inst, n+1)
	switch mode {

	case modePosition:
		*h.ptr(v) = val
		return

	}

	h.fail("Invalid opcode mode %d", mode)
}

func (h *oper) finish(n int) error {
	if h.err == nil {
		h.c.IC += n
	}
	return h.err
}

func (h *oper) jump(ic int) error {
	if h.err != nil {
		return h.err
	}

	if ic < 0 || ic >= len(h.c.Mem) {
		return fmt.Errorf("Invalid jump %d to %d at %d", h.inst, ic, h.c.IC)
	}

	h.c.IC = ic
	return nil
}

func (h *oper) cond_jump(leninst int, cond bool, ic int) error {
	if h.err != nil {
		return h.err
	}

	if cond {
		return h.jump(ic)
	}

	return h.finish(leninst)
}

var decexp = []int{1, 10, 100, 1000, 10000, 100000}

func digit(value, digit int) int {
	if digit >= len(decexp) {
		return 0
	}
	return (value / decexp[digit]) % 10
}

func fAdd(c *Comp, inst int) error {
	h := oper{c: c, inst: inst}
	h.warg(3, h.rarg(1)+h.rarg(2))
	return h.finish(4)
}

func fMult(c *Comp, inst int) error {
	h := oper{c: c, inst: inst}
	h.warg(3, h.rarg(1)*h.rarg(2))
	return h.finish(4)
}

func fInput(c *Comp, inst int) error {
	v, err := c.Input.ReadInt()
	if err != nil {
		return err
	}

	h := oper{c: c, inst: inst}
	h.warg(1, v)
	return h.finish(2)
}

func fOutput(c *Comp, inst int) error {
	h := oper{c: c, inst: inst}
	v := h.rarg(1)

	if h.err != nil {
		return h.err
	}

	if err := c.Output.WriteInt(v); err != nil {
		return err
	}

	return h.finish(2)
}

func fJumpIfTrue(c *Comp, inst int) error {
	h := oper{c: c, inst: inst}
	return h.cond_jump(3, h.rarg(1) != 0, h.rarg(2))
}

func fJumpIfFalse(c *Comp, inst int) error {
	h := oper{c: c, inst: inst}
	return h.cond_jump(3, h.rarg(1) == 0, h.rarg(2))
}

func fLessThan(c *Comp, inst int) error {
	h := oper{c: c, inst: inst}
	h.warg(3, bool_int(h.rarg(1) < h.rarg(2)))
	return h.finish(4)
}

func fEquals(c *Comp, inst int) error {
	h := oper{c: c, inst: inst}
	h.warg(3, bool_int(h.rarg(1) == h.rarg(2)))
	return h.finish(4)
}

func bool_int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func fHalt(c *Comp, inst int) error {
	return nil
}

func fInvalid(c *Comp, inst int) error {
	return fmt.Errorf("Invalid instruction %d at %d", inst, c.IC)
}
