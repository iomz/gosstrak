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

// AddSubscription adds a set of subscriptions if not exists yet
func (obst *OptimalBST) AddSubscription(sub Subscriptions) {
	for _, fs := range sub.keys() {
		obst.add(fs, sub[fs].NotificationURI)
	}
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

// DeleteSubscription deletes a set of subscriptions if already exist
func (obst *OptimalBST) DeleteSubscription(sub Subscriptions) {
	for _, fs := range sub.keys() {
		obst.delete(fs, sub[fs].NotificationURI)
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

// add a set of subscriptions if not exists yet
func (obst *OptimalBST) add(fs string, notificationURI string) {
	if strings.HasPrefix(fs, obst.filterObject.String) { // fs \in obst.FilterObject.String
		if fs == obst.filterObject.String { // the identical filter
			// if the notificationURI is different, update it
			if obst.notificationURI != notificationURI {
				obst.notificationURI = notificationURI
			}
		} else { // it's a subset (matching branch) of the current node, and there's no matchNext node
			if obst.matchNext == nil {
				obst.matchNext = &OptimalBST{}
				obst.matchNext.filterObject = NewFilter(fs[obst.filterObject.Size:], obst.filterObject.Offset+obst.filterObject.Size)
				obst.matchNext.notificationURI = notificationURI
			} else { // if there's already matchNext node
				obst.matchNext.add(fs[obst.filterObject.Size:], notificationURI)
			}
		}
	} else { // doesn't match with the current node, traverse the mismatchNext node
		if obst.mismatchNext == nil { // there's no mismatchNext node
			obst.mismatchNext = &OptimalBST{}
			obst.mismatchNext.filterObject = NewFilter(fs, obst.filterObject.Offset)
			obst.mismatchNext.notificationURI = notificationURI
		} else { // if there's already mismatchNext node
			obst.mismatchNext.add(fs, notificationURI)
		}
	}
	return
}

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

// delete a set of subscriptions if not exists yet
func (obst *OptimalBST) delete(fs string, notificationURI string) {
	if strings.HasPrefix(fs, obst.filterObject.String) { // fs \in obst.FilterObject.String
		if fs == obst.filterObject.String { // this node is to delete
			if obst.matchNext == nil && obst.mismatchNext == nil { // something wrong
			} else if obst.matchNext != nil { // if there is subset, keep the node as an aggregation node
				if obst.matchNext.mismatchNext != nil {
					obst.notificationURI = ""
				} else { // if none other mismatch branch, concatenate the matchNext with to-be-deleted node
					obst.filterObject = NewFilter(fs+obst.matchNext.filterObject.String, obst.filterObject.Offset)
					obst.notificationURI = obst.matchNext.notificationURI
					obst.matchNext = obst.matchNext.matchNext
				}
			} else if obst.mismatchNext != nil { // replace this node with mismatchNext
				obst.filterObject = obst.mismatchNext.filterObject
				obst.notificationURI = obst.mismatchNext.notificationURI
				obst.matchNext = obst.mismatchNext.matchNext
				obst.mismatchNext = obst.mismatchNext.mismatchNext
			}
		} else { // it's a subset (matching branch) of the current node, and there's no matchNext node
			if obst.matchNext != nil { // if there's a matchNext node
				if fs[obst.filterObject.Size:] == obst.matchNext.filterObject.String &&
					obst.matchNext.matchNext == nil && obst.matchNext.mismatchNext == nil { // the matchNext is to delete
					obst.matchNext = nil
				} else {
					obst.matchNext.delete(fs[obst.filterObject.Size:], notificationURI)
				}
			}
		}
	} else { // doesn't match with the current node, traverse the mismatchNext node
		if obst.mismatchNext != nil { // there's a mismatchNext node
			if fs == obst.mismatchNext.filterObject.String &&
				obst.mismatchNext.matchNext == nil && obst.mismatchNext.mismatchNext == nil {
				obst.mismatchNext = nil
			} else {
				obst.mismatchNext.delete(fs, notificationURI)
			}
		}
	}
	return
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
	fmt.Fprintf(writer, "--%s %s\n", obst.filterObject.ToString(), n)
	if obst.matchNext != nil {
		fmt.Fprintf(writer, "%sok", strings.Repeat(" ", indent+2))
		obst.matchNext.print(writer, indent+2)
	}
	if obst.mismatchNext != nil {
		fmt.Fprintf(writer, "%sng", strings.Repeat(" ", indent+2))
		obst.mismatchNext.print(writer, indent+2)
	}
}

// BuildOptimalBST builds OptimalBST from Subscriptions
// returns the pointer to the node node
func BuildOptimalBST(sub *Subscriptions) *OptimalBST {
	// make subsets to the child subscriptions of the corresponding parents
	sub.linkSubset()

	// recalculate the cumulative evs if it has subset subscriptions
	for _, fs := range (*sub).keys() {
		info := (*sub)[fs]
		if info.Subset != nil {
			info.EntropyValue += recalculateEntropyValue(info.Subset)
		}
	}

	nds := NewNodes(sub)
	obst := &OptimalBST{}
	return obst.build(sub, nds)
}
