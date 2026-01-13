package algo

import (
	"sync"
)

// hourWeights: 시간대별 가중치 (로직에서 제외했으나 구조 유지를 위해 남겨둠)
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

// Calculate: 현재 인원과 최대 정원을 받아 0.0 ~ 1.0 사이의 혼잡도를 반환
func Calculate(currentUsers, maxCapacity int) float64 {
	// 정원이 0인 경우 0.0 반환 (0으로 나누기 방지)
	if maxCapacity <= 0 {
		return 0.0
	}

	// 1. 순수 비율 계산 (예: 50 / 100 = 0.5)
	baseRatio := float64(currentUsers) / float64(maxCapacity)

	// 2. [수정] 70% 결과를 만들었던 가중치 적용 로직을 제거했습니다.
	// 기존: finalCongestion := baseRatio * hourWeights[time.Now().Hour()]
	finalCongestion := baseRatio

	// 3. 결과값 범위 제한 (0.0 ~ 1.0)
	if finalCongestion > 1.0 {
		finalCongestion = 1.0
	}
	if finalCongestion < 0.0 {
		finalCongestion = 0.0
	}

	return finalCongestion
}

// ApplyEMA: 지수 이동 평균 적용 (필요 시 사용)
func ApplyEMA(prevEMA, newVal float64) float64 {
	alpha := 0.2 // 최신 데이터에 20%의 중요도를 둠
	return (newVal * alpha) + (prevEMA * (1 - alpha))
}

type CongestionCalculator interface {
	Calculate(currentUsers, maxCapacity int) float64
}

type RealTimeCalculator struct{}

// Calculate: 인터페이스 구현체 메서드
func (c *RealTimeCalculator) Calculate(current, max int) float64 {
	return Calculate(current, max)
}
