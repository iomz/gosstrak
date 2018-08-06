// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/iomz/go-llrp"
)

// DummyEngine is a dummy engine
type DummyEngine struct {
}

// AddSubscription does nothing
func (de *DummyEngine) AddSubscription(sub Subscriptions) {
}

// DeleteSubscription does nothing
func (de *DummyEngine) DeleteSubscription(sub Subscriptions) {
}

// Dump does nothing
func (de *DummyEngine) Dump() (s string) {
	return s
}

// MarshalBinary does nothing
func (de *DummyEngine) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.DummyEngine")

	return buf.Bytes(), err
}

// Name returs the name of this engine type
func (de *DummyEngine) Name() string {
	return "DummyEngine"
}

// Search does nothing
func (de *DummyEngine) Search(re llrp.ReadEvent) (string, []string, error) {
	return "", []string{}, fmt.Errorf("no match found for %v", re.ID)
}

// UnmarshalBinary does nothing
func (de *DummyEngine) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.DummyEngine" {
		return fmt.Errorf("Wrong Filtering Engine: %s", typeOfEngine)
	}

	return
}

// NewDummyEngine returns a dummy engine
func NewDummyEngine(sub Subscriptions) Engine {
	return &DummyEngine{}
}
