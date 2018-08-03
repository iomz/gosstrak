// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package monitoring

// StatMessageType is a type for StatMessage
type StatMessageType int

const (
	// Traffic message
	Traffic StatMessageType = iota
	// EngineThroughput message
	EngineThroughput
	// SelectedEngine message
	SelectedEngine
)

// StatMessage carries stat
type StatMessage struct {
	Type  StatMessageType
	Value []interface{}
	Name  string
}
