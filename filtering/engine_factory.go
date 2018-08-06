// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/iomz/go-llrp"
)

// EngineFactory manages the FC's subscriptions and engine instances
type EngineFactory struct {
	mainChannel          chan ManagementMessage
	generatorChannels    []chan ManagementMessage
	currentSubscriptions Subscriptions
	productionSystem     map[string]*EngineGenerator
	deploymentPriority   map[string]uint8
	enginePerformance    sync.Map
	currentEngineName    string
	statInterval         int
}

// IsActive returns false if no engine is available
func (ef *EngineFactory) IsActive() bool {
	if len(ef.currentEngineName) == 0 {
		return false
	}
	return true
}

// Search is a wrapper for Search() with the current EngineGenerator
func (ef *EngineFactory) Search(re llrp.ReadEvent) (string, []string, error) {
	for name, eg := range ef.productionSystem {
		if name != ef.currentEngineName && eg.FSM.Is("ready") {
			_, _, _ = eg.Search(re)
		}
	}
	if !ef.IsActive() || !ef.productionSystem[ef.currentEngineName].FSM.Is("ready") {
		return "", []string{}, fmt.Errorf("%v is not ready", ef.currentEngineName)
	}
	return ef.productionSystem[ef.currentEngineName].Search(re)
}

// NewEngineFactory returns the pointer to a new EngineFactory instance
func NewEngineFactory(sub Subscriptions, statInterval int, mc chan ManagementMessage) *EngineFactory {
	ef := &EngineFactory{
		mainChannel:  mc,
		statInterval: statInterval,
	}

	// Load saved subscriptions?
	ef.currentSubscriptions = sub

	// Load all the possible engines
	ef.productionSystem = make(map[string]*EngineGenerator)
	ef.enginePerformance = sync.Map{}
	ef.generatorChannels = []chan ManagementMessage{}
	for name, constructor := range AvailableEngines {
		ch := make(chan ManagementMessage)
		ef.generatorChannels = append(ef.generatorChannels, ch)
		eg := NewEngineGenerator(name, constructor, statInterval, ch)
		ef.productionSystem[name] = eg
		ef.enginePerformance.Store(name, float64(0))
	}

	// Calculate the priority of deployment
	ef.deploymentPriority = map[string]uint8{
		"LegacyEngine": 0,
		"List":         1,
		"PatriciaTrie": 3,
		"SplayTree":    2,
	}

	log.Printf("[EngineFactory] deploymentPriority: %v", ef.deploymentPriority)

	return ef
}

