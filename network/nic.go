package network

import (
	"fmt"

	"github.com/tajtiattila/aoc2019/intcomp"
)

type Network struct {
	nic []*NIC

	// used in Run()
	handler func(Packet) bool
	done    bool
}

func New(rom []int, n int) *Network {
	nw := &Network{
		nic: make([]*NIC, n),
	}
	for i := range nw.nic {
		nic := &NIC{
			addr:  i,
			queue: []int{i},
		}
		nic.c = intcomp.New(rom,
			&nicInput{nic},
			&nwOutput{nw: nw})
		nw.nic[i] = nic
	}
	return nw
}

// PacketHandler is used in Run to process network messages.
// Network.Run() exits when it returns false.
type PacketHandler func(Packet) bool

func (nw *Network) Run(h func(Packet) bool) error {
	if h == nil {
		return fmt.Errorf("Run() needs handler")
	}

	nw.handler = h
	defer func() {
		nw.handler = nil
	}()

	nw.done = false
	for !nw.done {
		if err := nw.Step(); err != nil {
			return err
		}
	}
	return nil
}

func (nw *Network) NAT() (int, error) {
	var last Packet

	nw.handler = func(p Packet) bool {
		last = p
		return true
	}
	defer func() {
		nw.handler = nil
	}()

	for !nw.allIdle() {
		if err := nw.Step(); err != nil {
			return 0, err
		}
	}

	last.Addr = 0
	nw.send(last)
	return last.Y, nil
}

func (nw *Network) Step() error {
	for _, nic := range nw.nic {
		if err := nic.c.Step(); err != nil {
			return err
		}
	}
	return nil
}

func (nw *Network) allIdle() bool {
	for _, nic := range nw.nic {
		if !nic.idle {
			return false
		}
	}
	return true
}

func (nw *Network) send(p Packet) {
	if nw.handler != nil && !nw.handler(p) {
		nw.done = true
	}
	if p.Addr < len(nw.nic) {
		nw.nic[p.Addr].recv(p)
	}
}

type nwOutput struct {
	nw *Network

	buf []int
}

func (o *nwOutput) WriteInt(n int) error {
	o.buf = append(o.buf, n)
	if len(o.buf) < 3 {
		return nil
	}

	o.nw.send(Packet{
		Addr: o.buf[0],
		X:    o.buf[1],
		Y:    o.buf[2],
	})

	o.buf = o.buf[:0]

	return nil
}

type NIC struct {
	addr int

	c *intcomp.Comp

	queue []int

	idle bool
}

func (nic *NIC) recv(p Packet) {
	if p.Addr != nic.addr {
		panic("impossible")
	}
	nic.queue = append(nic.queue, p.X, p.Y)
	nic.idle = false
}

type nicInput struct {
	nic *NIC
}

func (i *nicInput) ReadInt() (int, error) {
	q := i.nic.queue

	n := len(q)
	if n == 0 {
		i.nic.idle = true
		return -1, nil
	}

	r := q[0]
	copy(q, q[1:])
	i.nic.queue = q[:n-1]
	return r, nil
}

type Packet struct {
	Addr int // recipient address

	X, Y int
}
