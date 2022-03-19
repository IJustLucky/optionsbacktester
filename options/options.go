package options

import (
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/chobie/go-gaussian"
	"github.com/piquette/finance-go/equity"
	"github.com/piquette/finance-go/options"
	"log"
	"math"
	"time"
)

const (
	interest = 7.9
)

func Alpaca(str string) (x []float64, p1 float64, std float64) {
	apiKey := "PKSFCPJRTS1RQU81LHIX"
	apiSecret := "QcCwcGWATUHfiT2xDWj2cb9xNua98QjyDz5Y9tHH"
	// baseURL := "https://paper-api.Alpaca.markets"
	feed := "sip"
	equity, err := marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		Feed:      feed,
	}).GetBars(str, marketdata.GetBarsParams{
		TimeFrame: marketdata.OneMin,
		Start:     time.Date(2022, 3, 18, 9, 30, 0, 0, time.UTC),
		End:       time.Now(),
	})

	if err != nil {
		log.Fatal(err)
	}

	var eq []float64
	for _, v := range equity {
		eq = append(eq, v.Close)
	}

	var stddev float64
	for _, v := range eq {
		stddev += v
	}

	stddev = math.Sqrt(stddev / float64(len(equity)))
	return eq, eq[0], stddev
}

func GetOptionsData(str string) (K float64, q float64, t float64) {
	s, err := equity.Get(str)
	if err != nil {
		log.Fatal(err)
	}
	o := options.GetStraddle(str)
	K = o.Meta().Strikes[1]
	q = s.TrailingAnnualDividendYield
	t1 := o.Meta().ExpirationDate
	t = float64(t1 / 60 / 60 / 24 / 365)
	return K, q, t
}

type Option struct {
	StrikePrice      float64
	TimeToExpiration float64
	Type             string
}

type Underlying struct {
	Symbol     string
	Price      float64
	Volatility float64
}

type BS struct {
	StrikePrice          float64
	UnderlyingPrice      float64
	RiskFreeInterestRate float64
	Volatility           float64
	TimeToExpiration     float64
	Type                 string
	D1                   float64
	D2                   float64
	Delta                float64
	Theta                float64
	Gamma                float64
	Kappa                float64
	Rho                  float64
	ImpliedVolatility    float64
	TheoPrice            float64
	HalfLifeSTD          float64
	HalfLifeTheta        float64
	Norm                 *gaussian.Gaussian
}