// Run starts the engine factory to react with the ManagementChannel
func (ef *EngineFactory) Run() {
	log.Println("[EngineFactory] start running")
	// set channels from EngineGenerators + main
	cases := make([]reflect.SelectCase, len(ef.generatorChannels)+1)
	for i, ch := range append(ef.generatorChannels, ef.mainChannel) {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	go func() { // simulation stat ticker
		log.Println("[EngineFactory] setting up selective adoption handler")
		intervalTicker := time.NewTicker(time.Duration(ef.statInterval) * time.Second)
		for {
			select {
			case <-intervalTicker.C:
				/*
					var ename string
					perf := float64(0)
					ef.enginePerformance.Range(func(k, v interface{}) bool {
						if reflect.ValueOf(v).Float() > perf {
							ename = reflect.ValueOf(k).String()
							perf = reflect.ValueOf(v).Float()
						}
						return true
					})
					if ef.currentEngineName != ename && len(ename) != 0 {
						log.Printf("[EngineFactory] %s replaces the currentEngine %s due to performance", ename, ef.currentEngineName)
						ef.currentEngineName = ename
					}
					ef.mainChannel <- ManagementMessage{
						Type:       SelectedEngine,
						EngineName: ef.currentEngineName,
					}
				*/
				v, ok := ef.enginePerformance.Load(ef.currentEngineName)
				if !ok {
					v = reflect.ValueOf(float64(0)).Interface()
				}
				ef.mainChannel <- ManagementMessage{
					Type:              SimulationStat,
					EngineName:        ef.currentEngineName,
					CurrentThroughput: reflect.ValueOf(v).Float(),
				}
				endFlag := true
				for _, eg := range ef.productionSystem {
					if !eg.FSM.Is("ready") {
						endFlag = false
					}
				}
				if endFlag {
					log.Print("\n\n\n\n\nall the engines ready\n\n\n\n\n")
				}
			}
		}
	}()

	go func() {
		log.Println("[EngineFactory] setting up managementChannel listener")

		for {
			_, val, ok := reflect.Select(cases)
			if !ok {
				break
			}
			//msg, _ := reflect.ValueOf(val).Interface().(ManagementMessage)
			msg := ManagementMessage{
				Type:                    val.FieldByName("Type").Interface().(ManagementMessageType),
				Pattern:                 val.FieldByName("Pattern").String(),
				ReportURI:               val.FieldByName("ReportURI").String(),
				EngineGeneratorInstance: (*EngineGenerator)(unsafe.Pointer(val.FieldByName("EngineGeneratorInstance").Pointer())),
				CurrentThroughput:       val.FieldByName("CurrentThroughput").Float(),
				EventCount:              val.FieldByName("EventCount").Int(),
				MatchedCount:            val.FieldByName("MatchedCount").Int(),
				EngineName:              val.FieldByName("EngineName").String(),
			}
			switch msg.Type {
			case AddSubscription:
				/*
					if _, ok := ef.currentByteSubscriptions[msg.FilterString]; !ok {
						ef.currentByteSubscriptions[msg.FilterString] = &PartialSubscription{
							Offset:    0,
							ReportURI: msg.ReportURI,
						}
						for _, eg := range ef.productionSystem {
							err := eg.FSM.Event("update", &msg)
							if err != nil {
								log.Println(err)
							}
						}
					}
				*/
			case DeleteSubscription:
				/*
					if _, ok := ef.currentByteSubscriptions[msg.FilterString]; ok {
						delete(ef.currentByteSubscriptions, msg.FilterString)
						for _, eg := range ef.productionSystem {
							err := eg.FSM.Event("update", &msg)
							if err != nil {
								log.Println(err)
							}
						}
					}
				*/
			case OnEngineGenerated:
				log.Printf("[EngineFactory] received OnEngineGenerated from %s", msg.EngineGeneratorInstance.Engine.Name())
				if len(ef.currentEngineName) == 0 {
					log.Printf("[EngineFactory] set %s as an initial engine", msg.EngineGeneratorInstance.Name)
					ef.currentEngineName = msg.EngineGeneratorInstance.Name
					ef.mainChannel <- ManagementMessage{
						Type:       SelectedEngine,
						EngineName: ef.currentEngineName,
					}
					continue
				}
				if ef.deploymentPriority[ef.currentEngineName] < ef.deploymentPriority[msg.EngineGeneratorInstance.Name] {
					log.Printf("[EngineFactory] %s replaces the currentEngine %s", msg.EngineGeneratorInstance.Name, ef.currentEngineName)
					ef.currentEngineName = msg.EngineGeneratorInstance.Name
					ef.mainChannel <- ManagementMessage{
						Type:       SelectedEngine,
						EngineName: ef.currentEngineName,
					}
					continue
				}
				log.Printf("[EngineFactory] %s didn't replace the currentEngine %s", msg.EngineGeneratorInstance.Name, ef.currentEngineName)
			case TrafficStatus:
				ef.mainChannel <- msg // bypass the status message from generators to main
			case EngineStatus:
				ef.enginePerformance.Store(msg.EngineName, msg.CurrentThroughput)
				ef.mainChannel <- msg // bypass the status message from generators to main
			}
		}
		log.Fatalln("mainChannel listener exited in gosstrak-fc")
	}()

	// initialize the engines
	log.Println("[EngineFactory] initializing engines")
	for _, eg := range ef.productionSystem {
		// pass the cloned subscriptions
		eg.FSM.Event("init", ef.currentSubscriptions)
	}
}
