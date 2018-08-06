// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"fmt"
	"log"
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
	//"List":         NewList,
	"PatriciaTrie": NewPatriciaTrie,
	//"SplayTree":    NewSplayTree,
	//"LegacyEngine": NewLegacyEngine,
	"DummyEngine": NewDummyEngine,
	//"NoEngine": NewNoEngine, // comment out like this to disable an engine
}

/* internal helper func */
// timeTrack measures the time it taken from the start
func timeTrack(start time.Time, tpech chan time.Duration) {
	tpech <- time.Since(start)
}

// findPartialLCP discovers the the commonPrefix of which the majority share
// reutrns the prefix and the size of the majority in len(l)
func findPartialLCP(l []string) (string, int) {
	if len(l) < 2 {
		return "", 0
	}
	log.Printf("findPartialLCP: %q", l)
	m := lcp(l)
	if len(m) != 0 {
		return m, len(l)
	}
	zeros := []string{}
	ones := []string{}
	for _, pl := range l {
		if len(pl) < 1 {
			continue
		}
		switch pl[0] {
		case '0':
			zeros = append(zeros, pl[1:])
		case '1':
			ones = append(ones, pl[1:])
		}
	}
	if len(zeros) > 1 {
		zeroPL, zeroSize := findPartialLCP(zeros)
		if len(zeroPL) != 0 && zeroSize > 1 {
			return "0" + zeroPL, zeroSize
		}

	}
	if len(ones) > 1 {
		onePL, oneSize := findPartialLCP(ones)
		if len(onePL) != 0 && oneSize > 1 {
			return "1" + onePL, oneSize
		}
	}

	return "", 0
}

// lcp finds the longest common prefix of the input strings.
// It compares by bytes instead of runes (Unicode code points).
// It's up to the caller to do Unicode normalization if desired
// (e.g. see golang.org/x/text/unicode/norm).
func lcp(l []string) string {
	// Special cases first
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0]
	}
	// LCP of min and max (lexigraphically)
	// is the LCP of the whole set.
	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}

// getNextBit gets the bit of specific bit offset
func getNextBit(id []byte, nbo int) (rune, error) {
	o := nbo / ByteLength
	// No more bit in the ID
	if len(id) == o && nbo%ByteLength == 0 {
		return 'x', nil
	}
	if len(id) <= o {
		return '?', fmt.Errorf("invalid offset: %v", nbo)
	}
	if (uint8(id[o])>>uint8((7-(nbo%ByteLength))))%2 == 0 {
		return '0', nil
	}
	return '1', nil
}
