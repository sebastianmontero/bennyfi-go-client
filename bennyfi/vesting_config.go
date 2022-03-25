// The MIT License (MIT)

// Copyright (c) 2020, Digital Scarcity

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package bennyfi

import (
	"fmt"
	"math"
	"time"

	"github.com/sebastianmontero/eos-go-toolbox/util"
)

var (
	VestingPeriod     = "period_hrs"
	VestingPercentage = "percentage_x100000"
)

type VestingConfig struct {
	Config Config `json:"config"`
}

func NewNoVestingConfig() *VestingConfig {
	return &VestingConfig{
		Config: make(Config, 0),
	}
}

func NewImmediateVestingConfig() *VestingConfig {
	return NewVestingConfig(0, uint32(PercentageAdjustment))
}

func NewVestingConfig(period, percentage uint32) *VestingConfig {
	return &VestingConfig{
		Config: Config{
			{
				Key:   VestingPeriod,
				Value: FlexValueFromUint32(period),
			},
			{
				Key:   VestingPercentage,
				Value: FlexValueFromUint32(percentage),
			},
		},
	}
}

func (m *VestingConfig) HasVesting() bool {
	return m.Config.HasConfig()
}

func (m *VestingConfig) GetPeriod() uint32 {
	return m.Config.Get(VestingPeriod).Uint32()
}

func (m *VestingConfig) GetPercentage() uint32 {
	return m.Config.Get(VestingPercentage).Uint32()
}

func (m *VestingConfig) TotalCycles() uint32 {
	return uint32(math.Ceil(PercentageAdjustment / float64(m.GetPercentage())))
}

type VestingContext struct {
	VestingConfigs map[string]*VestingTracker
	VestingCycle   uint16
	VestingTime    time.Time
}

func (m *VestingContext) HasConfigs() bool {
	return len(m.VestingConfigs) > 0
}

func (m *VestingContext) Process(time time.Time, tracker *VestingTracker) {
	if m.HasConfigs() && time.After(m.VestingTime) {
		return
	}
	if !m.HasConfigs() || time.Before(m.VestingTime) {
		m.VestingConfigs = make(map[string]*VestingTracker, 0)
		m.VestingTime = time
	}
	m.VestingConfigs[tracker.DistName] = tracker
}

func (m *VestingContext) IncreaseCycle() {
	for _, tracker := range m.VestingConfigs {
		tracker.IncreaseCycle()
	}
}

type VestingTracker struct {
	VestingConfig *VestingConfig
	DistName      string
	Cycle         uint16
}

func (m *VestingTracker) HasNextVestingTime() bool {
	return m.Cycle+1 <= uint16(m.VestingConfig.TotalCycles())
}
func (m *VestingTracker) GetNextVestingTime(startTime time.Time) time.Time {
	return startTime.Add(time.Hour * (time.Duration(m.VestingConfig.GetPeriod() * uint32(m.Cycle+1))))
}
func (m *VestingTracker) IncreaseCycle() {
	m.Cycle++
}

type VestingTrackers []*VestingTracker

func (m VestingTrackers) GetContextForCycle(cycle uint16, startTime string) *VestingContext {
	st, err := util.ToTime(startTime)
	if err != nil {
		panic(fmt.Sprintf("failed parsing vesting start time, error: %v", err))
	}
	var context *VestingContext
	for c := uint16(1); c <= cycle; c++ {
		context = m.FindNext(st)
		if !context.HasConfigs() {
			panic(fmt.Sprintf("There is no vesting cycle: %v, max vesting cycle:%v", cycle, c-1))
		}
	}
	return context
}

func (m VestingTrackers) FindNext(startTime time.Time) *VestingContext {
	context := &VestingContext{}
	for _, tracker := range m {
		if tracker.HasNextVestingTime() {
			context.Process(tracker.GetNextVestingTime(startTime), tracker)
		}
	}
	context.IncreaseCycle()
	return context
}
