// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"log"
	"math"
	"time"
	//"reflect"

	"github.com/iomz/go-llrp"
	"github.com/looplab/fsm"
)

// EngineGenerator produce an engine according to the FSM
type EngineGenerator struct {
	FSM                 *fsm.FSM
	Name                string
	Engine              Engine
	managementChannel   chan ManagementMessage
	timePerEventChannel chan time.Duration
	totalTime           int64
	CurrentThroughput   float64
	EventCount          int64
	MatchedCount        int64
	statInterval        int
}

// NewEngineGenerator returns the pointer to a new EngineGenerator instance
func NewEngineGenerator(name string, ec EngineConstructor, statInterval int, mc chan ManagementMessage) *EngineGenerator {
	eg := &EngineGenerator{
		Name:              name,
		managementChannel: mc,
		totalTime:         0,
		CurrentThroughput: 0,
		EventCount:        0,
		MatchedCount:      0,
		statInterval:      statInterval,
	}

	eg.FSM = fsm.NewFSM(
		"unavailable",
		fsm.Events{
			{Name: "init", Src: []string{"unavailable"}, Dst: "generating"},
			{Name: "deploy", Src: []string{"generating", "rebuilding"}, Dst: "ready"},
			{Name: "update", Src: []string{"ready"}, Dst: "pending"},
			{Name: "rebuild", Src: []string{"pending"}, Dst: "rebuilding"},
		},
		fsm.Callbacks{
			"enter_state":      func(e *fsm.Event) { eg.enterState(e) },
			"enter_generating": func(e *fsm.Event) { eg.enterGenerating(e) },
			"enter_ready":      func(e *fsm.Event) { eg.enterReady(e) },
			"enter_pending":    func(e *fsm.Event) { eg.enterPending(e) },
			"enter_rebuilding": func(e *fsm.Event) { eg.enterRebuilding(e) },
		},
	)

	eg.timePerEventChannel = make(chan time.Duration)
	go func() {
		intervalTicker := time.NewTicker(time.Duration(eg.statInterval) * time.Second)

		for {
			select {
			case t, ok := <-eg.timePerEventChannel:
				if !ok {
					log.Fatalf("throughput monitor in EngingGenerator[%s] died", eg.Name)
				}
				//log.Printf("[EngineGenerator] %s: %v us/event", eg.Name, t.Nanoseconds())
				eg.totalTime += t.Nanoseconds() / 1000 // microseconds
				eg.EventCount++
			case <-intervalTicker.C:
				//log.Printf("%v, %v, %v", eg.Name, eg.EventCount, eg.MatchedCount)
				eg.managementChannel <- ManagementMessage{
					Type:         TrafficStatus,
					EngineName:   eg.Name,
					EventCount:   eg.EventCount,
					MatchedCount: eg.MatchedCount,
				}
				throughput := float64(eg.totalTime) / float64(eg.EventCount)
				if throughput != 0 && !math.IsNaN(throughput) {
					eg.CurrentThroughput = throughput
					eg.managementChannel <- ManagementMessage{
						Type:              EngineStatus,
						EngineName:        eg.Name,
						CurrentThroughput: eg.CurrentThroughput,
					}
				}
				eg.EventCount = 0
				eg.MatchedCount = 0
				eg.totalTime = 0
			}
		}
	}()

	return eg
}

// Search do search in the generated engine
func (eg *EngineGenerator) Search(re llrp.ReadEvent) (string, []string, error) {
	defer timeTrack(time.Now(), eg.timePerEventChannel)
	pureIdentity, reportURIs, err := eg.Engine.Search(re)
	if len(reportURIs) != 0 {
		eg.MatchedCount++
	}
	return pureIdentity, reportURIs, err
}

func (eg *EngineGenerator) enterState(e *fsm.Event) {
	log.Printf("[EngineGenerator] %s event, %s entering %s", e.Event, eg.Name, e.Dst)
}

func (eg *EngineGenerator) enterGenerating(e *fsm.Event) {
	go func() {
		//log.Printf("[EngineGenerator] start generating %s engine", eg.Name)
		sub := e.Args[0].(Subscriptions)
		eg.Engine = AvailableEngines[eg.Name](sub)
		eg.FSM.Event("deploy")
	}()
}

func (eg *EngineGenerator) enterRebuilding(e *fsm.Event) {
	msg := e.Args[0].(*ManagementMessage)
	switch msg.Type {
	case AddSubscription:
		/*
			eg.Engine.AddSubscription(ByteSubscriptions{
				msg.FilterString: &Info{
					Offset:          0,
					ReportURI: msg.ReportURI,
				},
			})
		*/
	case DeleteSubscription:
		/*
			eg.Engine.DeleteSubscription(ByteSubscriptions{
				msg.FilterString: &Info{
					Offset:          0,
					ReportURI: msg.ReportURI,
				},
			})
		*/
	}
	eg.FSM.Event("deploy")
}

func (eg *EngineGenerator) enterReady(e *fsm.Event) {
	log.Printf("[EngineGenerator] finished gererating %s engine", eg.Name)
	eg.managementChannel <- ManagementMessage{
		Type: OnEngineGenerated,
		EngineGeneratorInstance: eg,
	}
}

func (eg *EngineGenerator) enterPending(e *fsm.Event) {
	// Wait until the engine finishes the current execution
	eg.FSM.Event("rebuild", e.Args[0].(*ManagementMessage))
}
