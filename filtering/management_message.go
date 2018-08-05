// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

// ManagementMessageType is to indicate the type of ManagementMessage
type ManagementMessageType int

// ManagementMessage types
const (
	AddSubscription ManagementMessageType = iota
	DeleteSubscription
	OnEngineGenerated
	DeployEngine
	TrafficStatus
	EngineStatus
	SelectedEngine
	SimulationStat
)

// ManagementMessage holds management action for the EngineFactory
type ManagementMessage struct {
	Type                    ManagementMessageType
	Pattern                 string
	ReportURI               string
	EngineGeneratorInstance *EngineGenerator
	CurrentThroughput       float64
	EventCount              int64
	MatchedCount            int64
	EngineName              string
}
