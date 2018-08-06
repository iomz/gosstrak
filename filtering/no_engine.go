// Copyright (c) 2017 Iori Mizutani
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

// NoEngine is a dummy engine
type NoEngine struct {
}

// AddSubscription does nothing
func (ne *NoEngine) AddSubscription(sub Subscriptions) {
}

// DeleteSubscription does nothing
func (ne *NoEngine) DeleteSubscription(sub Subscriptions) {
}

// Dump does nothing
func (ne *NoEngine) Dump() (s string) {
	return s
}

// MarshalBinary does nothing
func (ne *NoEngine) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.NoEngine")

	return buf.Bytes(), err
}

// Name returs the name of this engine type
func (ne *NoEngine) Name() string {
	return "NoEngine"
}

// Search does nothing
func (ne *NoEngine) Search(re llrp.ReadEvent) (string, []string, error) {
	return "", []string{}, fmt.Errorf("no match found for %v", re.ID)
}

// UnmarshalBinary does nothing
func (ne *NoEngine) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.NoEngine" {
		return fmt.Errorf("Wrong Filtering Engine: %s", typeOfEngine)
	}

	return
}

// NewNoEngine returns a dummy engine
func NewNoEngine(sub Subscriptions) Engine {
	return &NoEngine{}
}
