package optimization

import (
	"fmt"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/backtesting"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/metrics"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/money"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/strategy"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
	"github.com/shopspring/decimal"
	"github.com/tomcraven/goga"
	"math/rand"
)

type StrategyOptimizer struct {
}

type gogaOptimizer struct {
	Strategy         strategy.Strategy
	iterationCount   uint
	iterationsLimit  uint
	providersFunc    func() []tick.Provider
	fitnessFunc      func(metrics metrics.TradeMetrics) int
	fitnessThreshold int
	initialBalance   int64
}

func (b *gogaOptimizer) Go() goga.Bitset {
	bitset := goga.Bitset{}
	bitset.Create(len(b.Strategy.GetParameters()))

	for i, parameter := range b.Strategy.GetParameters() {
		val := rand.Float32()*float32(parameter.Max-parameter.Min) + float32(parameter.Min)

		bitset.Set(i, int(val))
	}

	return bitset
}

func (s *gogaOptimizer) OnBeginSimulation() {
}

func (s *gogaOptimizer) Simulate(genome goga.Genome) {
	s.iterationCount++
	fmt.Printf("Simulation #%d...\n", s.iterationCount)

	bits := genome.GetBits()

	for i, parameter := range s.Strategy.GetParameters() {
		bit := bits.Get(i)
		p := strategy.NewParameter(parameter.Name, parameter.Min, parameter.Max)

		err := p.SetValue(bit)

		if err != nil {
			panic(err)
		}
	}

	balanceManager := money.BalanceManager{}
	balanceManager.SetInitialBalance(decimal.NewFromInt(s.initialBalance))
	balanceManager.Reset()

	tester := backtesting.StrategyTester{}
	tester.Init(&balanceManager, s.providersFunc())

	err := tester.Run(&s.Strategy)
	if err != nil {
		panic(err)
	}

	tradeHistories := tester.GetTradeHistories()
	tradeMetrics := metrics.CalculateMetrics(tradeHistories, balanceManager.GetInitialBalance(), balanceManager.GetTotalBalance())

	fitness := s.fitnessFunc(tradeMetrics)

	if fitness > 0 {
		genome.SetFitness(fitness)

	} else {
		genome.SetFitness(0)
	}
}

func (s *gogaOptimizer) OnEndSimulation() {
}

func (s *gogaOptimizer) ExitFunc(genome goga.Genome) bool {
	if s.iterationsLimit <= s.iterationCount {
		return true
	}

	return genome.GetFitness() >= s.fitnessThreshold
}

func (e *gogaOptimizer) OnElite(genome goga.Genome) {
	e.iterationCount++

	fmt.Printf("Iteration #%d, fitness = %d\n", e.iterationCount, genome.GetFitness())
}

func (optimizer *StrategyOptimizer) Run(strategy strategy.Strategy, providersFunc func() []tick.Provider, fitnessFunc func(metrics metrics.TradeMetrics) int, fitnessThreshold int, iterationsLimit uint, initialBalance int64, parallelSimulations int) []strategy.Parameter {
	gogaOpt := gogaOptimizer{
		Strategy:         strategy,
		providersFunc:    providersFunc,
		fitnessFunc:      fitnessFunc,
		iterationsLimit:  iterationsLimit,
		fitnessThreshold: fitnessThreshold,
		initialBalance:   initialBalance,
	}

	genAlgo := goga.NewGeneticAlgorithm()
	genAlgo.BitsetCreate = &gogaOpt
	genAlgo.Selector = goga.NewSelector(
		[]goga.SelectorFunctionProbability{
			{
				P: 1,
				F: goga.Roulette,
			},
		},
	)
	genAlgo.Mater = goga.NewMater([]goga.MaterFunctionProbability{
		{
			P:        1,
			F:        goga.UniformCrossover,
			UseElite: false,
		},
	})
	genAlgo.Simulator = &gogaOpt
	genAlgo.EliteConsumer = &gogaOpt

	ps := 1
	for _, parameter := range strategy.GetParameters() {
		ps *= parameter.Max - parameter.Min
	}

	fmt.Printf("population size = %d\n", ps)

	genAlgo.Init(ps, parallelSimulations)
	genAlgo.Simulate()
	population := genAlgo.GetPopulation()

	for _, genome := range population {
		fmt.Println(genome.GetFitness())

		bits := genome.GetBits()

		for i, parameter := range strategy.GetParameters() {
			bit := bits.Get(i)
			err := parameter.SetValue(bit)

			if err != nil {
				panic(err)
			}
		}
	}

	return strategy.GetParameters()
}