func Mean(x []float64) float64 {
	var Mu float64
	for _, v := range x {
		Mu += v
	}
	return Mu / float64(len(x))
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

func Alpha(p []float64, b float64, r float64, x []float64) float64 {
	mr := Mean(p)
	mm := Mean(x)
	lele := p[len(p)-1:]
	Ru := lele[0] - mr
	lele2 := x[len(x)-1:]
	Rmu := lele2[0] - mm
	R1 := Ru / mr
	Rf := r / mr
	Rm := Rmu / mr
	bdif := Rm - Rf
	a := R1 - (Rf + b*bdif)
	return a
}

func Sharpe(a float64, r float64, v float64) float64 {
	dif := a - r
	return dif / v
}

func Beta(p []float64, x []float64) float64 {
	Mu := Mean(x)
	return Covariance(p, x) / Variance(Mu, x)
}

func Variance(Mu float64, x []float64) float64 {
	var sigmasq float64
	for _, v := range x {
		sigmasq += math.Pow(v-Mu, 2)
	}
	return sigmasq / float64(len(x))
}

func Pdf(x []float64, d float64) float64 {
	var expon float64
	Mu := Mean(x)
	Sigma := math.Sqrt(Variance(Mu, x))
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
	Sigma := math.Sqrt(Variance(Mu, x))
	dist := d - Mu
	errf := Errf(dist / (Sigma * math.Sqrt(2)))
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

func NewBlackScholes(str string, option *Option, underlying *Underlying, riskFreeInterestRate float64) *BS {
	bs := &BS{
		StrikePrice:          option.StrikePrice,
		UnderlyingPrice:      underlying.Price,
		RiskFreeInterestRate: riskFreeInterestRate,
		Volatility:           underlying.Volatility,
		TimeToExpiration:     option.TimeToExpiration / 365,
		Type:                 option.Type,
	}

	bs.Initialize(str)

	return bs
}

func R(interest float64) float64 {
	r, _ := math.Lgamma(interest)
	return r
}

func (bs *BS) Initialize(str string) {
	eq, _, std := Alpaca(str)
	bs.Norm = gaussian.NewGaussian(0, 1)
	bs.D1 = bs.calcD1(bs.UnderlyingPrice, bs.StrikePrice, bs.RiskFreeInterestRate, bs.TimeToExpiration, bs.Volatility)
	bs.D2 = bs.calcD2(bs.D1, bs.Volatility, bs.TimeToExpiration)
	bs.Delta = bs.calcDelta()
	bs.TheoPrice = bs.calcTheoreticalPrice()
	bs.ImpliedVolatility = bs.calcIv()
	bs.Theta = bs.calcTheta()
	bs.Gamma = bs.calcGamma(eq)
	bs.Kappa = bs.calcKappa()
	bs.Rho = bs.calcRho()
	bs.HalfLifeSTD, bs.HalfLifeTheta = bs.calcHalfLife(std)
}

func (bs *BS) HistoricalVolatility() {

}

func (bs *BS) StandardDeviation(days int, dataPoints []float64) float64 {
	data := dataPoints[len(dataPoints)-days:]

	var total float64

	for _, d := range data {
		total += d
	}

	mean := total / float64(days)

	var temp float64

	for _, d := range data {
		temp += math.Pow(d-mean, 2)
	}

	return math.Sqrt(temp / float64(days))
}

func (bs *BS) calcD1(underlyingPrice float64, strikePrice float64, riskFreeInterestRate float64, timeToExpiration float64, volatility float64) float64 {
	return (math.Log(underlyingPrice/strikePrice) + (riskFreeInterestRate+math.Pow(volatility, 2)/2)*timeToExpiration) / (volatility * math.Sqrt(timeToExpiration))
}

func (bs *BS) calcD2(d1 float64, volatility float64, timeToExpiration float64) float64 {
	return d1 - (volatility * math.Sqrt(timeToExpiration))
}

func (bs *BS) calcDelta() float64 {
	return bs.Norm.Cdf(bs.D1)
}

func (bs *BS) calcGamma(x []float64) float64 {
	e := math.E
	ex := math.Pow(0.5*bs.D1, 2)
	top := math.Pow(e, ex)
	Mu := Mean(x)
	bottom := bs.StrikePrice * Variance(Mu, x) * math.Sqrt(2*math.Pi*bs.TimeToExpiration)
	return top / bottom
}

func (bs *BS) calcKappa() float64 {
	n := math.Pow(math.E, math.Pow(bs.D1, 2)/-2) / math.Sqrt(2*math.Pi)
	kappa := bs.UnderlyingPrice * n * math.Sqrt(bs.TimeToExpiration)
	return kappa
}

func (bs *BS) calcRho() float64 {
	r := R(interest)
	rho := math.Pow(bs.StrikePrice*bs.TimeToExpiration*math.E, r*bs.TimeToExpiration*bs.Norm.Cdf(bs.D2))
	return rho
}

func (bs *BS) calcTheta() float64 {
	return -((bs.UnderlyingPrice*bs.Volatility*bs.Norm.Cdf(bs.D1))/(2*math.Sqrt(bs.TimeToExpiration)) - bs.RiskFreeInterestRate*bs.StrikePrice*math.Exp(-bs.RiskFreeInterestRate*(bs.TimeToExpiration))*bs.Norm.Cdf(bs.D2)) / 365
}

func (bs *BS) calcIv() float64 {
	vol := math.Sqrt(2*math.Pi/bs.TimeToExpiration) * bs.TheoPrice / bs.UnderlyingPrice

	for i := 0; i < 100; i++ {

		d1 := bs.calcD1(bs.UnderlyingPrice, bs.StrikePrice, bs.RiskFreeInterestRate, bs.TimeToExpiration, vol)
		d2 := bs.calcD2(d1, vol, bs.TimeToExpiration)
		vega := bs.UnderlyingPrice * bs.Norm.Cdf(d1) * math.Sqrt(bs.TimeToExpiration)

		cp := 1.0
		if bs.Type == "PUT" {
			cp = -1
		}

		price0 := cp*bs.UnderlyingPrice*bs.Norm.Cdf(cp*d1) - cp*bs.StrikePrice*math.Exp(bs.RiskFreeInterestRate*bs.TimeToExpiration)*bs.Norm.Cdf(cp*d2)
		vol = vol - (price0-bs.TheoPrice)/vega

		if math.Abs(price0-bs.TheoPrice) < math.Pow(10, -25) {
			break
		}
	}
	return vol
}

func (bs *BS) calcTheoreticalPrice() float64 {
	normD1 := bs.Norm.Cdf(bs.D1)
	normD2 := bs.Norm.Cdf(bs.D2)

	return bs.UnderlyingPrice*normD1 - bs.StrikePrice*math.Pow(math.E, -bs.RiskFreeInterestRate*bs.TimeToExpiration)*normD2
}

func (bs *BS) calcHalfLife(std float64) (float64, float64) {
	nl := math.Gamma(2)
	return (nl / std) * 365, (nl / bs.Theta) * 365
}

func Intrinsic(b bool, str string) float64 {
	r := R(interest)
	_, S, _ := Alpaca(str)
	K, Q, t := GetOptionsData(str)
	p1 := math.Exp(-Q*t)*S - math.Exp(-r*t)*K
	switch b {
	case true:
		return math.Max(0, +p1)
	case false:
		return math.Max(0, -p1)
	}
	return math.Abs(p1)
}

func Extrinsic(b bool, str string) float64 {
	_, S, _ := Alpaca(str)
	i1 := Intrinsic(true, str)
	i2 := Intrinsic(false, str)
	switch b {
	case true:
		return S - i1
	case false:
		return S - i2
	}
	return 0.0
}
