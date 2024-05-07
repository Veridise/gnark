package rangecheck

import (
	"github.com/Veridise/gnark/frontend"
	"github.com/Veridise/gnark/std/math/bits"
)

type plainChecker struct {
	api frontend.API
}

func (pl plainChecker) Check(v frontend.Variable, nbBits int) {
	bits.ToBinary(pl.api, v, bits.WithNbDigits(nbBits))
}
