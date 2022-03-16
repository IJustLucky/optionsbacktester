package equities

import (
	movingaverage "github.com/RobinUS2/golang-moving-average"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata/stream"
	_ "github.com/shopspring/decimal"
)

const (
	windowSize = 20
)

type alpacaClientContainer struct {
	tradeClient   alpaca.Client
	dataClient    marketdata.Client
	streamClient  stream.StocksClient
	feed          string
	movingAverage *movingaverage.MovingAverage
	lastOrder     string
	stock         string
}

type strat struct {
	stocks            []string
	crypto            []string
	symbols           []string
	strat_mode        string
	strat_cond        bool
	timeframe         []marketdata.TimeFrame
	time_start        string
	time_end          string
	equity            int
	tick_size         []int
	index             int
	tick              int
	tradesize_qty     []int
	tradesize_percent []float64
	order_mode        []string
	lookback          []int
	last_order        string
}

func (s strat) getSymbols() (symbols []string) {
	for _, v := range s.stocks {
		s.symbols = append(s.symbols, v)
	}
	for _, v := range s.crypto {
		s.symbols = append(s.crypto, v)
	}
	return s.symbols
}

var algo alpacaClientContainer

var s strat
