package fields_bw6761

import (
	"math/big"

	bw6761 "github.com/consensys/gnark-crypto/ecc/bw6-761"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
)

type curveF = emulated.Field[emulated.BW6761Fp]
type baseEl = emulated.Element[emulated.BW6761Fp]

type E3 struct {
	A0, A1, A2 baseEl
}

type Ext3 struct {
	api frontend.API
	fp  *curveF
}

func NewExt3(api frontend.API) *Ext3 {
	fp, err := emulated.NewField[emulated.BW6761Fp](api)
	if err != nil {
		panic(err)
	}
	return &Ext3{
		api: api,
		fp:  fp,
	}
}

func (e Ext3) Reduce(x *E3) *E3 {
	var z E3
	z.A0 = *e.fp.Reduce(&x.A0)
	z.A1 = *e.fp.Reduce(&x.A1)
	z.A2 = *e.fp.Reduce(&x.A2)
	return &z
}

func (e Ext3) Zero() *E3 {
	zero := e.fp.Zero()
	return &E3{
		A0: *zero,
		A1: *zero,
		A2: *zero,
	}
}

func (e Ext3) One() *E3 {
	one := e.fp.One()
	zero := e.fp.Zero()
	return &E3{
		A0: *one,
		A1: *zero,
		A2: *zero,
	}
}

func (e Ext3) Neg(x *E3) *E3 {
	a0 := e.fp.Neg(&x.A0)
	a1 := e.fp.Neg(&x.A1)
	a2 := e.fp.Neg(&x.A2)
	return &E3{
		A0: *a0,
		A1: *a1,
		A2: *a2,
	}
}

func (e Ext3) Add(x, y *E3) *E3 {
	a0 := e.fp.Add(&x.A0, &y.A0)
	a1 := e.fp.Add(&x.A1, &y.A1)
	a2 := e.fp.Add(&x.A2, &y.A2)
	return &E3{
		A0: *a0,
		A1: *a1,
		A2: *a2,
	}
}

func (e Ext3) Sub(x, y *E3) *E3 {
	a0 := e.fp.Sub(&x.A0, &y.A0)
	a1 := e.fp.Sub(&x.A1, &y.A1)
	a2 := e.fp.Sub(&x.A2, &y.A2)
	return &E3{
		A0: *a0,
		A1: *a1,
		A2: *a2,
	}
}

func (e Ext3) Double(x *E3) *E3 {
	two := big.NewInt(2)
	a0 := e.fp.MulConst(&x.A0, two)
	a1 := e.fp.MulConst(&x.A1, two)
	a2 := e.fp.MulConst(&x.A2, two)
	return &E3{
		A0: *a0,
		A1: *a1,
		A2: *a2,
	}
}

func mulFpByNonResidue(fp *curveF, x *baseEl) *baseEl {

	z := fp.Neg(x)
	z = fp.MulConst(z, big.NewInt(4))
	return z
}

func (e Ext3) Conjugate(x *E3) *E3 {
	a1 := e.fp.Neg(&x.A1)
	return &E3{
		A0: x.A0,
		A1: *a1,
		A2: x.A2,
	}
}

func (e Ext3) MulByElement(x *E3, y *baseEl) *E3 {
	a0 := e.fp.Mul(&x.A0, y)
	a1 := e.fp.Mul(&x.A1, y)
	a2 := e.fp.Mul(&x.A2, y)
	z := &E3{
		A0: *a0,
		A1: *a1,
		A2: *a2,
	}
	return z
}

func (e Ext3) MulByConstElement(x *E3, y *big.Int) *E3 {
	a0 := e.fp.MulConst(&x.A0, y)
	a1 := e.fp.MulConst(&x.A1, y)
	a2 := e.fp.MulConst(&x.A2, y)
	return &E3{
		A0: *a0,
		A1: *a1,
		A2: *a2,
	}
}

// MulBy01 multiplication by sparse element (c0,c1,0)
func (e Ext3) MulBy01(z *E3, c0, c1 *baseEl) *E3 {

	a := e.fp.Mul(&z.A0, c0)
	b := e.fp.Mul(&z.A1, c1)

	tmp := e.fp.Add(&z.A1, &z.A2)
	t0 := e.fp.Mul(c1, tmp)
	t0 = e.fp.Sub(b, t0)
	t0 = e.fp.MulConst(t0, big.NewInt(4))
	t0 = e.fp.Add(t0, a)

	// for t2, schoolbook is faster than karatsuba
	// c2 = a0b2 + a1b1 + a2b0,
	// c2 = a2b0 + b ∵ b2 = 0, b = a1b1
	t2 := e.fp.Mul(&z.A2, c0)
	t2 = e.fp.Add(t2, b)

	t1 := e.fp.Add(c0, c1)
	tmp = e.fp.Add(&z.A0, &z.A1)
	t1 = e.fp.Mul(t1, tmp)
	tmp = e.fp.Add(b, a)
	t1 = e.fp.Sub(t1, tmp)

	return &E3{
		A0: *t0,
		A1: *t1,
		A2: *t2,
	}
}

