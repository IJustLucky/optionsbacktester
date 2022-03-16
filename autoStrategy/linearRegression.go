package autoStrategy

import (
	"alpaca_dynamic_alpha/options"
	"math"
)

type LinearRegression struct {
	intercept float64
	slope     float64
}

func (l LinearRegression) Regression(xd []float64, yd []float64) (a float64, b float64) {
	x := options.Sum(xd)
	y := options.Sum(yd)
	var xy []float64
	for _, v := range xd {
		for _, j := range yd {
			xy = append(xy, v*j)
		}
	}

	a := y*math.Pow(x, 2) - x*options.Sum(xy)/options.Length()

}
