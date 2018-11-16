// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package monitoring

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

// StatManager receives stat and publish them to InfluxDB
type StatManager struct {
	StatMessageChannel chan StatMessage
}

// NewStatManager creates a new instance of StatManager
func NewStatManager(mode string, addr string, user string, pass string, db string) *StatManager {
	// create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: pass,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "ns",
	})
	if err != nil {
		log.Fatal(err)
	}

	// make the stat message channel
	smc := make(chan StatMessage)

	go func() {
		for {
			msg, ok := <-smc
			if !ok {
				break
			}

			tags := make(map[string]string)
			fields := make(map[string]interface{})
			var measurement string

			// create a point
			switch msg.Type {
			case Traffic:
				ingress, ok := msg.Value[0].(int64)
				if !ok {
					continue
				}
				fields["incoming_events"] = ingress
				matches, ok := msg.Value[1].(int64)
				if !ok {
					continue
				}
				fields["matched_events"] = matches
				if ingress != 0 {
					mp := float64(matches) / float64(ingress) * 100.0
					if mp > 100 {
						mp = 100
					}
					fields["matching_probability"] = mp
				}
				tags["engine"] = msg.Name
				measurement = "traffic"
				//log.Printf("==> Traffic (incoming, matched) = %v, %v", ingress, matches)
			case EngineThroughput:
				fields["event_per_us"] = msg.Value[0]
				tags["engine"] = msg.Name
				measurement = "throughput"
				log.Printf("==> Throughput [%v]: %v", msg.Name, msg.Value[0])
			case SelectedEngine:
				engineType := 0
				switch msg.Name {
				case "Legacy":
					engineType = 0
				case "List":
					engineType = 1
				case "PatriciaTrie":
					engineType = 2
				case "SplayTree":
					engineType = 3
				case "DummyEngine":
					engineType = 4
				}
				fields["selected"] = engineType
				measurement = "engine"
				log.Printf("==> Engine: %v", msg.Name)
			}
			pt, err := client.NewPoint(measurement, tags, fields, time.Now())
			if err != nil {
				log.Fatal(err)
			}
			bp.AddPoint(pt)

			// write the batch
			if err := c.Write(bp); err != nil {
				log.Fatal(err)
			}
		}
		// close client resources
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
		log.Fatalln("StatMessageChannel closed, dying...")
	}()

	return &StatManager{StatMessageChannel: smc}
}
