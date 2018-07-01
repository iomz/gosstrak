// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"log"
	"reflect"
	"unsafe"
)

// ManagementMessageType is to indicate the type of ManagementMessage
type ManagementMessageType int

// ManagementMessage types
const (
	AddSubscription ManagementMessageType = iota
	DeleteSubscription
	OnEngineGenerated
	DeployEngine
)

// ManagementMessage holds management action for the EngineFactory
type ManagementMessage struct {
	Type                    ManagementMessageType
	FilterString            string
	NotificationURI         string
	EngineGeneratorInstance *EngineGenerator
}

// EngineFactory manages the FC's subscriptions and engine instances
type EngineFactory struct {
	mainChannel          chan ManagementMessage
	generatorChannels    []chan ManagementMessage
	currentSubscriptions Subscriptions
	productionLines      []*EngineGenerator
	deploymentPriority   map[string]uint8
	currentEngineName    string
}

// NewEngineFactory returns the pointer to a new EngineFactory instance
func NewEngineFactory(sub Subscriptions, mc chan ManagementMessage) *EngineFactory {
	ef := &EngineFactory{
		mainChannel: mc,
	}

	// Load saved subscriptions?
	ef.currentSubscriptions = sub

	// Load all the possible engines
	ef.productionLines = []*EngineGenerator{}
	ef.generatorChannels = []chan ManagementMessage{}
	for name, constructor := range AvailableEngines {
		ch := make(chan ManagementMessage)
		ef.generatorChannels = append(ef.generatorChannels, ch)
		eg := NewEngineGenerator(name, constructor, ch)
		ef.productionLines = append(ef.productionLines, eg)
	}

	// Calculate the priority of deployment
	ef.deploymentPriority = map[string]uint8{}
	priority := uint8(0)
	for name, _ := range AvailableEngines {
		ef.deploymentPriority[name] = priority
		priority++
	}

	log.Print("[EngineFactory] deploymentPriority: %v", ef.deploymentPriority)

	return ef
}

// Run starts the engine factory to react with the ManagementChannel
func (ef *EngineFactory) Run() {
	// set channels from EngineGenerators + main
	cases := make([]reflect.SelectCase, len(ef.generatorChannels)+1)
	for i, ch := range append(ef.generatorChannels, ef.mainChannel) {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	go func() {
		for {
			_, val, ok := reflect.Select(cases)
			if !ok {
				break
			}
			//msg, _ := reflect.ValueOf(val).Interface().(ManagementMessage)
			msg := ManagementMessage{
				Type:                    val.FieldByName("Type").Interface().(ManagementMessageType),
				FilterString:            val.FieldByName("FilterString").String(),
				NotificationURI:         val.FieldByName("FilterString").String(),
				EngineGeneratorInstance: (*EngineGenerator)(unsafe.Pointer(val.FieldByName("EngineGeneratorInstance").Pointer())),
			}
			switch msg.Type {
			case AddSubscription:
				if _, ok := ef.currentSubscriptions[msg.FilterString]; !ok {
					ef.currentSubscriptions[msg.FilterString] = &Info{
						Offset:          0,
						NotificationURI: msg.NotificationURI,
					}
					for _, eg := range ef.productionLines {
						err := eg.FSM.Event("update", &msg)
						if err != nil {
							log.Println(err)
						}
					}
				}
			case DeleteSubscription:
				if _, ok := ef.currentSubscriptions[msg.FilterString]; ok {
					delete(ef.currentSubscriptions, msg.FilterString)
					for _, eg := range ef.productionLines {
						err := eg.FSM.Event("update", &msg)
						if err != nil {
							log.Println(err)
						}
					}
				}
			case OnEngineGenerated:
				log.Printf("[EngineFactory] received OnEngineGenerated from %s\n", msg.EngineGeneratorInstance.Name)
				if len(ef.currentEngineName) == 0 {
					ef.currentEngineName = msg.EngineGeneratorInstance.Name
					ef.mainChannel <- ManagementMessage{
						Type: DeployEngine,
						EngineGeneratorInstance: msg.EngineGeneratorInstance,
					}
					continue
				}
				if ef.deploymentPriority[ef.currentEngineName] < ef.deploymentPriority[msg.EngineGeneratorInstance.Name] {
					ef.currentEngineName = msg.EngineGeneratorInstance.Name
					ef.mainChannel <- ManagementMessage{
						Type: DeployEngine,
						EngineGeneratorInstance: msg.EngineGeneratorInstance,
					}
					continue
				}
				log.Printf("[EngineFactory] %s didn't replace the currentEngine %s", msg.EngineGeneratorInstance.Name, ef.currentEngineName)
			}
		}
		log.Fatalln("mainChannel listener exited in gosstrak-fc")
	}()

	// initialize the engines
	for _, eg := range ef.productionLines {
		// pass the cloned subscriptions
		eg.FSM.Event("init", ef.currentSubscriptions.Clone())
	}
}
