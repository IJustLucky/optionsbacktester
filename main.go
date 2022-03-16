package main

import (
	a "alpaca_dynamic_alpha/options"
	"fmt"
)

const (
	interest = 7.9
)

func main() {
	r := a.R(interest)
	ben, eq, S, _, std1, _ := a.Alpaca()
	K, q, t := a.GetOptionsData()
	m := a.Mean(ben)
	fmt.Printf("m %f\n", m)
	c := a.Covariance(ben, eq)
	fmt.Printf("c %f\n", c)
	v := a.Variance(ben)
	fmt.Printf("v %f\n", v)
	b := a.Beta(ben, eq)
	fr := a.R(interest)
	fmt.Printf("fr %f\n", fr)
	fmt.Printf("b %f\n", b)
	alp := a.Alpha(ben, b, fr, eq)
	fmt.Printf("a %f\n", alp)
	delta1, delta2 := a.D(ben)
	fmt.Printf("d1 %f\n", delta1)
	fmt.Printf("d2 %f\n", delta2)
	dc := a.DeltaC(eq)
	fmt.Printf("dc %f\n", dc)
	dp := a.DeltaP(eq)
	fmt.Printf("dp %f\n", dp)
	g := a.Gamma()
	fmt.Printf("g %f\n", g)
	k := a.Kappa(true)
	fmt.Printf("k %f\n", k)
	rc := a.RhoC(eq)
	fmt.Printf("rc %f\n", rc)
	rp := a.RhoP(eq)
	fmt.Printf("rp %f\n", rp)
	iv := a.ImpliedVol(ben)
	fmt.Printf("iv %f\n", iv)
	s := a.Sharpe(alp, fr, v)
	fmt.Printf("s %f\n", s)
	i := a.Intrinsic(true)
	fmt.Printf("i %f\n", i)
	e := a.Extrinsic(true)
	fmt.Printf("e %f\n", e)
	T := a.Theta(eq, true, S, K, t, std1, r, q)
	fmt.Printf("t %f\n", T)
}
