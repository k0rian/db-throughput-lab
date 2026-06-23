package benchmark

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	Name       string        `json:"name"`
	Duration   string        `json:"duration"`
	TotalTasks int64         `json:"total_tasks"`
	Successes  int64         `json:"successes"`
	Failures   int64         `json:"failures"`
	QPS        float64       `json:"qps"`
	AvgLatency string        `json:"avg_latency"`
}

// RunThroughputTest runs a task concurrently for a specified duration and returns benchmark results.
// This allows testing and comparing different SQL statements or cached vs non-cached queries.
func RunThroughputTest(name string, duration time.Duration, concurrency int, task func() error) Result {
	var successes int64
	var failures int64
	var totalLatency int64 // in nanoseconds

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(concurrency)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					taskStart := time.Now()
					err := task()
					latency := time.Since(taskStart).Nanoseconds()
					atomic.AddInt64(&totalLatency, latency)

					if err != nil {
						atomic.AddInt64(&failures, 1)
					} else {
						atomic.AddInt64(&successes, 1)
					}
				}
			}
		}()
	}

	wg.Wait()
	actualDuration := time.Since(start)

	totalTasks := successes + failures
	qps := float64(totalTasks) / actualDuration.Seconds()
	var avgLatency time.Duration
	if totalTasks > 0 {
		avgLatency = time.Duration(totalLatency / totalTasks)
	}

	return Result{
		Name:       name,
		Duration:   actualDuration.String(),
		TotalTasks: totalTasks,
		Successes:  successes,
		Failures:   failures,
		QPS:        qps,
		AvgLatency: avgLatency.String(),
	}
}
