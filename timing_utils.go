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
	lastPCR := tp.lastVideoTime.Load().(time.Duration)
	lastTime := tp.lastKeyframe.Load().(time.Time)

	if !lastTime.IsZero() {
		expected := currentPCR - lastPCR
		actual := now.Sub(lastTime)
		drift := actual - expected

		if absDuration(drift) > maxDrift+warnBuffer {
			tp.stats.PCRDrifts++
			log.Printf("[PCR] Drift: %v (allowed ±%v)", drift, maxDrift)
		}
	}
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

		expected := currentPCR - prevPCR
		actual := now.Sub(prevTime)
		drift := actual - expected

		if absDuration(drift) > maxDrift+warnBuffer {
			log.Printf("[PCR] Significant drift detected: %v (allowed: ±%v)",
				drift, maxDrift)
		}
	}

	lastPCR.Store(now)
	lastPCRVal.Store(currentPCR)
}
