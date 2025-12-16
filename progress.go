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

type StageTimer struct {
	StartTime time.Time
	Stages    []time.Duration
}

func NewStageTimer() *StageTimer {
	return &StageTimer{StartTime: time.Now()}
}

func (t *StageTimer) StageComplete() {
	t.Stages = append(t.Stages, time.Since(t.StartTime))
	t.StartTime = time.Now()
}

func (t *StageTimer) Elapsed() time.Duration {
	total := time.Duration(0)
	for _, d := range t.Stages {
		total += d
	}
	return total + time.Since(t.StartTime)
}

func (t *StageTimer) EstimateRemaining(currentStage, totalStages int) string {
	if currentStage == 0 || len(t.Stages) == 0 {
		return "estimating..."
	}
	avgPerStage := t.Elapsed() / time.Duration(currentStage)
	remaining := avgPerStage * time.Duration(totalStages-currentStage)
	return formatDuration(remaining)
}
