// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	//"strconv"
	"strings"

	"github.com/iomz/gosstrak/tdt"
)

// ByteSubscriptions contains filter string as key and PartialSubscription as value
type ByteSubscriptions map[string]*PartialSubscription

// PartialSubscription contains reportURI and pValue for a filter
type PartialSubscription struct {
	Offset    int
	ReportURI string
	Subset    ByteSubscriptions
}

// Subscriptions contains a slice of urn:epc:pat as values and a URI to report events as keys
type Subscriptions map[string][]string

// Clone retuns a new copy of subscriptions
func (sub Subscriptions) Clone() Subscriptions {
	clone := Subscriptions{}
	for _, f := range sub.Keys() {
		clone[f] = []string{}
		for _, dest := range sub[f] {
			clone[f] = append(clone[f], dest)
		}
	}
	return clone
}

// Keys return a slice of keys in Subscriptions
func (sub Subscriptions) Keys() []string {
	ks := make([]string, len(sub))
	i := 0
	for k := range sub {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	return ks
}

// ToByteSubscriptions preprocess the subscription and convert them in bytes
func (sub Subscriptions) ToByteSubscriptions() ByteSubscriptions {
	bsub := ByteSubscriptions{}
	for reportURI, patterns := range sub {
		for _, pat := range patterns {
			tf := strings.Split(strings.TrimPrefix(pat, "urn:epc:pat:"), ":")
			if len(tf) != 2 { // should only containts a type and fields
				continue
			}
			fields := strings.Split(strings.ToUpper(tf[1]), ".")
			pfs, err := tdt.MakePrefixFilterString(tf[0], fields)
			if err != nil {
				log.Print(err)
			}
			bsub[pfs] = &PartialSubscription{
				Offset:    0,
				ReportURI: reportURI,
				Subset:    ByteSubscriptions{},
			}
		}
	}
	return bsub
}

// LoadSubscriptionsFromCSVFile takes a csv file name and returns Subscriptions
func LoadSubscriptionsFromCSVFile(f string) Subscriptions {
	sub := Subscriptions{}
	fp, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	numFilters := 0
	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = false
	reader.FieldsPerRecord = -1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		reportURI := strings.ToLower(record[0])
		if !strings.HasPrefix(reportURI, "http") {
			continue
		}
		for i := 1; i < len(record); i++ {
			pat := record[i]
			if strings.HasPrefix(strings.ToLower(pat), "urn:epc:pat:") {
				if _, ok := sub[reportURI]; !ok {
					sub[reportURI] = []string{}
				}
				sub[reportURI] = append(sub[reportURI], pat)
				numFilters++
			}
		}
	}
	log.Printf("%v filtering patterns loaded from %s", numFilters, f)
	return sub
}

// MarshalBinary overwrites the marshaller in gob encoding *PartialSubscription
func (psub *PartialSubscription) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Offset
	enc.Encode(psub.Offset)

	// ReportURI
	enc.Encode(psub.ReportURI)

	// Subset
	enc.Encode(psub.Subset)

	//buf.Encode
	return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PartialSubscription
func (psub *PartialSubscription) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Offset
	if err = dec.Decode(&psub.Offset); err != nil {
		return
	}

	// ReportURI
	if err = dec.Decode(&psub.ReportURI); err != nil {
		return
	}

	// Subset
	if err = dec.Decode(&psub.Subset); err != nil {
		return
	}

	return
}

// Dump retuns a string representation of the ByteSubscriptions
func (bsub ByteSubscriptions) Dump() string {
	writer := &bytes.Buffer{}
	bsub.print(writer, 0)
	return writer.String()
}

