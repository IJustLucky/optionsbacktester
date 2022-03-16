package main

import (
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

type login struct {
	api_key    string
	secret_key string
	base_url   string
	data_url   string
	data_feed  string
}

type strat struct {
	stocks            []string
	crypto            []string
	symbols           []string
	strat_mode        []string
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
