// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package metrics

import (
	"time"
	"sync"
)

var (
	cpuStats CPUStats
	cpuMetrics struct {
		GlobalTime GaugeFloat64
		GlobalWait GaugeFloat64
		LocalTime GaugeFloat64
	}
	registerCPUMetricsOnce = sync.Once{}
)

// CPUStats is the system and process CPU stats.
// All values are in seconds.
type CPUStats struct {
	GlobalTime float64 // Time spent by the CPU working on all processes.
	GlobalWait float64 // Time spent by waiting on disk for all processes.
	LocalTime  float64 // Time spent by the CPU working on this process.
}

// CaptureCPUStats captures new values for the Go process CPU usage
// statistics exported in cpu.CPUStats. This is designed to be called as a
// goroutine.
func CaptureCPUStats(d time.Duration) {
	for range time.Tick(d) {
		CaptureCPUStatsOnce()
	}
}

// CaptureCPUStatsOnce captures new values for the Go process CPU usage
// statistics exported in cpu.CPUStats. This is designed to be called in a
// background goroutine.
func CaptureCPUStatsOnce() {
	err := ReadCPUStats(&cpuStats)
	if err != nil {
		panic(err)
	}
	cpuMetrics.GlobalTime.Update(cpuStats.GlobalTime)
	cpuMetrics.GlobalWait.Update(cpuStats.GlobalWait)
	cpuMetrics.LocalTime.Update(cpuStats.LocalTime)
}

// RegisterCPUStats registers metrics for the Go process CPU usage statistics
// exported in cpu.CPUStats.
func RegisterCPUStats(r Registry) {
	if r == nil {
		r = DefaultRegistry
	}
	registerCPUMetricsOnce.Do(func() {
		cpuMetrics.GlobalTime = NewGaugeFloat64()
		cpuMetrics.GlobalWait = NewGaugeFloat64()
		cpuMetrics.LocalTime = NewGaugeFloat64()
		r.Register("cpu.CPUStats.GlobalTime", cpuMetrics.GlobalTime)
		r.Register("cpu.CPUStats.GlobalWait", cpuMetrics.GlobalWait)
		r.Register("cpu.CPUStats.LocalTime", cpuMetrics.LocalTime)
	})
}
