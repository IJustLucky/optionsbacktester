package options

import (
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/piquette/finance-go/equity"
	"github.com/piquette/finance-go/options"
	"log"
	"math"
	"time"
)

const (
	interest = 7.9
	ticker   = "SPY"
	bench    = "UVXY"
)

func Alpaca() (x []float64, y []float64, p1 float64, p2 float64, std1 float64, std2 float64) {
	apiKey := "PKRBH1534VMCJOEETZFW"
	apiSecret := "1UFIFRbyo50wbeb96OuBO7JfhgeaoPVXBw7wFBOw"
	// baseURL := "https://paper-api.Alpaca.markets"
	feed := "sip"
	spy, err := marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		Feed:      feed,
	}).GetBars(ticker, marketdata.GetBarsParams{
		TimeFrame: marketdata.OneMin,
		Start:     time.Date(2022, 3, 1, 13, 30, 0, 0, time.UTC),
		End:       time.Date(2022, 3, 13, 13, 30, 1, 0, time.UTC),
	})

	if err != nil {
		log.Fatal(err)
	}

	var eq []float64
	for _, v := range spy {
		eq = append(eq, v.Close)
	}

	var stddev float64
	for _, v := range eq {
		stddev += v
	}

	stddev = math.Sqrt(stddev / float64(len(spy)))

	tic, err := marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		Feed:      feed,
	}).GetBars(bench, marketdata.GetBarsParams{
		TimeFrame: marketdata.OneMin,
		Start:     time.Date(2022, 3, 1, 13, 30, 0, 0, time.UTC),
		End:       time.Date(2022, 3, 13, 13, 30, 1, 0, time.UTC),
	})
	if err != nil {
		log.Fatal(err)
	}
	var earr []float64
	for _, j := range tic {
		earr = append(earr, j.Close)
	}

	var stddev2 float64
	for _, v := range earr {
		stddev2 += v
	}
	stddev2 = math.Sqrt(stddev2 / float64(len(earr)))
	return eq, earr, eq[0], earr[0], stddev, stddev2
}

func GetOptionsData() (K float64, q float64, t float64) {
	s, err := equity.Get(ticker)
	if err != nil {
		log.Fatal(err)
	}
	o := options.GetStraddle(ticker)
	K = o.Meta().Strikes[1]
	q = s.TrailingAnnualDividendYield
	days := float64(1)
	t = days
	return K, q, 90
}

func Length(p []float64) float64 {
	return float64(len(p))
}

func Sum(p []float64) float64 {
	s := 0.0
	for _, v := range p {
		s = s + v
	}
	return s
}

func Mean(p []float64) float64 {
	return Sum(p) / Length(p)
}

func Covariance(e1 []float64, e2 []float64) float64 {
	var c1c []float64
	var c2c []float64
	var cmul []float64
	var csum float64
	for _, v := range e1 {
		c1c = append(c1c, v-Mean(e1))
	}
	for _, v := range e2 {
		c2c = append(c2c, v-Mean(e2))
	}

	if Length(e1) > Length(e2) {
		for i := range e2 {
			cmul = append(cmul, c1c[i]*c2c[i])
		}
	} else if Length(e1) < Length(e2) {
		for i := range e1 {
			cmul = append(cmul, c1c[i]*c2c[i])
		}
	} else {
		return 0
	}

	for _, v := range cmul {
		csum += v
	}

	l := (Length(e1)+Length(e2))/2 - 1
	cov := csum / l
	return cov
}

func Variance(p []float64) float64 {
	l := Length(p) - 1
	m := Mean(p)
	var sigmasq float64
	for i := range p {
		sigmasq += math.Pow(p[i]-m, 2)
	}
	return sigmasq / (l - 1)
}

func Alpha(p []float64, b float64, r float64, x []float64) float64 {
	mr := Mean(p)
	mm := Mean(x)
	lele := p[len(p)-1:]
	Ru := lele[0] - mr
	lele2 := x[len(x)-1:]
	Rmu := lele2[0] - mm
	R := Ru / mr
	Rf := r / mr
	Rm := Rmu / mr
	bdif := Rm - Rf
	a := R - (Rf + b*bdif)
	return a
}

