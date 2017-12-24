// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

// Engine provides interface for the filtering engines
type Engine interface {
	AnalyzeLocality(id []byte, prefix string, lm *LocalityMap)
	Dump() string
	MarshalBinary() ([]byte, error)
	Search(id []byte) []string
	UnmarshalBinary(data []byte) error
}

// NotifyMap contains notificationURI sring as key and slice of ids in [][]byte
type NotifyMap map[string][][]byte
