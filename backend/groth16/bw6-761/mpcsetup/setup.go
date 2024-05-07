// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gnark DO NOT EDIT

package mpcsetup

import (
	curve "github.com/consensys/gnark-crypto/ecc/bw6-761"
	"github.com/consensys/gnark-crypto/ecc/bw6-761/fr/fft"
	groth16 "github.com/Veridise/gnark/backend/groth16/bw6-761"
)

func ExtractKeys(srs1 *Phase1, srs2 *Phase2, evals *Phase2Evaluations, nConstraints int) (pk groth16.ProvingKey, vk groth16.VerifyingKey) {
	_, _, _, g2 := curve.Generators()

	// Initialize PK
	pk.Domain = *fft.NewDomain(uint64(nConstraints))
	pk.G1.Alpha.Set(&srs1.Parameters.G1.AlphaTau[0])
	pk.G1.Beta.Set(&srs1.Parameters.G1.BetaTau[0])
	pk.G1.Delta.Set(&srs2.Parameters.G1.Delta)
	pk.G1.Z = srs2.Parameters.G1.Z
	bitReverse(pk.G1.Z)

	pk.G1.K = srs2.Parameters.G1.L
	pk.G2.Beta.Set(&srs1.Parameters.G2.Beta)
	pk.G2.Delta.Set(&srs2.Parameters.G2.Delta)

	// Filter out infinity points
	nWires := len(evals.G1.A)
	pk.InfinityA = make([]bool, nWires)
	A := make([]curve.G1Affine, nWires)
	j := 0
	for i, e := range evals.G1.A {
		if e.IsInfinity() {
			pk.InfinityA[i] = true
			continue
		}
		A[j] = evals.G1.A[i]
		j++
	}
	pk.G1.A = A[:j]
	pk.NbInfinityA = uint64(nWires - j)

	pk.InfinityB = make([]bool, nWires)
	B := make([]curve.G1Affine, nWires)
	j = 0
	for i, e := range evals.G1.B {
		if e.IsInfinity() {
			pk.InfinityB[i] = true
			continue
		}
		B[j] = evals.G1.B[i]
		j++
	}
	pk.G1.B = B[:j]
	pk.NbInfinityB = uint64(nWires - j)

	B2 := make([]curve.G2Affine, nWires)
	j = 0
	for i, e := range evals.G2.B {
		if e.IsInfinity() {
			// pk.InfinityB[i] = true should be the same as in B
			continue
		}
		B2[j] = evals.G2.B[i]
		j++
	}
	pk.G2.B = B2[:j]

	// Initialize VK
	vk.G1.Alpha.Set(&srs1.Parameters.G1.AlphaTau[0])
	vk.G1.Beta.Set(&srs1.Parameters.G1.BetaTau[0])
	vk.G1.Delta.Set(&srs2.Parameters.G1.Delta)
	vk.G2.Beta.Set(&srs1.Parameters.G2.Beta)
	vk.G2.Delta.Set(&srs2.Parameters.G2.Delta)
	vk.G2.Gamma.Set(&g2)
	vk.G1.K = evals.G1.VKK

	// sets e, -[δ]2, -[γ]2
	if err := vk.Precompute(); err != nil {
		panic(err)
	}

	return pk, vk
}