// MulBy1 multiplication of E6 by sparse element (0, c1, 0)
func (e Ext3) MulBy1(z *E3, c1 *baseEl) *E3 {

	b := e.fp.Mul(&z.A1, c1)
	tmp := e.fp.Add(&z.A1, &z.A2)
	t0 := e.fp.Mul(c1, tmp)
	t0 = e.fp.Sub(b, t0)
	t0 = e.fp.MulConst(t0, big.NewInt(4))
	tmp = e.fp.Add(&z.A0, &z.A1)
	t1 := e.fp.Mul(c1, tmp)
	t1 = e.fp.Sub(t1, b)

	return &E3{
		A0: *t0,
		A1: *t1,
		A2: *b,
	}
}

// MulBy12 multiplication by sparse element (0,b1,b2)
func (e Ext3) MulBy12(x *E3, b1, b2 *baseEl) *E3 {
	t1 := e.fp.Mul(&x.A1, b1)
	t2 := e.fp.Mul(&x.A2, b2)
	c0 := e.fp.Add(&x.A1, &x.A2)
	tmp := e.fp.Add(b1, b2)
	c0 = e.fp.Mul(c0, tmp)
	c0 = e.fp.Sub(c0, t1)
	c0 = e.fp.Sub(t2, c0)
	c0 = e.fp.MulConst(c0, big.NewInt(4))
	c1 := e.fp.Add(&x.A0, &x.A1)
	c1 = e.fp.Mul(c1, b1)
	c1 = e.fp.Sub(c1, t1)
	tmp = mulFpByNonResidue(e.fp, t2)
	c1 = e.fp.Add(c1, tmp)
	tmp = e.fp.Add(&x.A0, &x.A2)
	c2 := e.fp.Mul(b2, tmp)
	c2 = e.fp.Sub(c2, t2)
	c2 = e.fp.Add(c2, t1)
	return &E3{
		A0: *c0,
		A1: *c1,
		A2: *c2,
	}
}

// Mul01By01 multiplies two E3 sparse element of the form:
//
//	E3{
//		A0: c0,
//		A1: c1,
//		A2: 0,
//	}
//
// and
//
//	E3{
//		A0: d0,
//		A1: d1,
//		A2: 0,
//	}
func (e Ext3) Mul01By01(c0, c1, d0, d1 *baseEl) *E3 {
	a := e.fp.Mul(d0, c0)
	b := e.fp.Mul(d1, c1)
	t1 := e.fp.Add(c0, c1)
	tmp := e.fp.Add(d0, d1)
	t1 = e.fp.Mul(t1, tmp)
	tmp = e.fp.Add(b, a)
	t1 = e.fp.Sub(t1, tmp)
	return &E3{
		A0: *a,
		A1: *t1,
		A2: *b,
	}
}

func (e Ext3) Mul(x, y *E3) *E3 {
	return e.MulKaratsuba(x, y)
}

func (e Ext3) MulKaratsuba(x, y *E3) *E3 {
	// Algorithm 13 from https://eprint.iacr.org/2010/354.pdf
	t0 := e.fp.Mul(&x.A0, &y.A0)
	t1 := e.fp.Mul(&x.A1, &y.A1)
	t2 := e.fp.Mul(&x.A2, &y.A2)

	c0 := e.fp.Add(&x.A1, &x.A2)
	tmp := e.fp.Add(&y.A1, &y.A2)
	c0 = e.fp.Mul(c0, tmp)
	tmp = e.fp.Add(t2, t1)
	c0 = e.fp.Sub(c0, tmp)
	c0 = mulFpByNonResidue(e.fp, c0)
	c0 = e.fp.Add(c0, t0)

	c1 := e.fp.Add(&x.A0, &x.A1)
	tmp = e.fp.Add(&y.A0, &y.A1)
	c1 = e.fp.Mul(c1, tmp)
	tmp = e.fp.MulConst(t2, big.NewInt(4))
	tmp = e.fp.Add(tmp, t1)
	tmp = e.fp.Add(tmp, t0)
	c1 = e.fp.Sub(c1, tmp)

	c2 := e.fp.Add(&y.A0, &y.A2)
	tmp = e.fp.Add(&x.A0, &x.A2)
	c2 = e.fp.Mul(c2, tmp)
	c2 = e.fp.Add(c2, t1)
	tmp = e.fp.Add(t2, t0)
	c2 = e.fp.Sub(c2, tmp)

	return &E3{
		A0: *c0,
		A1: *c1,
		A2: *c2,
	}
}

