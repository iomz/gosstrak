// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"log"
	//"reflect"

	"github.com/looplab/fsm"
)

// EngineGenerator produce an engine according to the FSM
type EngineGenerator struct {
	managementChannel chan ManagementMessage
	Name              string
	Engine            Engine
	FSM               *fsm.FSM
}

// NewEngineGenerator returns the pointer to a new EngineGenerator instance
func NewEngineGenerator(name string, ec EngineConstructor, mc chan ManagementMessage) *EngineGenerator {
	eg := &EngineGenerator{
		managementChannel: mc,
		Name:              name,
	}

	eg.FSM = fsm.NewFSM(
		"unavailable",
		fsm.Events{
			{Name: "init", Src: []string{"unavailable"}, Dst: "generating"},
			{Name: "generated", Src: []string{"generating", "pending"}, Dst: "ready"},
			{Name: "deploy", Src: []string{"ready"}, Dst: "deployed"},
			{Name: "update", Src: []string{"generating", "ready", "deployed"}, Dst: "pending"},
		},
		fsm.Callbacks{
			"enter_state":      func(e *fsm.Event) { eg.enterState(e) },
			"enter_generating": func(e *fsm.Event) { eg.enterGenerating(e) },
			"enter_ready":      func(e *fsm.Event) { eg.enterReady(e) },
			//"enter_deployed":   func(e *fsm.Event) { eg.enterDeployed(e) },
			"enter_pending": func(e *fsm.Event) { eg.enterPending(e) },
		},
	)

	return eg
}

func (eg *EngineGenerator) enterState(e *fsm.Event) {
	log.Printf("[EngineGenerator] %s, entering %s\n", e.Event, e.Dst)
}

func (eg *EngineGenerator) enterGenerating(e *fsm.Event) {
	go func() {
		sub := e.Args[0].(Subscriptions)
		eg.Engine = AvailableEngines[eg.Name](sub)
		eg.FSM.Event("generated")
	}()
}

func (eg *EngineGenerator) enterReady(e *fsm.Event) {
	log.Printf("[EngineGenerator] finished gererating %s engine\n", eg.Name)
	eg.managementChannel <- ManagementMessage{
		Type: OnEngineGenerated,
		EngineGeneratorInstance: eg,
	}
	err := eg.FSM.Event("deploy")
	if err != nil {
		log.Fatal(err)
	}
}

func (eg *EngineGenerator) enterDeployed(e *fsm.Event) {
	// do nothing
}

func (eg *EngineGenerator) enterPending(e *fsm.Event) {
	msg := e.Args[0].(*ManagementMessage)
	switch msg.Type {
	case AddSubscription:
		eg.Engine.AddSubscription(Subscriptions{
			m: SubMap{
				msg.FilterString: &Info{
					Offset:          0,
					NotificationURI: msg.NotificationURI,
				},
			},
		})
	case DeleteSubscription:
		eg.Engine.DeleteSubscription(Subscriptions{
			m: SubMap{
				msg.FilterString: &Info{
					Offset:          0,
					NotificationURI: msg.NotificationURI,
				},
			},
		})
	}
	eg.FSM.Event("generated")
}
