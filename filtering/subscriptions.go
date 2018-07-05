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
	"os"
	"sort"
	"strconv"
	"strings"
)

// Subscriptions contains filter string as key and Info as value
type Subscriptions map[string]*Info

// Info contains notificationURI and pValue for a filter
type Info struct {
	Offset          int
	NotificationURI string
	EntropyValue    float64
	Subset          Subscriptions
}

// MarshalBinary overwrites the marshaller in gob encoding *Info
func (info *Info) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Offset
	enc.Encode(info.Offset)

	// NotificationURI
	enc.Encode(info.NotificationURI)

	// EntropyValue
	enc.Encode(info.EntropyValue)

	// Subset
	enc.Encode(info.Subset)

	//buf.Encode
	return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *Info
func (info *Info) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Offset
	if err = dec.Decode(&info.Offset); err != nil {
		return
	}

	// NotificationURI
	if err = dec.Decode(&info.NotificationURI); err != nil {
		return
	}

	// EntropyValue
	if err = dec.Decode(&info.EntropyValue); err != nil {
		return
	}
	// Subset
	if err = dec.Decode(&info.Subset); err != nil {
		return
	}

	return
}

// Clone retuns a new copy of subscriptions
func (sub Subscriptions) Clone() Subscriptions {
	clone := Subscriptions{}

	for _, fs := range sub.Keys() {
		info := Info{}
		info.Offset = sub[fs].Offset
		info.NotificationURI = sub[fs].NotificationURI
		info.EntropyValue = sub[fs].EntropyValue
		if len(sub[fs].Subset) == 0 {
			info.Subset = Subscriptions{}
		} else {
			info.Subset = sub[fs].Subset.Clone()
		}
		clone[fs] = &info
	}
	return clone
}

// Dump retuns a string representation of the Subscriptions
func (sub Subscriptions) Dump() string {
	writer := &bytes.Buffer{}
	sub.print(writer, 0)
	return writer.String()
}

// Keys return a slice of keys in m
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

// linkSubset finds subsets and nest them under the parents
func (sub Subscriptions) linkSubset() {
	nds := NewNodes(sub)
	for _, nd := range nds {
		for _, fs := range sub.Keys() {
			info := sub[fs]
			linkCandidate := nd.filter
			// check if fs is a subset of the linkCandidate
			if strings.HasPrefix(fs, linkCandidate) &&
				fs != linkCandidate { // they shouldn't be the same
				// if there is no subset already
				if len(sub[linkCandidate].Subset) == 0 {
					sub[linkCandidate].Subset = Subscriptions{}
					sub[linkCandidate].Subset[fs[len(linkCandidate):]] = &Info{
						Offset:          info.Offset + len(linkCandidate),
						NotificationURI: info.NotificationURI,
						EntropyValue:    info.EntropyValue,
					}
				} else {
					sub[linkCandidate].Subset[fs[len(linkCandidate):]] = &Info{
						Offset:          info.Offset + len(linkCandidate),
						NotificationURI: info.NotificationURI,
						EntropyValue:    info.EntropyValue,
					}
				}
				// recursively link the subset
				sub[linkCandidate].Subset.linkSubset()
				// finaly delete the filter from the upper Subscriptions
				delete(sub, fs)
			}
		}
	}
	return
}

// MarshalBinary overwrites the marshaller in gob encoding *Subscriptions
func (sub *Subscriptions) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Size of Subscriptions
	enc.Encode(len(*sub))
	for _, fs := range (*sub).Keys() {
		// FilterString
		enc.Encode(fs)
		// Info
		err = enc.Encode((*sub)[fs])
	}

	return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding List
func (sub *Subscriptions) UnmarshalBinary(data []byte) (err error) {
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

		/// Info
		info := Info{
			Subset: Subscriptions{},
		}
		err = dec.Decode(&info)

		// add to the sub
		(*sub)[fs] = &info
	}

	return
}

// LoadFiltersFromCSVFile takes a csv file name and generate Subscriptions
func LoadFiltersFromCSVFile(f string) Subscriptions {
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
		if len(record) < 3 {
			// Default case
			// prefix as key, *filtering.Info as value
			sub[record[1]] = &Info{
				Offset:          0,
				NotificationURI: record[0],
				EntropyValue:    0,
				Subset:          Subscriptions{},
			}
		} else {
			// For OptimalBST, filter with EntropyValue
			// prefix as key, *filtering.Info as value
			pValue, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				panic(err)
			}
			fs := record[1]
			uri := record[0]
			sub[fs] = &Info{
				Offset:          0,
				NotificationURI: uri,
				EntropyValue:    pValue,
				Subset:          Subscriptions{},
			}
		}
	}
	return sub
}

func recalculateEntropyValue(sub Subscriptions) float64 {
	ev := float64(0)
	for _, fs := range sub.Keys() {
		ss := sub[fs].Subset
		if len(ss) != 0 {
			sub[fs].EntropyValue += recalculateEntropyValue(ss)
		}
		ev += sub[fs].EntropyValue
	}
	return ev
}

func (sub Subscriptions) print(writer io.Writer, indent int) {
	for _, fs := range sub.Keys() {
		fmt.Fprintf(writer, "%s--%s %f\n", strings.Repeat(" ", indent), fs, sub[fs].EntropyValue)
		ss := sub[fs].Subset
		if len(ss) != 0 {
			ss.print(writer, indent+2)
		}
	}
}
