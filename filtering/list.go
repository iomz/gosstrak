package filtering

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

// List is a slice of pointers to ExactMatch
type List []*ExactMatch

// ExactMatch is a raw filter directly taken from Subscriptions
type ExactMatch struct {
	notificationURI string
	filter          *FilterObject
}

// MarshalBinary overwrites the marshaller in gob encoding *List
func (list *List) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.List")

	// Size of List
	enc.Encode(len(*list))
	for _, em := range *list {
		// Notify
		enc.Encode(em.notificationURI)
		// Filter
		err = enc.Encode(em.filter)
	}

	return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding List
func (list *List) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.List" {
		return errors.New("Wrong Filtering Engine: " + typeOfEngine)
	}

	// Size of List
	var listSize int
	if err = dec.Decode(&listSize); err != nil {
		return
	}

	for i := 0; i < listSize; i++ {
		em := ExactMatch{}
		// Notify
		if err = dec.Decode(&em.notificationURI); err != nil {
			return
		}
		// Filter
		err = dec.Decode(&em.filter)
		*list = append(*list, &em)
	}

	return
}

// AnalyzeLocality increments the locality per node for the specific id
func (list *List) AnalyzeLocality(id []byte, prefix string, lm *LocalityMap) {
}

// Dump returs a string representation of the PatriciaTrie
func (list *List) Dump() string {
	writer := &bytes.Buffer{}
	for _, em := range *list {
		fmt.Fprintf(writer, "--%s %s\n", em.filter.ToString(), em.notificationURI)
	}
	return writer.String()
}

// Search returns a slice of notificationURI
func (list *List) Search(id []byte) (matches []string) {
	for _, em := range *list {
		if em.filter.Match(id) {
			matches = append(matches, em.notificationURI)
		}
	}
	return
}

// BuildList builds a simple list of filters from filter.Subscriptions
// returns the pointer to the slice of ExactMatch struct
func BuildList(sub Subscriptions) *List {
	list := make(List, len(sub))

	// store ExactMatch in sorted order from sub
	for i, fs := range sub.keys() {
		list[i] = &ExactMatch{sub[fs].NotificationURI, NewFilter(fs, 0)}
	}

	return &list
}