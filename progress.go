package main

import (
	"fmt"
	"strings"
	"time"
)

func progressBar(current, total int, width int) string {
	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s] %d/%d", bar, current, total)
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
}

// StageTimer tracks per-stage durations and can estimate remaining time for a workflow.
//
// It is intentionally simple: each call to StageComplete() records the time since the previous
// stage and resets the timer baseline.
type StageTimer struct {
	StartTime time.Time
	Stages    []time.Duration
}

// NewStageTimer creates a timer starting "now".
func NewStageTimer() *StageTimer {
	return &StageTimer{StartTime: time.Now()}
}

// StageComplete records the duration since the last stage boundary and resets the baseline.
func (t *StageTimer) StageComplete() {
	t.Stages = append(t.Stages, time.Since(t.StartTime))
	t.StartTime = time.Now()
}

// Elapsed returns total elapsed time across completed stages plus the current in-progress stage.
func (t *StageTimer) Elapsed() time.Duration {
	total := time.Duration(0)
	for _, d := range t.Stages {
		total += d
	}
	return total + time.Since(t.StartTime)
}

// EstimateRemaining returns a human-readable ETA for the remaining stages.
func (t *StageTimer) EstimateRemaining(currentStage, totalStages int) string {
	if currentStage == 0 || len(t.Stages) == 0 {
		return "estimating..."
	}
	avgPerStage := t.Elapsed() / time.Duration(currentStage)
	remaining := avgPerStage * time.Duration(totalStages-currentStage)
	return formatDuration(remaining)
}
