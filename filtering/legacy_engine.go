package filtering

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strings"

	"github.com/iomz/go-llrp"
	"github.com/iomz/gosstrak/tdt"
)

// LegacyEngine is a engine based-on text match
type LegacyEngine struct {
	filters Subscriptions
	tdtCore *tdt.Core
}

// AddSubscription adds a set of subscriptions if not exists yet
func (le *LegacyEngine) AddSubscription(sub Subscriptions) {
	for _, f := range sub.Keys() {
		if reportURIs, ok := le.filters[f]; !ok {
			le.filters[f] = sub[f]
		} else {
			for _, dest := range sub[f] {
				if stringIndexInSlice(dest, reportURIs) < 0 {
					le.filters[f] = append(le.filters[f], dest)
				}
			}
		}
	}
}

// DeleteSubscription deletes a set of subscriptions if already exist
func (le *LegacyEngine) DeleteSubscription(sub Subscriptions) {
	for _, f := range sub.Keys() {
		if reportURIs, ok := le.filters[f]; !ok {
			continue
		} else {
			for _, dest := range sub[f] {
				if i := stringIndexInSlice(dest, reportURIs); i > -1 {
					le.filters[f] = append(le.filters[f][:i], le.filters[f][i+1:]...)
				}
			}
			if len(le.filters[f]) == 0 {
				delete(le.filters, f)
			}
		}
	}
}

// Dump returs a string representation of the PatriciaTrie
func (le *LegacyEngine) Dump() string {
	writer := &bytes.Buffer{}
	for _, f := range le.filters.Keys() {
		fmt.Fprintf(writer, "--%s %q\n", f, le.filters[f])
	}
	return writer.String()
}

// MarshalBinary overwrites the marshaller in gob encoding LegacyEngine
func (le *LegacyEngine) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// type of Engine
	enc.Encode("Engine:filtering.LegacyEngine")

	// size of LegacyEngine
	enc.Encode(len(le.filters.Keys()))
	for _, f := range le.filters.Keys() {
		// filter
		enc.Encode(f)
		// size of reportURIs
		enc.Encode(len(le.filters[f]))
		for _, reportURI := range le.filters[f] {
			enc.Encode(reportURI)
		}
	}

	return buf.Bytes(), err
}

// Name returs the name of this engine type
func (le *LegacyEngine) Name() string {
	return "LegacyEngine"
}

// Search returns a pureIdentity of the llrp.ReadEvent if found any subscription without err
func (le *LegacyEngine) Search(re llrp.ReadEvent) (pureIdentity string, reportURIs []string, err error) {
	// translate the readevent to a PureIdentity
	pureIdentity, err = le.tdtCore.Translate(re.PC, re.ID)
	if err != nil {
		return
	}

	for reportURI, patterns := range le.filters {
		for _, pattern := range patterns {
			seq := strings.Split(pattern, ":")
			if len(seq) != 5 {
				continue
			}
			patternType := seq[3]
			pattern := seq[4]

			switch patternType {
			case "giai-96":
				fields := strings.Split(seq[4], ".")
				// remove filter value in tag uri to match with the received PureIdentity
				pattern = "giai:" + strings.Join(fields[1:], ".")
			case "grai-96":
				fields := strings.Split(seq[4], ".")
				// remove filter value in tag uri to match with the received PureIdentity
				pattern = "grai:" + strings.Join(fields[1:], ".")
			case "sgtin-96":
				fields := strings.Split(seq[4], ".")
				// remove filter value in tag uri to match with the received PureIdentity
				pattern = "sgtin:" + strings.Join(fields[1:], ".")
			case "sscc-96":
				fields := strings.Split(seq[4], ".")
				// remove filter value in tag uri to match with the received PureIdentity
				pattern = "sscc:" + strings.Join(fields[1:], ".")
			case "iso17363":
				pattern = patternType + ":" + strings.Replace(pattern, ".", "", -1)
			case "iso17365":
				pattern = patternType + ":" + strings.Replace(pattern, ".", "", -1)
			}
			if strings.HasPrefix(strings.TrimPrefix(pureIdentity, "urn:epc:id:"), pattern) {
				reportURIs = append(reportURIs, reportURI)
			}
		}
	}
	if len(reportURIs) == 0 {
		log.Printf("%v not found", pureIdentity)
		return pureIdentity, reportURIs, fmt.Errorf("no match found for %v", pureIdentity)
	}
	return
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding LegacyEngine
func (le *LegacyEngine) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.LegacyEngine" {
		return fmt.Errorf("Wrong Filtering Engine: %s", typeOfEngine)
	}

	// Size of LegacyEngine
	var legacyEngineSize int
	if err = dec.Decode(&legacyEngineSize); err != nil {
		return
	}

	for i := 0; i < legacyEngineSize; i++ {
		var f string
		// filter
		if err = dec.Decode(&f); err != nil {
			return
		}
		var reportURIsSize int
		if err = dec.Decode(&reportURIsSize); err != nil {
			return
		}
		var reportURIs []string
		for j := 0; j < legacyEngineSize; j++ {
			// reportURI
			var dest string
			err = dec.Decode(&dest)
			reportURIs = append(reportURIs, dest)
		}
	}

	// tdt.Core
	le.tdtCore = tdt.NewCore()

	return
}

// NewLegacyEngine builds a LegacyEngine
func NewLegacyEngine(sub Subscriptions) Engine {
	// initialize LegacyEngine
	le := &LegacyEngine{}

	// load up the subscriptions
	le.filters = sub.Clone()

	// initialize tdt.Core
	le.tdtCore = tdt.NewCore()

	return le
}

// Internal helper methods -----------------------------------------------------

// check if string is in a slice
func stringIndexInSlice(a string, list []string) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}