func (e Ext3) MulToomCook3(x, y *E3) *E3 {
	// Toom-Cook-3x:
	// We start by computing five interpolation points – these are evaluations of
	// the product x(u)y(u) with u ∈ {0, ±1, 2, ∞}:
	//
	// v0 = x(0)y(0) = x.A0 * y.A0
	// v1 = x(1)y(1) = (x.A0 + x.A1 + x.A2)(y.A0 + y.A1 + y.A2)
	// v2 = x(−1)y(−1) = (x.A0 − x.A1 + x.A2)(y.A0 − y.A1 + y.A2)
	// v3 = x(2)y(2) = (x.A0 + 2x.A1 + 4x.A2)(y.A0 + 2y.A1 + 4y.A2)
	// v4 = x(∞)y(∞) = x.A2 * y.A2
	//
	// Then the interpolation is performed as:
	//
	// a0 = v0 + β((1/2)v0 − (1/2)v1 − (1/6)v2 + (1/6)v3 − 2v4)
	// a1 = −(1/2)v0 + v1 − (1/3)v2 − (1/6)v3 + 2v4 + βv4
	// a2 = −v0 + (1/2)v1 + (1/2)v2 − v4
	//
	// where is β=-4 the cubic non-residue (mulFpByNonResidue).
	//
	// In-circuit, we compute 6*x*y as
	// a0 = 6v0 - β(3(v1 - v0 + 4v4) + v2 - v3)
	// a1 = -(3v0 + 2v2 + v3) + 6(v1 + 2v4 + βv4)
	// a2 = 3(v1 + v2 - 2(v0 + v4))
	//
	// and then divide a0, a1 and a2 by 6 using a hint.
	//
	// This costs 5M + 22A.

	two := big.NewInt(2)
	three := big.NewInt(3)
	four := big.NewInt(4)
	six := big.NewInt(6)

	v0 := e.fp.Mul(&x.A0, &y.A0)
	t1 := e.fp.Add(&x.A0, &x.A2)
	t2 := e.fp.Add(&y.A0, &y.A2)
	t3 := e.fp.Add(t2, &y.A1)
	v1 := e.fp.Add(t1, &x.A1)
	v1 = e.fp.Mul(v1, t3)
	t3 = e.fp.Sub(t2, &y.A1)
	v2 := e.fp.Sub(t1, &x.A1)
	v2 = e.fp.Mul(v2, t3)
	t1 = e.fp.MulConst(&x.A1, two)
	t2 = e.fp.MulConst(&x.A2, four)
	v3 := e.fp.Add(t1, t2)
	v3 = e.fp.Add(v3, &x.A0)
	t1 = e.fp.MulConst(&y.A1, two)
	t2 = e.fp.MulConst(&y.A2, four)
	t3 = e.fp.Add(t1, t2)
	t3 = e.fp.Add(t3, &y.A0)
	v3 = e.fp.Mul(v3, t3)
	v4 := e.fp.Mul(&x.A2, &y.A2)

	a0 := e.fp.Sub(v1, v0)
	t1 = e.fp.MulConst(v4, four)
	a0 = e.fp.Add(a0, t1)
	a0 = e.fp.MulConst(a0, three)
	a0 = e.fp.Sub(a0, v3)
	a0 = e.fp.Add(a0, v2)
	a0 = e.fp.MulConst(a0, four)
	t1 = e.fp.MulConst(v0, six)
	a0 = e.fp.Add(a0, t1)

	t1 = e.fp.MulConst(v0, three)
	t2 = e.fp.MulConst(v2, two)
	t1 = e.fp.Add(t1, t2)
	t1 = e.fp.Add(t1, v3)
	a1 := e.fp.MulConst(v4, two)
	a1 = e.fp.Sub(v1, a1)
	a1 = e.fp.MulConst(a1, six)
	a1 = e.fp.Sub(a1, t1)

	t1 = e.fp.Add(v0, v4)
	t1 = e.fp.MulConst(t1, two)
	a2 := e.fp.Add(v1, v2)
	a2 = e.fp.Sub(a2, t1)
	a2 = e.fp.MulConst(a2, three)

	return e.divE3By6(
		&E3{A0: *a0, A1: *a1, A2: *a2},
	)
}

