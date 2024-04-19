// Copyright 2015 The go-ethereum Authors
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
	diskStats  DiskStats
	diskMetrics struct {
		ReadCount  Gauge
		ReadBytes  Gauge
		WriteCount Gauge
		WriteBytes Gauge
	}
	registerDiskMetricsOnce = sync.Once{}
)

// DiskStats is the per process disk io stats.
type DiskStats struct {
	ReadCount  int64 // Number of read operations executed
	ReadBytes  int64 // Total number of bytes read
	WriteCount int64 // Number of write operations executed
	WriteBytes int64 // Total number of byte written
}

// CaptureDiskStats captures new values for the Go process disk usage
// statistics exported in disk.DiskStats. This is designed to be called as a
// goroutine.
func CaptureDiskStats(d time.Duration) {
	for range time.Tick(d) {
		CaptureDiskStatsOnce()
	}
}

// CaptureDiskStatsOnce captures new values for the Go process disk usage
// statistics exported in disk.DiskStats. This is designed to be called in a
// background goroutine.
func CaptureDiskStatsOnce() {
	err := ReadDiskStats(&diskStats)
	if err != nil {
		panic(err)
	}
	diskMetrics.ReadCount.Update(diskStats.ReadCount)
	diskMetrics.ReadBytes.Update(diskStats.ReadBytes)
	diskMetrics.WriteCount.Update(diskStats.WriteCount)
	diskMetrics.WriteBytes.Update(diskStats.WriteBytes)
}

// RegisterDiskStats registers metrics for the Go process disk usage statistics
// exported in disk.DiskStats.
func RegisterDiskStats(r Registry) {
	if r == nil {
		r = DefaultRegistry
	}
	registerDiskMetricsOnce.Do(func() {
		diskMetrics.ReadCount = NewGauge()
		diskMetrics.ReadBytes = NewGauge()
		diskMetrics.WriteCount = NewGauge()
		diskMetrics.WriteBytes = NewGauge()
		r.Register("disk.DiskStats.ReadCount", diskMetrics.ReadCount)
		r.Register("disk.DiskStats.ReadBytes", diskMetrics.ReadBytes)
		r.Register("disk.DiskStats.WriteCount", diskMetrics.WriteCount)
		r.Register("disk.DiskStats.WriteBytes", diskMetrics.WriteBytes)
	})
}
