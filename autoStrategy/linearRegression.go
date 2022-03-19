package autoStrategy

import "math"

type Regression struct {
	n, sx, sy, sxx, sxy, syy float64
}

func (r *Regression) Count() int {
	return int(r.n)
}

func (r *Regression) Size() int {
	return int(r.n)
}

func (r *Regression) Update(x, y float64) {
	r.n++
	r.sx += x
	r.sy += y
	r.sxx += x * x
	r.sxy += x * y
	r.syy += y * y
}

func (r *Regression) UpdateArray(xData, yData []float64) {
	if len(xData) != len(yData) {
		if len(xData) > len(yData) {
			for i := 0; i < len(yData); i++ {
				r.Update(xData[i], yData[i])
			}
		} else {
			for i := 0; i < len(xData); i++ {
				r.Update(xData[i], yData[i])
			}
		}
	} else {
		for i := 0; i < len(yData); i++ {
			r.Update(xData[i], yData[i])
		}
	}
}

func (r *Regression) Slope() float64 {
	ss_xy := r.n*r.sxy - r.sx*r.sy
	ss_xx := r.n*r.sxx - r.sx*r.sx
	return ss_xy / ss_xx
}

func (r *Regression) Intercept() float64 {
	return (r.sy - r.Slope()*r.sx) / r.n
}

func (r *Regression) RSquared() float64 {
	ss_xy := r.n*r.sxy - r.sx*r.sy
	ss_xx := r.n*r.sxx - r.sx*r.sx
	ss_yy := r.n*r.syy - r.sy*r.sy
	return ss_xy * ss_xy / ss_xx / ss_yy
}

func (r *Regression) InterceptStandardError() float64 {
	if r.n <= 2 {
		return math.NaN()
	}
	ss_xy := r.n*r.sxy - r.sx*r.sy
	ss_xx := r.n*r.sxx - r.sx*r.sx
	ss_yy := r.n*r.syy - r.sy*r.sy
	s := math.Sqrt((ss_yy - ss_xy*ss_xy/ss_xx) / (r.n - 2.0))
	mean_x := r.sx / r.n
	return s * math.Sqrt(1.0/r.n+mean_x*mean_x/ss_xx)
}

func LinearRegression(xData, yData []float64) (slope, intercept, rsquared float64,
	count int, slopeStdErr, interceptStdErr float64) {
	var r Regression
	r.UpdateArray(xData, yData)
	ss_xy := r.n*r.sxy - r.sx*r.sy
	ss_xx := r.n*r.sxx - r.sx*r.sx
	ss_yy := r.n*r.syy - r.sy*r.sy
	slope = ss_xy / ss_xx
	intercept = (r.sy - r.Slope()*r.sx) / r.n
	rsquared = ss_xy * ss_xy / ss_xx / ss_yy
	if r.n <= 2 {
		slopeStdErr = math.NaN()
		interceptStdErr = math.NaN()
	} else {
		s := math.Sqrt((ss_yy - ss_xy*ss_xy/ss_xx) / (r.n - 2.0))
		slopeStdErr = s / math.Sqrt(ss_xx)
		mean_x := r.sx / r.n
		interceptStdErr = s * math.Sqrt(1.0/r.n+mean_x*mean_x/ss_xx)
	}
	count = int(r.n)
	return
}

func Predict(m, x1, y1, b float64) (x2, y2 float64) {
	if y1 == 0 {
		return x1, b + x1*m
	} else if x1 == 0 {
		return (y1 - b) / m, y1
	} else {
		return x1, y1
	}
}
