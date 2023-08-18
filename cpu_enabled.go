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

//go:build !ios && !js
// +build !ios,!js

package metrics

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
)

// ReadCPUStats retrieves the current CPU stats.
func ReadCPUStats(stats *CPUStats) error {
	// passing false to request all cpu times.
	timeStats, err := cpu.Times(false)
	if err != nil {
		return err
	}
	if len(timeStats) == 0 {
		return fmt.Errorf("Error: (metrics) Empty cpu stats")
	}
	// requesting all cpu times will always return an array with only one
	// time stats entry.
	timeStat := timeStats[0]
	localTime, err := getProcessCPUTime()
	if err != nil {
		return err
	}
	stats.GlobalTime = timeStat.User + timeStat.Nice + timeStat.System
	stats.GlobalWait = timeStat.Iowait
	stats.LocalTime = localTime
	return nil
}
