// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// OptimalBST struct
type OptimalBST struct {
	notificationURI string
	filterObject    *FilterObject
	matchNext       *OptimalBST
	mismatchNext    *OptimalBST
}

// AnalyzeLocality increments the locality per node for the specific id
func (obst *OptimalBST) AnalyzeLocality(id []byte, path string, lm *LocalityMap) {
	if len(path) == 0 {
		path = obst.notificationURI
	} else {
		path += "," + obst.notificationURI
	}

	// count up the traffic per node
	if _, ok := (*lm)[path]; !ok {
		(*lm)[path] = 1
	} else {
		(*lm)[path]++
	}

	if obst.filterObject.Match(id) {
		if obst.matchNext != nil {
			obst.matchNext.AnalyzeLocality(id, path, lm)
		} else {
			if _, ok := (*lm)[path+",Match"]; !ok {
				(*lm)[path+",Match"] = 1
			} else {
				(*lm)[path+",Match"]++
			}
		}
	} else if obst.mismatchNext != nil {
		obst.mismatchNext.AnalyzeLocality(id, path, lm)
	} else {
		if _, ok := (*lm)[path+",Mismatch"]; !ok {
			(*lm)[path+",Mismatch"] = 1
		} else {
			(*lm)[path+",Mismatch"]++
		}
	}
}

// Dump returs a string representation of the PatriciaTrie
func (obst *OptimalBST) Dump() string {
	writer := &bytes.Buffer{}
	obst.print(writer, 0)
	return writer.String()
}

// MarshalBinary overwrites the marshaller in gob encoding *OptimalBST
func (obst *OptimalBST) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.OptimalBST")

	// Notify
	enc.Encode(obst.notificationURI)

	// Filter
	hasFilter := obst.filterObject != nil
	enc.Encode(hasFilter)
	if hasFilter {
		err = enc.Encode(obst.filterObject)
	}

	// matchNext
	hasMatchNext := obst.matchNext != nil
	enc.Encode(hasMatchNext)
	if hasMatchNext {
		err = enc.Encode(obst.matchNext)
	}

	// mismatchNext
	hasMismatchNext := obst.mismatchNext != nil
	enc.Encode(hasMismatchNext)
	if hasMismatchNext {
		err = enc.Encode(obst.mismatchNext)
	}

	return buf.Bytes(), err
}

// Search returns a slice of notificationURI
func (obst *OptimalBST) Search(id []byte) []string {
	matches := []string{}
	if obst.filterObject.Match(id) {
		matches = append(matches, obst.notificationURI)
		if obst.matchNext != nil {
			matches = append(matches, obst.matchNext.Search(id)...)
		}
		return matches
	}
	if obst.mismatchNext != nil {
		return obst.mismatchNext.Search(id)
	}
	return matches
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (obst *OptimalBST) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.OptimalBST" {
		return errors.New("Wrong Filtering Engine: " + typeOfEngine)
	}

	// Notify
	if err = dec.Decode(&obst.notificationURI); err != nil {
		return
	}

	// FilterObject
	var hasFilterObject bool
	if err = dec.Decode(&hasFilterObject); err != nil {
		return
	}
	if hasFilterObject {
		err = dec.Decode(&obst.filterObject)
	}

	// matchNext
	var hasMatchNext bool
	if err = dec.Decode(&hasMatchNext); err != nil {
		return
	}
	if hasMatchNext {
		err = dec.Decode(&obst.matchNext)
	} else {
		obst.matchNext = nil
	}

	// mismatchNext
	var hasMismatchNext bool
	if err = dec.Decode(&hasMismatchNext); err != nil {
		return
	}
	if hasMismatchNext {
		err = dec.Decode(&obst.mismatchNext)
	} else {
		obst.mismatchNext = nil
	}

	return
}

// Internal helper methods -----------------------------------------------------

func (obst *OptimalBST) build(sub *Subscriptions, nds *Nodes) *OptimalBST {
	current := obst
	for i, nd := range *nds {
		current.filterObject = NewFilter(nd.filter, nd.offset)
		current.notificationURI = nd.notificationURI
		// if this node has subset
		if _, ok := (*sub)[nd.filter]; ok && (*sub)[nd.filter].Subset != nil {
			subset := (*sub)[nd.filter].Subset
			subnds := NewNodes(subset)
			subobst := &OptimalBST{}
			current.matchNext = subobst.build(subset, subnds)
		} else {
			current.matchNext = nil
		}
		if i+1 < len(*nds) {
			current.mismatchNext = &OptimalBST{}
			current = current.mismatchNext
		} else {
			current.mismatchNext = nil
		}
	}
	return obst
}

func (obst *OptimalBST) equal(want *OptimalBST) (ok bool, got *OptimalBST, wanted *OptimalBST) {
	if obst.notificationURI != want.notificationURI ||
		!reflect.DeepEqual(obst.filterObject, want.filterObject) {
		return false, obst, want
	}
	if want.matchNext != nil {
		if obst.matchNext == nil {
			return false, nil, want.matchNext
		}
		res, cgot, cwanted := obst.matchNext.equal(want.matchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	if want.mismatchNext != nil {
		if obst.mismatchNext == nil {
			return false, nil, want.mismatchNext
		}
		res, cgot, cwanted := obst.mismatchNext.equal(want.mismatchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	return true, nil, nil
}

func (obst *OptimalBST) print(writer io.Writer, indent int) {
	var n string
	if len(obst.notificationURI) != 0 {
		n = "-> " + obst.notificationURI
	}
	fmt.Fprintf(writer, "%s--%s %s\n", strings.Repeat(" ", indent), obst.filterObject.ToString(), n)
	if obst.matchNext != nil {
		obst.matchNext.print(writer, indent+2)
	}
	if obst.mismatchNext != nil {
		obst.mismatchNext.print(writer, indent+2)
	}
}

// BuildOptimalBST builds OptimalBST from Subscriptions
// returns the pointer to the node node
func BuildOptimalBST(sub *Subscriptions) *OptimalBST {
	// make subsets to the child subscriptions of the corresponding parents
	sub.linkSubset()

	// recalculate the cumulative evs if it has subset subscriptions
	for _, info := range *sub {
		if info.Subset != nil {
			info.EntropyValue += recalculateEntropyValue(info.Subset)
		}
	}

	nds := NewNodes(sub)
	obst := &OptimalBST{}
	return obst.build(sub, nds)
}