func Sharpe(a float64, r float64, v float64) float64 {
	dif := a - r
	return dif / v
}

func Beta(p []float64, x []float64) float64 {
	return Covariance(p, x) / Variance(x)
}

func D(p []float64) (d1 float64, d2 float64) {
	_, _, S, _, _, _ := Alpaca()
	K, _, _ := GetOptionsData()
	R := R(interest)
	_, Q, _ := GetOptionsData()
	_, _, T := GetOptionsData()
	V := Variance(p)
	V = math.Sqrt(V)
	div := S / K
	g, _ := math.Lgamma(div)
	b := R - Q
	w := math.Pow(V, 2)
	w = w / 2
	c := b + w
	j := T * c
	o := S / K
	d1 = (g*o + j) / (V * math.Sqrt(T))
	d2 = d1 - (V * math.Sqrt(T))
	return d1, d2
}

func DeltaC(x []float64) float64 {
	_, Q, _ := GetOptionsData()
	_, _, T := GetOptionsData()
	e := math.E
	d1, _ := D(x)
	l := math.Pow(e, -T*Q)
	c := Cdf(x, d1)
	return l * c
}

func DeltaP(x []float64) float64 {
	_, Q, _ := GetOptionsData()
	_, _, T := GetOptionsData()
	e := math.E
	d1, _ := D(x)
	c := Cdf(x, d1)
	l := math.Pow(e, -T*Q)
	return -(l * c)
}

func Gamma() float64 {
	R := R(interest)
	_, _, S, _, std, _ := Alpaca()
	K, Q, t := GetOptionsData()
	drq := math.Exp(-Q * t)
	drd := S * std * math.Sqrt(t)
	d1pdf := D1pdff(S, K, t, std, Q, R)
	gamma1 := (drq / drd) * d1pdf
	return gamma1
}

func Kappa(b bool) float64 {
	R := R(interest)
	_, _, S, _, std, _ := Alpaca()
	K, Q, t := GetOptionsData()
	d1 := D1pdff(S, K, t, std, Q, R)
	t = math.Sqrt(t)
	if b == true || b == false {
		return S * math.Exp(-Q*t-d1*d1/2) * math.Sqrt(t) * math.Sqrt(2*math.Pi) / 100
	}

	return 2 * math.Exp(-Q*t-d1*d1/2) * math.Sqrt(t) * math.Sqrt(2*math.Pi) / 100
}

func RhoC(x []float64) float64 {
	K, _, _ := GetOptionsData()
	R := R(interest)
	R = R / 100
	_, d2 := D(x)
	_, _, T := GetOptionsData()
	e := math.E
	l := math.Pow(e, -1.0*R*T)
	N := Cdf(x, d2)
	return K * T * l * N / 100
}

func RhoP(x []float64) float64 {
	K, _, _ := GetOptionsData()
	R := R(interest)
	_, _, T := GetOptionsData()
	_, d2 := D(x)
	e := math.E
	R = R / 100
	l := math.Pow(e, -R*T)
	N := Cdf(x, d2)
	return K * T * l * -(N) / 100
}

var sqtwopi = math.Sqrt(2 * math.Pi)
var IVPrecision = 0.00001

func D1f(underlying float64, strike float64, timeToExpiration float64, volatility float64, riskFreeInterest float64, dividend float64, volatilityWithExpiration float64) float64 {
	d1 := math.Log(underlying/strike) + (timeToExpiration * (riskFreeInterest - dividend + ((volatility * volatility) * 0.5)))
	d1 = d1 / volatilityWithExpiration
	return d1
}

func D2f(d1 float64, volatilityWithExpiration float64) float64 {
	d2 := d1 - volatilityWithExpiration
	return d2
}
func D1pdff(underlying float64, strike float64, timeToExpiration float64, volatility float64, riskFreeInterest float64, dividend float64) float64 {
	vt := volatility * (math.Sqrt(timeToExpiration))
	d1 := D1f(underlying, strike, timeToExpiration, volatility, riskFreeInterest, dividend, vt)
	d1pdf := math.Exp(-(d1 * d1) * 0.5)
	d1pdf = d1pdf / sqtwopi
	return d1pdf
}

