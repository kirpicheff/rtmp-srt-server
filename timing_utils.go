package main

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/datarhei/joy4/av"
)

const (
	maxPTS              = 1<<33 - 1
	defaultCTS          = 2 * time.Millisecond
	maxDrift            = 10 * time.Millisecond
	ptsClockRate        = 90000
	maxPTSDuration      = time.Duration(maxPTS) * time.Second / time.Duration(ptsClockRate)
	warnBuffer          = 80 * time.Millisecond
	minKeyframeInterval = 2 * time.Second
)

type TimingProcessor struct {
	baseTime      atomic.Value // time.Duration
	lastVideoTime atomic.Value // time.Duration
	lastAudioTime atomic.Value // time.Duration
	lastKeyframe  atomic.Value // time.Time
	lastPCRTime   atomic.Value // time.Time - время последнего PCR
	lastPCRValue  atomic.Value // time.Duration - значение последнего PCR
	stats         TimingStats
}

type TimingStats struct {
	TotalPackets     uint64
	CorrectedPackets uint64
	PCRDrifts        uint64
}

func NewTimingProcessor() *TimingProcessor {
	tp := &TimingProcessor{}
	tp.baseTime.Store(time.Duration(0))
	tp.lastVideoTime.Store(time.Duration(0))
	tp.lastAudioTime.Store(time.Duration(0))
	tp.lastKeyframe.Store(time.Time{})
	tp.lastPCRTime.Store(time.Time{})
	tp.lastPCRValue.Store(time.Duration(0))
	return tp
}

func (tp *TimingProcessor) Process(pkt *av.Packet) {
	tp.stats.TotalPackets++

	// Этап 1: Базовая валидация
	tp.validateBasicTiming(pkt)

	// Этап 2: Установка базового времени
	tp.setBaseTime(pkt)

	// Этап 3: Нормализация временных меток
	tp.normalizeTimestamps(pkt)

	// Этап 4: Обеспечение монотонности
	tp.enforceMonotonicity(pkt)

	// Этап 5: Финальная валидация
	tp.validateFinalTiming(pkt)

	// Этап 6: Мониторинг PCR дрейфа
	if pkt.Idx == 0 {
		tp.monitorPCR(pkt.Time)
	}
}

func (tp *TimingProcessor) validateBasicTiming(pkt *av.Packet) {
	if pkt.Time < 0 {
		pkt.Time = 0
		tp.stats.CorrectedPackets++
	}

	if pkt.CompositionTime < 0 {
		pkt.CompositionTime = 0
		tp.stats.CorrectedPackets++
	} else if pkt.IsKeyFrame && pkt.CompositionTime == 0 {
		pkt.CompositionTime = defaultCTS
	}
}

func (tp *TimingProcessor) setBaseTime(pkt *av.Packet) {
	if pkt.IsKeyFrame && tp.baseTime.Load().(time.Duration) == 0 {
		lastKF := tp.lastKeyframe.Load().(time.Time)
		if time.Since(lastKF) > minKeyframeInterval || lastKF.IsZero() {
			tp.baseTime.Store(pkt.Time)
			tp.lastKeyframe.Store(time.Now())
			log.Printf("[Timing] Base time set: %v", pkt.Time)
		}
	}
}

func (tp *TimingProcessor) normalizeTimestamps(pkt *av.Packet) {
	base := tp.baseTime.Load().(time.Duration)
	if base > 0 {
		normalized := pkt.Time - base
		if normalized < 0 {
			normalized = 0
		}
		pkt.Time = normalized
	}
}

func (tp *TimingProcessor) enforceMonotonicity(pkt *av.Packet) {
	var lastTime *atomic.Value
	if pkt.Idx == 0 {
		lastTime = &tp.lastVideoTime
	} else {
		lastTime = &tp.lastAudioTime
	}

	current := lastTime.Load().(time.Duration)
	if pkt.Time < current {
		pkt.Time = current + time.Millisecond
		tp.stats.CorrectedPackets++
	}
	lastTime.Store(pkt.Time)
}

func (tp *TimingProcessor) validateFinalTiming(pkt *av.Packet) {
	if pkt.Time > maxPTSDuration {
		pkt.Time %= maxPTSDuration
		log.Printf("[Timing] PTS wrap: %v", pkt.Time)
	}

	if pkt.CompositionTime > 0 && pkt.Time < pkt.CompositionTime {
		pkt.CompositionTime = 0
		tp.stats.CorrectedPackets++
	}
}

