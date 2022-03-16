package autoStrategy

type cNode struct {
	name string
	cond bool
}

type vNode struct {
	name  string
	value float64
}

type StrategyNode struct {
	conditionNodes []cNode
	valuesNodes    []vNode
}

type StrategyUnfinished struct {
	name             string
	ticks            []int
	tradeSizePercent []float32
	tradeSizeQty     []float32
	symbols          []string
	lookback         []int
	time             []string
	takeProfit       []float32
	stopLoss         []float32
	trail            []int
	slippage         []int
	strategyType     []string
}

type StrategyFinished struct {
	name             string
	ticks            int
	tradeSizePercent float32
	tradeSizeQty     float32
	symbols          string
	lookback         int
	time             string
	takeProfit       float32
	stopLoss         float32
	trail            int
	slippage         int
	strategyType     string
}

type Results struct {
	equity      float32
	netProfit   float32
	grossProfit float32
	expenses    float32
	dropDown    float32
	sharpe      float32
	alpha       float32
	beta        float32
	delta       float32
	gamma       float32
	theta       float32
	kappa       float32
	rho         float32
}

var maCross StrategyUnfinished

type strategyMaker interface {
	DeriveValues() StrategyUnfinished
	ChooseValues() StrategyFinished
}

func (s StrategyUnfinished) DeriveValues(n string) StrategyUnfinished {
	s.name = n
	s.ticks = []int{1, 5, 15, 25, 50, 75, 100, 120, 180, 300, 600, 3600}
	s.tradeSizePercent = []float32{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1.0}
	s.tradeSizeQty = []float32{500, 1000, 2500, 5000, 10000}
	s.symbols = []string{"SPY"}
	s.lookback = []int{1, 2, 3, 5, 8, 10, 15, 25, 50, 100, 200, 500}
	s.time = []string{"1m", "5m", "15m", "30m", "1h", "4h", "1d"}
	s.takeProfit = []float32{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1.0}
	s.stopLoss = []float32{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1.0}
	s.trail = []int{1, 2, 3, 4, 5, 10, 15, 25, 50}
	s.slippage = []int{1, 2, 3, 4, 5, 10, 15, 25, 50}
	s.strategyType = []string{"both", "long", "short"}
	return s
}

func (s StrategyFinished) ChooseValues(n string, tick int, tsp float32, qty float32, sym string, lb int, t string, tp float32, sl float32, tl int, sli int, st string) StrategyFinished {
	s.name = n
	s.ticks = tick
	s.tradeSizePercent = tsp
	s.tradeSizeQty = qty
	s.symbols = sym
	s.lookback = lb
	s.time = t
	s.takeProfit = tp
	s.stopLoss = sl
	s.trail = tl
	s.slippage = sli
	s.strategyType = st
	return s
}