func Theta(x []float64, callType bool, underlying float64, strike float64, timeToExpiration float64, volatility float64, riskFreeInterest float64, dividend float64) float64 {

	var sign float64
	if !callType {
		sign = -1
	} else {
		sign = 1
	}

	sqt := math.Sqrt(timeToExpiration)
	drq := math.Exp(-dividend * timeToExpiration)
	dr := math.Exp(-riskFreeInterest * timeToExpiration)
	d1pdf := D1pdff(underlying, strike, timeToExpiration, volatility, riskFreeInterest, dividend)
	twosqt := 2 * sqt
	p1 := -1 * ((underlying * volatility * drq) / twosqt) * d1pdf

	vt := volatility * (sqt)
	d1 := D1f(underlying, strike, timeToExpiration, volatility, riskFreeInterest, dividend, vt)
	d2 := D2f(d1, vt)
	var nd1, nd2 float64

	d1 = sign * d1
	d2 = sign * d2
	nd1 = Cdf(x, d1)
	nd2 = Cdf(x, d2)

	p2 := -sign * riskFreeInterest * strike * dr * nd2
	p3 := sign * dividend * underlying * drq * nd1
	theta := (p1 + p2 + p3) / 365
	return theta
}

func R(i float64) float64 {
	g, _ := math.Lgamma(1.0 + i)
	return g
}

var sqrt2 = math.Pow(2, 0.5)
var toomanydev float64 = 8

func Pdf(x []float64, d float64) float64 {
	var expon float64
	Mu := Mean(x)
	Sigma := math.Sqrt(Variance(x))
	stdsqpi := Sigma * math.Pow(2*math.Pi, 0.5)
	if Mu == 0 {
		expon = -(d * d) / Sigma
	} else {
		expon = -(math.Pow(d-Mu, 2)) / Sigma
	}
	probDist := math.Exp(expon) / stdsqpi
	return probDist
}

func Cdf(x []float64, d float64) float64 {
	Mu := Mean(x)
	Sigma := math.Sqrt(Variance(x))
	dist := d - Mu
	if math.Abs(dist) > toomanydev*Sigma {
		if d < Mu {
			return 0.0
		} else {
			return 1.0
		}
	}
	errf := Errf(dist / (Sigma * sqrt2))
	cdf := 0.5 * (1.0 + errf)
	return cdf
}

func Errf(z float64) float64 {
	var t float64
	t = 1.0 / (1.0 + 0.5*math.Abs(z))
	ans := 1 - t*math.Exp(-z*z-1.26551223+
		t*(1.00002368+
			t*(0.37409196+
				t*(0.09678418+
					t*(-0.18628806+
						t*(0.27886807+
							t*(-1.13520398+
								t*(1.48851587+
									t*(-0.82215223+
										t*(0.17087277))))))))))
	if z >= 0 {
		return ans
	}
	return -ans
}

func ImpliedVol(x []float64) float64 {
	_, _, S, _, _, _ := Alpaca()
	K, _, t := GetOptionsData()
	R := R(interest)

	d12, _ := math.Lgamma(S / K)
	s := math.Sqrt(Variance(x))
	ex := math.Pow(s, 2)
	d1 := t * (d12 + (R + ex/2))
	d2 := d1 - (s * math.Sqrt(t))
	cdf1 := Cdf(x, d1)
	cdf2 := Cdf(x, d2)
	pow := math.Pow(math.E, -R*t)
	call := (S * cdf1) - (K * cdf2 * pow)
	return call
}

func Intrinsic(b bool) float64 {
	R := R(interest)
	_, _, S, _, _, _ := Alpaca()
	K, Q, t := GetOptionsData()
	p := math.Exp(-Q*t)*S - math.Exp(-R*t)*K

	switch b {
	case true:
		return math.Max(0, +p)
	case false:
		return math.Max(0, -p)
	}
	return math.Abs(p)
}

func Extrinsic(b bool) float64 {
	_, _, S, _, _, _ := Alpaca()
	switch b {
	case true:
		return S - Intrinsic(true)
	case false:
		return S - Intrinsic(false)
	}
	return 0.0
}
