// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

// ReadEvent is the struct to hold data on RFTags
type ReadEvent struct {
	ID []byte
	PC []byte
}
