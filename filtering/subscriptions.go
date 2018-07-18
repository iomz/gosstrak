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

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		reportURI := strings.ToUpper(record[0])
		if !strings.HasPrefix(reportURI, "http") {
			continue
		}
		for i := 1; i < len(record); i++ {
			pat := strings.ToLower(record[i])
			if strings.HasPrefix(pat, "urn:epc:pat:") {
				if _, ok := sub[reportURI]; !ok {
					sub[reportURI] = []string{}
				}
				sub[reportURI] = append(sub[reportURI], pat)
			}
		}
	}
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

// Clone retuns a new copy of subscriptions
func (sub ByteSubscriptions) Clone() ByteSubscriptions {
	clone := ByteSubscriptions{}

	for _, fs := range sub.Keys() {
		psub := PartialSubscription{}
		psub.Offset = sub[fs].Offset
		psub.ReportURI = sub[fs].ReportURI
		if len(sub[fs].Subset) == 0 {
			psub.Subset = ByteSubscriptions{}
		} else {
			psub.Subset = sub[fs].Subset.Clone()
		}
		clone[fs] = &psub
	}
	return clone
}

// Dump retuns a string representation of the ByteSubscriptions
func (sub ByteSubscriptions) Dump() string {
	writer := &bytes.Buffer{}
	sub.print(writer, 0)
	return writer.String()
}

// Keys return a slice of keys in m
func (sub ByteSubscriptions) Keys() []string {
	ks := make([]string, len(sub))
	i := 0
	for k := range sub {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	return ks
}

// linkSubset finds subsets and nest them under the parents
func (sub ByteSubscriptions) linkSubset() {
	type element struct {
		filter    string
		offset    int
		reportURI string
	}

	var elements []*element
	for _, fs := range sub.Keys() {
		elements = append(elements, &element{
			filter:    fs,
			offset:    sub[fs].Offset,
			reportURI: sub[fs].ReportURI,
		})
	}

	for _, e := range elements {
		for _, fs := range sub.Keys() {
			psub := sub[fs]
			linkCandidate := e.filter
			// check if fs is a subset of the linkCandidate
			if strings.HasPrefix(fs, linkCandidate) &&
				fs != linkCandidate { // they shouldn't be the same
				// if there is no subset already
				if len(sub[linkCandidate].Subset) == 0 {
					sub[linkCandidate].Subset = ByteSubscriptions{}
					sub[linkCandidate].Subset[fs[len(linkCandidate):]] = &PartialSubscription{
						Offset:    psub.Offset + len(linkCandidate),
						ReportURI: psub.ReportURI,
					}
				} else {
					sub[linkCandidate].Subset[fs[len(linkCandidate):]] = &PartialSubscription{
						Offset:    psub.Offset + len(linkCandidate),
						ReportURI: psub.ReportURI,
					}
				}
				// recursively link the subset
				sub[linkCandidate].Subset.linkSubset()
				// finaly delete the filter from the upper ByteSubscriptions
				delete(sub, fs)
			}
		}
	}
	return
}

// MarshalBinary overwrites the marshaller in gob encoding *ByteSubscriptions
func (sub *ByteSubscriptions) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Size of ByteSubscriptions
	enc.Encode(len(*sub))
	for _, fs := range (*sub).Keys() {
		// FilterString
		enc.Encode(fs)
		// PartialSubscription
		err = enc.Encode((*sub)[fs])
	}

	return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding List
func (sub *ByteSubscriptions) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Size of List
	var subSize int
	if err = dec.Decode(&subSize); err != nil {
		return
	}

	for i := 0; i < subSize; i++ {
		// FilterString
		var fs string
		if err = dec.Decode(&fs); err != nil {
			return
		}

		/// PartialSubscription
		psub := PartialSubscription{
			Subset: ByteSubscriptions{},
		}
		err = dec.Decode(&psub)

		// add to the sub
		(*sub)[fs] = &psub
	}

	return
}

// LoadFiltersFromCSVFile takes a csv file name and generate ByteSubscriptions
func LoadFiltersFromCSVFile(f string) ByteSubscriptions {
	sub := ByteSubscriptions{}
	fp, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if len(record) == 2 {
			// Default case
			// prefix as key, *filtering.PartialSubscription as value
			sub[record[1]] = &PartialSubscription{
				Offset:    0,
				ReportURI: record[0],
				Subset:    ByteSubscriptions{},
			}
		}
	}
	return sub
}

func (sub ByteSubscriptions) print(writer io.Writer, indent int) {
	for _, fs := range sub.Keys() {
		fmt.Fprintf(writer, "%s--%s %v %s\n", strings.Repeat(" ", indent), fs, sub[fs].Offset, sub[fs].ReportURI)
		ss := sub[fs].Subset
		if len(ss) != 0 {
			ss.print(writer, indent+2)
		}
	}
}
