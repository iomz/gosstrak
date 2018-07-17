package filtering

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/iomz/go-llrp"
	"github.com/iomz/gosstrak/tdt"
)

type LegacyEngine struct {
	mainChannel         chan ManagementMessage
	subscriptions       map[string][]string
	tdtCore             *tdt.Core
	statInterval        int
	nEvent              int
	totalTime           int64
	timePerEventChannel chan time.Duration
	CurrentThroughput   float64
}

func (le *LegacyEngine) Search(re *llrp.ReadEvent) (matched []string, pureIdentity string, err error) {
	defer timeTrack(time.Now(), le.timePerEventChannel)
	// translate the readevent to a PureIdentity
	pureIdentity, err = le.tdtCore.Translate(re.PC, re.ID)
	if err != nil {
		return matched, pureIdentity, err
	}

	for prefix, dests := range le.subscriptions {
		if strings.HasPrefix(pureIdentity, prefix) {
			matched = append(matched, dests...)
		}
	}
	if len(matched) == 0 {
		return matched, pureIdentity, fmt.Errorf("no matching subscription for %v", pureIdentity)
	}
	return matched, pureIdentity, nil
}

func NewLegacyEngine(filename string, statInterval int, mc chan ManagementMessage) *LegacyEngine {
	// initialize LegacyEngine
	le := LegacyEngine{
		mainChannel:       mc,
		statInterval:      statInterval,
		nEvent:            0,
		totalTime:         0,
		CurrentThroughput: 0,
	}
	le.subscriptions = make(map[string][]string)

	// initialize tdt.Core
	le.tdtCore = tdt.NewCore()

	// read the file and load them up
	fp, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
		if len(record) == 2 {
			// convert the filter label to urn prefix
			urnPrefix, err := convertLabelToURNPrefix(record[0])
			if err != nil {
				continue
			}
			if _, ok := le.subscriptions[urnPrefix]; !ok {
				le.subscriptions[urnPrefix] = []string{record[0]}
			} else {
				le.subscriptions[urnPrefix] = append(le.subscriptions[urnPrefix], record[0])
			}
		}
	}
	log.Printf("%v filters loaded from %v", len(le.subscriptions), filename)

	le.timePerEventChannel = make(chan time.Duration)
	go func() {
		intervalTicker := time.NewTicker(time.Duration(le.statInterval) * time.Second)

		for {
			select {
			case t, ok := <-le.timePerEventChannel:
				if !ok {
					log.Fatalln("throughput monitor in LegacyEngine died")
				}
				//log.Printf("[EngineGenerator] %s: %v us/event", eg.Name, t.Nanoseconds())
				le.totalTime += t.Nanoseconds() / 1000 // microseconds
				le.nEvent++
			case <-intervalTicker.C:
				throughput := float64(le.totalTime) / float64(le.nEvent)
				//log.Printf("%s total: %v, n: %v", eg.Name, eg.totalTime, eg.nEvent)
				if throughput != 0 && !math.IsNaN(throughput) {
					le.CurrentThroughput = throughput
					le.mainChannel <- ManagementMessage{
						Type:              EngineStatus,
						CurrentThroughput: le.CurrentThroughput,
					}
				}
				le.nEvent = 0
				le.totalTime = 0
			}
		}
	}()

	return &le
}

func convertLabelToURNPrefix(label string) (urnPrefix string, err error) {
	elems := strings.Split(label, "_")
	if len(elems) < 1 {
		return urnPrefix, fmt.Errorf("%v cannot be parsed to an URN prefix", label)
	}
	switch elems[0] {
	case "GIAI-96":
		urnPrefix = "urn:epc:id:giai:"
		elems = append(elems[:2], elems[3:]...)
	case "GRAI-96":
		urnPrefix = "urn:epc:id:grai:"
		elems = append(elems[:2], elems[3:]...)
	case "SGTIN-96":
		urnPrefix = "urn:epc:id:sgtin:"
		elems = append(elems[:2], elems[3:]...)
	case "SSCC-96":
		urnPrefix = "urn:epc:id:sscc:"
		elems = append(elems[:2], elems[3:]...)
	case "ISO17363":
		urnPrefix = "urn:epc:id:iso:17363:"
	case "ISO17365":
		urnPrefix = "urn:epc:id:iso:17365:"
	}
	for i := 1; i < len(elems); i++ {
		if i != 1 {
			urnPrefix += "."
		}
		urnPrefix += elems[i]
	}

	return urnPrefix, nil
}
