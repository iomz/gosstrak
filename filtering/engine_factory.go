// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"log"
	"reflect"
	"unsafe"

	"github.com/iomz/go-llrp"
)

// ManagementMessageType is to indicate the type of ManagementMessage
type ManagementMessageType int

// ManagementMessage types
const (
	AddSubscription ManagementMessageType = iota
	DeleteSubscription
	OnEngineGenerated
	DeployEngine
	EngineStatus
)

// ManagementMessage holds management action for the EngineFactory
type ManagementMessage struct {
	Type                    ManagementMessageType
	FilterString            string
	ReportURI               string
	EngineGeneratorInstance *EngineGenerator
	CurrentThroughput       float64
}

// EngineFactory manages the FC's subscriptions and engine instances
type EngineFactory struct {
	mainChannel          chan ManagementMessage
	generatorChannels    []chan ManagementMessage
	currentSubscriptions Subscriptions
	productionSystem     map[string]*EngineGenerator
	deploymentPriority   map[string]uint8
	currentEngineName    string
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
		if name != ef.currentEngineName {
			_, _, _ = eg.Search(re)
		}
	}
	return ef.productionSystem[ef.currentEngineName].Search(re)
}

// NewEngineFactory returns the pointer to a new EngineFactory instance
func NewEngineFactory(sub Subscriptions, statInterval int, mc chan ManagementMessage) *EngineFactory {
	ef := &EngineFactory{
		mainChannel: mc,
	}

	// Load saved subscriptions?
	ef.currentSubscriptions = sub

	// Load all the possible engines
	ef.productionSystem = make(map[string]*EngineGenerator)
	ef.generatorChannels = []chan ManagementMessage{}
	for name, constructor := range AvailableEngines {
		ch := make(chan ManagementMessage)
		ef.generatorChannels = append(ef.generatorChannels, ch)
		eg := NewEngineGenerator(name, constructor, statInterval, ch)
		ef.productionSystem[name] = eg
	}

	// Calculate the priority of deployment
	ef.deploymentPriority = map[string]uint8{
		"List":         3,
		"PatriciaTrie": 1,
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
				FilterString:            val.FieldByName("FilterString").String(),
				ReportURI:               val.FieldByName("FilterString").String(),
				EngineGeneratorInstance: (*EngineGenerator)(unsafe.Pointer(val.FieldByName("EngineGeneratorInstance").Pointer())),
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
					continue
				}
				if ef.deploymentPriority[ef.currentEngineName] < ef.deploymentPriority[msg.EngineGeneratorInstance.Name] {
					log.Printf("[EngineFactory] %s replaces the currentEngine %s", msg.EngineGeneratorInstance.Name, ef.currentEngineName)
					ef.currentEngineName = msg.EngineGeneratorInstance.Name
					continue
				}
				log.Printf("[EngineFactory] %s didn't replace the currentEngine %s", msg.EngineGeneratorInstance.Name, ef.currentEngineName)
			case EngineStatus:
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
