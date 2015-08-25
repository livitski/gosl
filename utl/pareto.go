// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"math"
	"math/rand"

	"github.com/cpmech/gosl/chk"
)

// DblsParetoMin compares two vectors using Pareto's optimal criterion
//  Note: minimum dominates (is better)
func DblsParetoMin(u, v []float64) (u_dominates, v_dominates bool) {
	chk.IntAssert(len(u), len(v))
	u_has_all_leq := true // all u values are less-than or equal-to v values
	u_has_one_le := false // u has at least one value less-than v
	v_has_all_leq := true // all v values are less-than or equalt-to u values
	v_has_one_le := false // v has at least one value less-than u
	for i := 0; i < len(u); i++ {
		if u[i] > v[i] {
			u_has_all_leq = false
			v_has_one_le = true
		}
		if u[i] < v[i] {
			u_has_one_le = true
			v_has_all_leq = false
		}
	}
	if u_has_all_leq && u_has_one_le {
		u_dominates = true
	}
	if v_has_all_leq && v_has_one_le {
		v_dominates = true
	}
	return
}

// DblsParetoMinProb compares two vectors using Pareto's optimal criterion
// φ ∃ [0,1] is a scaling factor that helps v win even if it's not smaller.
// If φ==0, deterministic analysis is carried out. If φ==1, probabilistic analysis is carried out.
// As φ → 1, v "gets more help".
//  Note: (1) minimum dominates (is better)
//        (2) v dominates if !u_dominates
func DblsParetoMinProb(u, v []float64, φ float64) (u_dominates bool) {
	chk.IntAssert(len(u), len(v))
	var pu, pv float64
	for i := 0; i < len(u); i++ {
		pu += ProbContestSmall(u[i], v[i], φ)
		pv += ProbContestSmall(v[i], u[i], φ)
	}
	pu /= float64(len(u))
	if FlipCoin(pu) {
		u_dominates = true
	}
	return
}

// ProbContestSmall computes the probability for a contest between u and v where u wins if it's
// the smaller value. φ ∃ [0,1] is a scaling factor that helps v win even if it's not smaller.
// If φ==0, deterministic analysis is carried out. If φ==1, probabilistic analysis is carried out.
// As φ → 1, v "gets more help".
func ProbContestSmall(u, v, φ float64) float64 {
	u = math.Atan(u)/math.Pi + 1.5
	v = math.Atan(v)/math.Pi + 1.5
	if u < v {
		return v / (v + φ*u)
	}
	if u > v {
		return φ * v / (φ*v + u)
	}
	return 0.5
}

// FlipCoin generates a Bernoulli variable; throw a coin with probability p
func FlipCoin(p float64) bool {
	if p == 1.0 {
		return true
	}
	if p == 0.0 {
		return false
	}
	if rand.Float64() <= p {
		return true
	}
	return false
}

// ParetoFront2D computes the Pareto optimal front
//  Note: this function is slow for large sets
//  Output:
//   front -- indices of pareto front
func ParetoFront2D(X, Y []float64) (front []int) {
	dominated := map[int]bool{}
	n := len(X)
	for i := 0; i < n; i++ {
		dominated[i] = false
	}
	for i := 0; i < n; i++ {
		u := []float64{X[i], Y[i]}
		for j := i + 1; j < n; j++ {
			v := []float64{X[j], Y[j]}
			u_dominates, v_dominates := DblsParetoMin(u, v)
			if u_dominates {
				dominated[j] = true
			}
			if v_dominates {
				dominated[i] = true
			}
		}
	}
	ndom := 0
	for i := 0; i < n; i++ {
		if !dominated[i] {
			ndom++
		}
	}
	front = make([]int, ndom)
	k := 0
	for i := 0; i < n; i++ {
		if !dominated[i] {
			front[k] = i
			k++
		}
	}
	return
}