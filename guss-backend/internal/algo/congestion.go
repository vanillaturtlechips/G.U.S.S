package algo

import (
	"sync"
	"time"
)

var hourWeights = [24]float64{
	0: 0.5, 1: 0.4, 2: 0.3, 3: 0.3, 4: 0.4, 5: 0.6,
	6: 0.8, 7: 1.1, 8: 1.2, 9: 1.0, 10: 0.8, 11: 0.7,
	12: 0.8, 13: 0.8, 14: 0.7, 15: 0.7, 16: 0.9, 17: 1.1,
	18: 1.4, 19: 1.5, 20: 1.4, 21: 1.2, 22: 0.9, 23: 0.7,
}

var calculatorPool = sync.Pool{
	New: func() interface{} {
		return new(float64)
	},
}

func Calculate(currentUsers, maxCapacity int) float64 {
	// 정원이 0인 경우 방어 로직
	if maxCapacity <= 0 {
		return 0.0
	}

	baseRatio := float64(currentUsers) / float64(maxCapacity)

	currentHour := time.Now().Hour()
	weight := hourWeights[currentHour]
	
	finalCongestion := baseRatio * weight

	if finalCongestion > 1.0 {
		finalCongestion = 1.0
	}
	if finalCongestion < 0.0 {
		finalCongestion = 0.0
	}

	return finalCongestion
}

func ApplyEMA(prevEMA, newVal float64) float64 {
	alpha := 0.2 // 최신 데이터에 20%의 중요도를 둠
	return (newVal * alpha) + (prevEMA * (1 - alpha))
}

type CongestionCalculator interface {
	Calculate(currentUsers, maxCapacity int) float64
}

type RealTimeCalculator struct{}

func (c *RealTimeCalculator) Calculate(current, max int) float64 {
	return Calculate(current, max) // 기존에 만든 로직 호출
}