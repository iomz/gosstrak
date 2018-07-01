// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

// Engine provides interface for the filtering engines
type Engine interface {
	AddSubscription(Subscriptions)
	AnalyzeLocality([]byte, string, *LocalityMap)
	DeleteSubscription(Subscriptions)
	Dump() string
	MarshalBinary() ([]byte, error)
	Search([]byte) []string
	UnmarshalBinary([]byte) error
}

// EngineConstructor is a function signature for engine constructors
type EngineConstructor func(Subscriptions) Engine

// Engines is a map of Engne's name and its constructor
type Engines map[string]EngineConstructor

// AvailableEngines is a map of EngineConstructor with the engine's names as keys
var AvailableEngines = Engines{
	"List":         NewList,
	"PatriciaTrie": NewPatriciaTrie,
	"SplayTree":    NewSplayTree,
}
