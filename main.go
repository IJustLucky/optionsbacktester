package main

import (
	"alpaca_dynamic_alpha/autoStrategy"
	a "alpaca_dynamic_alpha/options"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	interest  = 7.9
	benchmark = "SPY"
	equity    = "TSLA"
)

func main() {
	//for ind := 1; ; ind++ {

	benmark, benPrice, benSTD := a.Alpaca(benchmark)
	eq, equityPrice, equitySTD := a.Alpaca(equity)
	benmarkK, _, benmarkt := a.GetOptionsData(benchmark)
	equityK, _, equityt := a.GetOptionsData(equity)
	Bopt1c := a.Option{
		StrikePrice:      benmarkK,
		TimeToExpiration: benmarkt,
		Type:             "Call",
	}

	Bunder := a.Underlying{
		Symbol:     benchmark,
		Price:      benPrice,
		Volatility: benSTD,
	}

	Eopt1c := a.Option{
		StrikePrice:      equityK,
		TimeToExpiration: equityt,
		Type:             "Call",
	}

	Eunder := a.Underlying{
		Symbol:     equity,
		Price:      equityPrice,
		Volatility: equitySTD,
	}

	r := a.R(interest)
	fmt.Printf("Risk Free Rate %f\n", r)
	fmt.Println()
	c := a.Covariance(benmark, eq)
	fmt.Println()
	fmt.Printf("Covariance %f\n", c)
	Mu1 := a.Mean(benmark)
	v1 := a.Variance(Mu1, benmark)
	fmt.Println("Benchmark~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println()
	fmt.Printf("Variance Equity %f\n", v1)
	fmt.Println()
	Mu2 := a.Mean(benmark)
	v2 := a.Variance(Mu2, benmark)
	fmt.Printf("Variance Benchmark %f\n", v2)
	b1 := a.Beta(benmark, eq)
	fmt.Println()
	fmt.Printf("Beta %f\n", a.Beta(benmark, eq))
	alpha1 := a.Alpha(benmark, b1, r, eq)
	fmt.Println()
	fmt.Printf("Alpha %f\n", alpha1)
	fmt.Println()
	fmt.Printf("Sharpe %f\n", a.Sharpe(alpha1, r, equitySTD))
	fmt.Println()
	benc := a.NewBlackScholes(benchmark, &Bopt1c, &Bunder, r)
	//benp := a.NewBlackScholes(&opt1p, &under, r)
	//fmt.Println()
	fmt.Printf("Call StrikePrice %f\n", benc.StrikePrice)
	fmt.Println()
	//fmt.Printf("Put StrikePrice %f\n", benp.StrikePrice)
	//fmt.Println()
	fmt.Printf("Call D1 %f\n", benc.D1)
	fmt.Println()
	//fmt.Printf("Put D1 %f\n", benp.D1)
	//fmt.Println()
	fmt.Printf("Call D2 %f\n", benc.D2)
	fmt.Println()
	//fmt.Printf("Put D2 %f\n", benp.D2)
	//fmt.Println()
	fmt.Printf("Call Delta %f\n", benc.Delta)
	fmt.Println()
	//fmt.Printf("Put Delta %f\n", benp.Delta)
	//fmt.Println()
	fmt.Printf("Call ImpliedVolatility %f\n", benc.ImpliedVolatility)
	fmt.Println()
	//fmt.Printf("Put %f\n", benp.ImpliedVolatility)
	//fmt.Println()
	fmt.Printf("Call Theta %f\n", benc.Theta)
	fmt.Println()
	//fmt.Printf("Put Theta %f\n", benp.Theta)
	//fmt.Println()
	fmt.Printf("Call Gamma %f\n", benc.Gamma)
	fmt.Println()
	//fmt.Printf("Put Gamma %f\n", benp.Gamma)
	//fmt.Println()
	fmt.Printf("Call Kappa %f\n", benc.Kappa)
	fmt.Println()
	//fmt.Printf("Put Kappa %f\n", benp.Kappa)
	//fmt.Println()
	fmt.Printf("Call Rho %f\n", benc.Rho)
	fmt.Println()
	//fmt.Printf("Put Rho %f\n", benp.Rho)
	//fmt.Println()
	fmt.Printf("Call TheoPrice %f\n", benc.TheoPrice)
	fmt.Println()
	//fmt.Printf("Put TheoPrice %f\n", benp.TheoPrice)
	//fmt.Println()
	fmt.Printf("Call UnderlyingPrice %f\n", benc.UnderlyingPrice)
	fmt.Println()
	//fmt.Printf("Put UnderlyingPrice %f\n", benp.UnderlyingPrice)
	//fmt.Println()
	fmt.Printf("Half Life Sigma %f\n", benc.HalfLifeSTD)
	fmt.Println()
	fmt.Printf("Half Life Theta %f\n", benc.HalfLifeTheta)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	Mu3 := a.Mean(eq)
	v3 := a.Variance(Mu3, eq)
	fmt.Println("Equity~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println()
	fmt.Printf("Variance Equity %f\n", v3)
	Mu4 := a.Mean(benmark)
	v4 := a.Variance(Mu4, benmark)
	fmt.Println()
	fmt.Printf("Variance eqchmark %f\n", v4)
	fmt.Println()
	b2 := a.Beta(eq, benmark)
	fmt.Printf("Beta %f\n", a.Beta(eq, benmark))
	alpha2 := a.Alpha(eq, b2, r, benmark)
	fmt.Println()
	fmt.Printf("Alpha %f\n", alpha2)
	fmt.Println()
	fmt.Printf("Sharpe %f\n", a.Sharpe(alpha2, r, equitySTD))
	eqc := a.NewBlackScholes(equity, &Eopt1c, &Eunder, r)
	//eqp := a.NewBlackScholes(&opt1p, &under, r)
	fmt.Printf("Call StrikePrice %f\n", eqc.StrikePrice)
	fmt.Println()
	//fmt.Printf("Put StrikePrice %f\n", eqp.StrikePrice)
	//fmt.Println()
	fmt.Printf("Call D1 %f\n", eqc.D1)
	fmt.Println()
	//fmt.Printf("Put D1 %f\n", eqp.D1)
	//fmt.Println()
	fmt.Printf("Call D2 %f\n", eqc.D2)
	fmt.Println()
	//fmt.Printf("Put D2 %f\n", eqp.D2)
	//fmt.Println()
	fmt.Printf("Call Delta %f\n", eqc.Delta)
	fmt.Println()
	//fmt.Printf("Put Delta %f\n", eqp.Delta)
	//fmt.Println()
	fmt.Printf("Call ImpliedVolatility %f\n", eqc.ImpliedVolatility)
	fmt.Println()
	//fmt.Printf("Put %f\n", eqp.ImpliedVolatility)
	//fmt.Println()
	fmt.Printf("Call Theta %f\n", eqc.Theta)
	//fmt.Println()
	//fmt.Printf("Put Theta %f\n", eqp.Theta)
	fmt.Println()
	fmt.Printf("Call Gamma %f\n", eqc.Gamma)
	fmt.Println()
	//fmt.Printf("Put Gamma %f\n", eqp.Gamma)
	//fmt.Println()
	fmt.Printf("Call Kappa %f\n", eqc.Kappa)
	fmt.Println()
	//fmt.Printf("Put Kappa %f\n", eqp.Kappa)
	//fmt.Println()
	fmt.Printf("Call Rho %f\n", eqc.Rho)
	fmt.Println()
	//fmt.Printf("Put Rho %f\n", eqp.Rho)
	//fmt.Println()
	fmt.Printf("Call TheoPrice %f\n", eqc.TheoPrice)
	fmt.Println()
	//fmt.Printf("Put TheoPrice %f\n", eqp.TheoPrice)
	//fmt.Println()
	fmt.Printf("Call UnderlyingPrice %f\n", eqc.UnderlyingPrice)
	fmt.Println()
	//fmt.Printf("Put UnderlyingPrice %f\n", eqp.UnderlyingPrice)
	//fmt.Println()
	fmt.Printf("Half Life Sigma %f\n", eqc.HalfLifeSTD)
	fmt.Println()
	fmt.Printf("Half Life Theta %f\n", eqc.HalfLifeTheta)
	fmt.Println()
	fmt.Printf("Strike %f\n", equityK)
	fmt.Println()
	fmt.Printf("Time to Expiration %f\n", equityt)
	fmt.Println()

	i1 := a.Intrinsic(true, benchmark)
	i2 := a.Intrinsic(false, equity)
	e1 := a.Extrinsic(true, benchmark)
	e2 := a.Extrinsic(false, equity)
	fmt.Printf("Benmark Intrinsic Value %f\n", i1)
	fmt.Println()
	fmt.Printf("Equity Intrinsic Value %f\n", i2)
	fmt.Println()
	fmt.Printf("Benmark Extrinsic Value %f\n", e1)
	fmt.Println()
	fmt.Printf("Equity Extrinsic Value %f\n", e2)
	fmt.Println()

	sl, in, rsq, cou, se, ise := autoStrategy.LinearRegression(eq, benmark)
	fmt.Printf("SL %f\n", sl)
	fmt.Println()
	fmt.Printf("IN %f\n", in)
	fmt.Println()
	fmt.Printf("RSQ %f\n", rsq)
	fmt.Println()
	fmt.Printf("COU %v\n", cou)
	fmt.Println()
	fmt.Printf("SE %f\n", se)
	fmt.Println()
	fmt.Printf("ISE %f\n", ise)

	cdfd1 := a.Cdf(benmark, benPrice)
	pdfd1 := a.Pdf(benmark, benPrice)
	cdfd2 := a.Cdf(eq, equityPrice)
	pdfd2 := a.Pdf(eq, equityPrice)

	fmt.Println()
	fmt.Printf("CDF eqchmark %f\n", cdfd1)
	fmt.Println()
	fmt.Printf("PDF eqchmark %f\n", pdfd1)
	fmt.Println()
	fmt.Printf("CDF Equity %f\n", cdfd2)
	fmt.Println()
	fmt.Printf("PDF Equity %f\n", pdfd2)
	fmt.Println()

	ben1 := eq[len(eq)-1]
	x1, y1 := autoStrategy.Predict(sl, ben1, 0, in)
	fmt.Printf("X %f, Y %f\n", x1, y1)
	fmt.Println()
	ben2 := benmark[len(benmark)-1]
	x2, y2 := autoStrategy.Predict(sl, 0, ben2, in)
	fmt.Printf("X %f, Y %f\n", x2, y2)
	p := plot.New()
	p.Title.Text = "SPY vs Asset"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	pts := make(plotter.XYs, len(eq))

	if len(benmark) > len(eq) {
		for index, _ := range eq {
			pts[index].X = eq[index]
			pts[index].Y = benmark[index]
		}
	} else {
		for index, _ := range benmark {
			pts[index].X = eq[index]
			pts[index].Y = benmark[index]
		}
	}

	p.Add(plotter.NewGrid())
	plotutil.AddLinePoints(p, pts)
	me, lo, hi := plotutil.MeanAndConf95(benmark)
	fmt.Println()
	fmt.Printf("Mean %f, Low %f, High %f\n", me, lo, hi)
	me2, lo2, hi2 := plotutil.MeanAndConf95(eq)
	fmt.Println()
	fmt.Printf("Mean %f, Low %f, High %f\n", me2, lo2, hi2)

	str := fmt.Sprintf("SPYvsTSLA2.png")
	if err := p.Save(4*vg.Inch, 4*vg.Inch, str); err != nil {
		panic(err)
	}

	//time.Sleep(1 * time.Minute)
	//}
}