func (e Ext3) Square(x *E3) *E3 {
	// Chung-Hasan (SQR2)
	// Algorithm 16 from https://eprint.iacr.org/2010/354.pdf

	c6 := e.fp.MulConst(&x.A1, big.NewInt(2))
	c4 := e.fp.Mul(&x.A0, c6) // x.A0 * xA1 * 2
	c5 := e.fp.Mul(&x.A2, &x.A2)
	c1 := mulFpByNonResidue(e.fp, c5)
	c1 = e.fp.Add(c1, c4)
	c2 := e.fp.Sub(c4, c5)

	c3 := e.fp.Mul(&x.A0, &x.A0)
	c4 = e.fp.Sub(&x.A0, &x.A1)
	c4 = e.fp.Add(c4, &x.A2)
	c5 = e.fp.Mul(c6, &x.A2) // x.A1 * xA2 * 2
	c4 = e.fp.Mul(c4, c4)
	c0 := mulFpByNonResidue(e.fp, c5)
	c4 = e.fp.Add(c4, c5)
	c4 = e.fp.Sub(c4, c3)

	a0 := e.fp.Add(c0, c3)
	a1 := c1
	a2 := e.fp.Add(c2, c4)

	return &E3{
		A0: *a0,
		A1: *a1,
		A2: *a2,
	}
}

func (e Ext3) Inverse(x *E3) *E3 {
	res, err := e.fp.NewHint(inverseE3Hint, 3, &x.A0, &x.A1, &x.A2)
	if err != nil {
		// err is non-nil only for invalid number of inputs
		panic(err)
	}

	inv := E3{
		A0: *res[0],
		A1: *res[1],
		A2: *res[2],
	}
	one := e.One()

	// 1 == inv * x
	_one := e.Mul(&inv, x)
	e.AssertIsEqual(one, _one)

	return &inv

}

func (e Ext3) DivUnchecked(x, y *E3) *E3 {
	res, err := e.fp.NewHint(divE3Hint, 6, &x.A0, &x.A1, &x.A2, &y.A0, &y.A1, &y.A2)
	if err != nil {
		// err is non-nil only for invalid number of inputs
		panic(err)
	}

	div := E3{
		A0: *res[0],
		A1: *res[1],
		A2: *res[2],
	}

	// x = div * y
	_x := e.Mul(&div, y)
	e.AssertIsEqual(x, _x)

	return &div

}

func (e Ext3) divE3By6(x *E3) *E3 {
	res, err := e.fp.NewHint(divE3By6Hint, 3, &x.A0, &x.A1, &x.A2)
	if err != nil {
		// err is non-nil only for invalid number of inputs
		panic(err)
	}

	y := E3{
		A0: *res[0],
		A1: *res[1],
		A2: *res[2],
	}

	// x == 6 * y
	_x := e.MulByConstElement(&y, big.NewInt(6))
	e.AssertIsEqual(x, _x)

	return &y
}

// MulByNonResidue mul x by (0,1,0)
func (e Ext3) MulByNonResidue(x *E3) *E3 {
	z := &E3{
		A0: x.A2,
		A1: x.A0,
		A2: x.A1,
	}
	z.A0 = *mulFpByNonResidue(e.fp, &z.A0)
	return z
}

func (e Ext3) AssertIsEqual(a, b *E3) {
	e.fp.AssertIsEqual(&a.A0, &b.A0)
	e.fp.AssertIsEqual(&a.A1, &b.A1)
	e.fp.AssertIsEqual(&a.A2, &b.A2)
}

func (e Ext3) Copy(x *E3) *E3 {
	return &E3{
		A0: x.A0,
		A1: x.A1,
		A2: x.A2,
	}
}

func FromE3(a *bw6761.E3) E3 {
	return E3{
		A0: emulated.ValueOf[emulated.BW6761Fp](a.A0),
		A1: emulated.ValueOf[emulated.BW6761Fp](a.A1),
		A2: emulated.ValueOf[emulated.BW6761Fp](a.A2),
	}
}

func (e Ext3) Select(selector frontend.Variable, z1, z0 *E3) *E3 {
	a0 := e.fp.Select(selector, &z1.A0, &z0.A0)
	a1 := e.fp.Select(selector, &z1.A1, &z0.A1)
	a2 := e.fp.Select(selector, &z1.A2, &z0.A2)
	return &E3{A0: *a0, A1: *a1, A2: *a2}
}