func (tp *TimingProcessor) monitorPCR(currentPCR time.Duration) {
	now := time.Now()
	lastPCRTime := tp.lastPCRTime.Load().(time.Time)
	lastPCRValue := tp.lastPCRValue.Load().(time.Duration)

	// Проверка разумности значения PCR
	if currentPCR < 0 {
		log.Printf("[PCR] Warning: negative PCR value: %v", currentPCR)
		return
	}

	if !lastPCRTime.IsZero() {
		// Проверка разумности интервала
		timeSinceLast := now.Sub(lastPCRTime)
		if timeSinceLast > 10*time.Second {
			log.Printf("[PCR] Warning: large time gap since last PCR: %v", timeSinceLast)
			// Сбрасываем состояние при слишком большом разрыве
			tp.resetPCRState(now, currentPCR)
			return
		}

		// Защита от переполнения - проверка монотонности PCR
		if currentPCR < lastPCRValue {
			log.Printf("[PCR] Non-monotonic: %v < %v", currentPCR, lastPCRValue)
			tp.resetPCRState(now, currentPCR)
			return
		}

		// PCR уже измеряется в единицах времени (time.Duration)
		// Ожидаемое время между PCR значениями
		expectedTimeInterval := currentPCR - lastPCRValue
		actualTimeInterval := now.Sub(lastPCRTime)
		drift := actualTimeInterval - expectedTimeInterval

		// Градуированная обработка дрейфа
		switch {
		case absDuration(drift) > 1*time.Second:
			// Критично - сбрасываем состояние
			log.Printf("[PCR] CRITICAL drift: %v (expected: %v, actual: %v)",
				drift, expectedTimeInterval, actualTimeInterval)
			tp.resetPCRState(now, currentPCR)
			return
		case absDuration(drift) > 100*time.Millisecond:
			// Предупреждение
			log.Printf("[PCR] WARNING drift: %v (expected: %v, actual: %v)",
				drift, expectedTimeInterval, actualTimeInterval)
			tp.stats.PCRDrifts++
		case absDuration(drift) > maxDrift+warnBuffer:
			// Обычное предупреждение
			log.Printf("[PCR] Drift: %v (expected: %v, actual: %v)",
				drift, expectedTimeInterval, actualTimeInterval)
			tp.stats.PCRDrifts++
		}
	} else {
		// Первое измерение - просто логируем
		log.Printf("[PCR] First PCR measurement: %v", currentPCR)
	}

	// Обновляем значения для следующего измерения
	tp.lastPCRTime.Store(now)
	tp.lastPCRValue.Store(currentPCR)
}

func (tp *TimingProcessor) resetPCRState(now time.Time, currentPCR time.Duration) {
	tp.lastPCRTime.Store(now)
	tp.lastPCRValue.Store(currentPCR)
	log.Printf("[PCR] State reset to: %v", currentPCR)
}

func (tp *TimingProcessor) GetStats() TimingStats {
	return tp.stats
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// Legacy functions for backward compatibility
func validateTiming(pkt *av.Packet) {
	tp := NewTimingProcessor()
	tp.validateBasicTiming(pkt)
	tp.validateFinalTiming(pkt)
}

func normalizeTimestamps(pkt *av.Packet, baseTime *time.Duration,
	baseTimeSet *bool, lastVideoTime, lastAudioTime *time.Duration) {

	tp := NewTimingProcessor()
	if *baseTimeSet {
		tp.baseTime.Store(*baseTime)
	}
	tp.lastVideoTime.Store(*lastVideoTime)
	tp.lastAudioTime.Store(*lastAudioTime)

	tp.Process(pkt)

	*baseTime = tp.baseTime.Load().(time.Duration)
	*baseTimeSet = tp.baseTime.Load().(time.Duration) > 0
	*lastVideoTime = tp.lastVideoTime.Load().(time.Duration)
	*lastAudioTime = tp.lastAudioTime.Load().(time.Duration)
}

func processPacketTiming(pkt *av.Packet, baseTime *time.Duration,
	baseTimeSet *bool, lastVideoTime, lastAudioTime *time.Duration) {

	tp := NewTimingProcessor()
	if *baseTimeSet {
		tp.baseTime.Store(*baseTime)
	}
	tp.lastVideoTime.Store(*lastVideoTime)
	tp.lastAudioTime.Store(*lastAudioTime)

	tp.Process(pkt)

	*baseTime = tp.baseTime.Load().(time.Duration)
	*baseTimeSet = tp.baseTime.Load().(time.Duration) > 0
	*lastVideoTime = tp.lastVideoTime.Load().(time.Duration)
	*lastAudioTime = tp.lastAudioTime.Load().(time.Duration)
}

// Global PCR monitoring (legacy)
var (
	lastPCR    atomic.Value // time.Time
	lastPCRVal atomic.Value // time.Duration
)

func checkPCRDrift(currentPCR time.Duration) {
	now := time.Now()

	if lastPCRVal.Load() != nil {
		prevPCR := lastPCRVal.Load().(time.Duration)
		prevTime := lastPCR.Load().(time.Time)

		// Защита от переполнения - проверка монотонности PCR
		if currentPCR < prevPCR {
			log.Printf("[PCR] Global non-monotonic: %v < %v", currentPCR, prevPCR)
			lastPCR.Store(now)
			lastPCRVal.Store(currentPCR)
			return
		}

		// PCR уже измеряется в единицах времени (time.Duration)
		// Ожидаемое время между PCR значениями
		expectedTimeInterval := currentPCR - prevPCR
		actualTimeInterval := now.Sub(prevTime)
		drift := actualTimeInterval - expectedTimeInterval

		// Градуированная обработка дрейфа
		switch {
		case absDuration(drift) > 1*time.Second:
			// Критично - сбрасываем состояние
			log.Printf("[PCR] Global CRITICAL drift: %v (expected: %v, actual: %v)",
				drift, expectedTimeInterval, actualTimeInterval)
			lastPCR.Store(now)
			lastPCRVal.Store(currentPCR)
			return
		case absDuration(drift) > 100*time.Millisecond:
			// Предупреждение
			log.Printf("[PCR] Global WARNING drift: %v (expected: %v, actual: %v)",
				drift, expectedTimeInterval, actualTimeInterval)
		case absDuration(drift) > maxDrift+warnBuffer:
			// Обычное предупреждение
			log.Printf("[PCR] Global drift: %v (expected: %v, actual: %v)",
				drift, expectedTimeInterval, actualTimeInterval)
		}
	}

	lastPCR.Store(now)
	lastPCRVal.Store(currentPCR)
}
