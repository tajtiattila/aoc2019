package spacecard

import "math/big"

type Calculation struct {
	ndeck int
	a, b  int
}

func (c Calculation) fixup() Calculation {
	c.a %= c.ndeck
	if c.a < 0 {
		c.a += c.ndeck
	}
	c.b %= c.ndeck
	if c.b < 0 {
		c.b += c.ndeck
	}
	return c
}

func (c Calculation) Index(n int) int {
	var r big.Int
	r.Mod(r.Add(r.Mul(bi(c.a), bi(n)), bi(c.b)), bi(c.ndeck))
	return int(r.Int64())
}

func bi(i int) *big.Int {
	return big.NewInt(int64(i))
}

func (c Calculation) Repeat(n int) Calculation {
	c = c.fixup()

	ndeck := bi(c.ndeck)

	/*
	 * 1: ax+b
	 * 2: a(ax+b)+b = a²x+ab+b = a²x+(a+1)b
	 * 3: a(a²x+(a+1)b)+b = a³x+a(a+1)b+b = a²x+(a²+a+1)b
	 */

	var a big.Int
	a.Exp(bi(c.a), bi(n), ndeck)

	var aminus1 big.Int
	aminus1.Sub(&a, bi(1))

	var b, mi big.Int
	b.Mul(bi(c.b), &aminus1)
	b.Mul(&b, mi.ModInverse(bi(c.a-1), ndeck))
	b.Mod(&b, ndeck)

	return Calculation{
		ndeck: c.ndeck,
		a:     int(a.Int64()),
		b:     int(b.Int64()),
	}.fixup()
}

func (c Calculation) Inv() Calculation {
	/*
		y=ax+b mod n
		y/a=x+b/a mod n
		x=y/a-b/a mod n
	*/

	ndeck := bi(c.ndeck)

	var a, b big.Int
	a.ModInverse(bi(c.a), ndeck)
	b.Mod(b.Mul(bi(-c.b), &a), ndeck)

	return Calculation{
		ndeck: c.ndeck,
		a:     int(a.Int64()),
		b:     int(b.Int64()),
	}.fixup()
}

func pow(base, exp, mod int) int {
	var r big.Int
	bi := func(i int) *big.Int {
		return big.NewInt(int64(i))
	}
	r.Exp(bi(base), bi(exp), bi(mod))
	return int(r.Int64())
}

func Calc(ndeck int, ops []Op) Calculation {
	c := Calculation{
		ndeck: ndeck,
		a:     1,
	}

	for _, op := range ops {
		op.comp(&c)
	}

	return c.fixup()
}

func (*dealIntoNew) comp(c *Calculation) {
	c.a = -c.a
	c.b = -c.b - 1
}

/* let n = 3
0123456789
9876543210
2109876543
3456789012

0123456789
3456789012
2109876543

(cut n) when reversed == cut -n
*/

func (op *cutop) comp(c *Calculation) {
	c.b -= op.n
}

/* let n = 3
0123456789
9876543210
9258147036
6307418529

0123456789
..0....... startofs = n-1
..0..1..2.
.30.41.52.
6307418529

9258147036
0741852963
3692581470

dwi n when reversed == offset dwi n with startofs = n-1
*/

func (op *dealWithInc) comp(c *Calculation) {
	// c(ax+b) = acx+bc
	c.a *= op.n
	c.b *= op.n

	c.a = c.a % c.ndeck
	c.b = c.b % c.ndeck
}