// Keys return a slice of keys in m
func (bsub ByteSubscriptions) Keys() []string {
	ks := make([]string, len(bsub))
	i := 0
	for k := range bsub {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	return ks
}

// linkSubset finds subsets and nest them under the parents
func (bsub ByteSubscriptions) linkSubset() {
	type element struct {
		filter    string
		offset    int
		reportURI string
	}

	var elements []*element
	for _, fs := range bsub.Keys() {
		elements = append(elements, &element{
			filter:    fs,
			offset:    bsub[fs].Offset,
			reportURI: bsub[fs].ReportURI,
		})
	}

	for _, e := range elements {
		for _, fs := range bsub.Keys() {
			psub := bsub[fs]
			linkCandidate := e.filter
			// check if fs is a subset of the linkCandidate
			if strings.HasPrefix(fs, linkCandidate) &&
				fs != linkCandidate { // they shouldn't be the same
				// if there is no subset already
				if len(bsub[linkCandidate].Subset) == 0 {
					bsub[linkCandidate].Subset = ByteSubscriptions{}
					bsub[linkCandidate].Subset[fs[len(linkCandidate):]] = &PartialSubscription{
						Offset:    psub.Offset + len(linkCandidate),
						ReportURI: psub.ReportURI,
					}
				} else {
					bsub[linkCandidate].Subset[fs[len(linkCandidate):]] = &PartialSubscription{
						Offset:    psub.Offset + len(linkCandidate),
						ReportURI: psub.ReportURI,
					}
				}
				// recursively link the subset
				bsub[linkCandidate].Subset.linkSubset()
				// finaly delete the filter from the upper ByteSubscriptions
				delete(bsub, fs)
			}
		}
	}
	return
}

// makePatriciaSubset finds subsets and nest them under the parents
func (bsub ByteSubscriptions) makePatriciaSubset() {
	// store the currentOffset of this ByteSubscriptions
	var currentOffset int

	// loop until there's no lcp in this ByteSubscriptions
	for {
		commonPrefix, subsetSize := findPartialLCP(bsub.Keys())
		if subsetSize < 1 {
			break
		}
		log.Printf("commonPrefix: %v", commonPrefix)
		// sbusb is a double pointer to the map
		// but a pointer to ByteSubscriptions
		// initialize an empty instance and a pointer
		sbsub := &ByteSubscriptions{}
		// just a pointer, set later in the loop
		var superset *PartialSubscription
		// iterate the bsub with this commonPrefix
		for i, fs := range bsub.Keys() {
			if i == 0 {
				// set currentOffset
				// All Offset is the same in a ByteSubscriptions
				// FIXME: Offset should be a member of ByteSubscriptions
				//        and not of PartialSubscription
				currentOffset = bsub[fs].Offset
			}
			if fs == commonPrefix {
				// if the commonPrefix itself is a subscription
				// make this a superset of subscirptions with current commonPrefix
				superset = &PartialSubscription{
					Offset:    currentOffset,
					ReportURI: bsub[fs].ReportURI,
					Subset:    bsub[fs].Subset,
				}
				// delete the superset
				delete(bsub, fs)
			} else if strings.HasPrefix(fs, commonPrefix) {
				// if this is PartialSubscription is a subset of this commonPrefix
				// check if this is not a superset?
				(*sbsub)[fs[len(commonPrefix):]] = &PartialSubscription{
					Offset:    currentOffset + len(commonPrefix),
					ReportURI: bsub[fs].ReportURI,
					Subset:    bsub[fs].Subset,
				}
				// delete the subset
				delete(bsub, fs)
			}
		}
		if superset == nil {
			// if there's no superset yet, create an empty PartialSubscription
			superset = &PartialSubscription{
				Offset: currentOffset,
			}
		}
		// group the subset and insert it to the superset
		superset.Subset = *sbsub
		// insert the superset to the current bsub
		bsub[commonPrefix] = superset
	}
	// finish finding lcp in this depth
	// now dive to the next depth and link recursively
	for _, fs := range bsub.Keys() {
		if len(bsub[fs].Subset) != 0 {
			bsub[fs].Subset.linkSubset()
		}
	}
	return
}

// print used for Dump()
func (bsub ByteSubscriptions) print(writer io.Writer, indent int) {
	for _, fs := range bsub.Keys() {
		fmt.Fprintf(writer, "%s--%s %v %s\n", strings.Repeat(" ", indent), fs, bsub[fs].Offset, bsub[fs].ReportURI)
		ss := bsub[fs].Subset
		if len(ss) != 0 {
			ss.print(writer, indent+2)
		}
	}
}
