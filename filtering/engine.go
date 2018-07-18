// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"time"

	"github.com/iomz/go-llrp"
)

// Engine provides interface for the filtering engines
type Engine interface {
	AddSubscription(Subscriptions)
	DeleteSubscription(Subscriptions)
	Dump() string
	MarshalBinary() ([]byte, error)
	Name() string
	Search(llrp.ReadEvent) (string, []string, error) // pureIdentity, reportURIs, err
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
	"LegacyEngine": NewLegacyEngine,
}

/* internal helper func */
// timeTrack measures the time it taken from the start
func timeTrack(start time.Time, tpech chan time.Duration) {
	tpech <- time.Since(start)
}
