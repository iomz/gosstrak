// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/docker/libchan/spdy"
	"github.com/iomz/go-llrp"
	"github.com/iomz/gosstrak/filtering"
	"github.com/iomz/gosstrak/monitoring"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Notification is the struct to send/receive captured ID
type Notification struct {
	ID []byte
}

// Constant Values
const (
	// BufferSize is a general size for a buffer
	BufferSize = 64 * 1024 // 64 KiB
	QueueSize  = 128
)

// Environmental variables
var (
	// Current Version
	version = "0.3.0"

	// app
	app = kingpin.
		New("gosstrak-fc", "An RFID middleware to replace Fosstrak F&C.")

	// common flag
	verbose = app.
		Flag("debug", "Enable verbose mode.").
		Short('v').
		Default("false").
		Bool()
	filterFile = app.
			Flag("filter-file", "A CSV file contains filter and notify.").
			Short('f').
			Default("filters.csv").
			String()

	// LLRP related values
	llrpInitialMessageID = app.
				Flag("initialMessageID", "The initial messageID to start from.").
				Short('m').
				Default("1000").
				Int()
	llrpAddr = app.
			Flag("ip", "LLRP emulator address.").
			Short('l').
			Default("127.0.0.1:5084").
			String()

	// ALE related values
	managementAddr = app.
			Flag("managementAddr", "Psuedo ALE management endpoint").
			Default("127.0.0.1:2784").
			String()

	// stat related values
	enableStat = app.
			Flag("enableStat", "Enable statistical monitoring.").
			Default("false").
			Bool()
	influxAddr = app.
			Flag("influxAddr", "The endpoint of influxdb.").
			Default("http://127.0.0.1:8086").
			String()
	influxUser = app.
			Flag("influxUser", "The username for influxdb.").
			Default("gosstrak").
			String()
	influxPass = app.
			Flag("influxPass", "The password for influxdb.").
			Default("gosstrak").
			String()
	influxDB = app.
			Flag("influxDB", "The database in influxdb.").
			Default("gosstrak").
			String()

	// start command
	cmdStart = app.Command("start", "Start the gosstrak-fc.")

	// Current messageID
	currentMessageID = uint32(*llrpInitialMessageID)
)

func getPackagePath() string {
	// Determine the package dir
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return path.Dir(filename)
}

func run() {
	log.Println("initializing gosstrak-fc for master mode...")

	// setup StatManager
	var sm *monitoring.StatManager
	if *enableStat {
		log.Println("setting up a stat manager for InfluxDB")
		sm = monitoring.NewStatManager(*influxAddr, *influxUser, *influxPass, *influxDB)
	}

	// load existing subscriptions from file
	log.Println("loading subscriptions from file")
	sub := filtering.LoadFiltersFromCSVFile(*filterFile)

	// receive the engine instance
	log.Println("setting up a management channel")
	mc := make(chan filtering.ManagementMessage, QueueSize)
	go func() {
		for {
			msg, ok := <-mc
			if !ok {
				break
			}
			switch msg.Type {
			case filtering.EngineStatus:
				if *enableStat {
					sm.StatMessageChannel <- monitoring.StatMessage{
						Type:  monitoring.EngineThroughput,
						Value: []interface{}{msg.EngineGeneratorInstance.CurrentThroughput},
						Name:  msg.EngineGeneratorInstance.Name,
					}
				}
			}
		}
		log.Fatalln("management channel closed, dying...")
	}()

	// set up an EngineFactory with a management channel
	log.Println("setting up an engine factory")
	engineFactory := filtering.NewEngineFactory(sub, mc)
	go engineFactory.Run()
	// wait until the first engine becomes available
	for !engineFactory.IsActive() {
		time.Sleep(time.Second)
	}

	// receive management access
	log.Println("setting up an management interface")
	go func() {
		managementListener, err := net.Listen("tcp", *managementAddr)
		if err != nil {
			log.Fatal(err)
		}
		for {
			c, err := managementListener.Accept()
			if err != nil {
				log.Fatal(err)
				break
			}
			p, err := spdy.NewSpdyStreamProvider(c, true)
			if err != nil {
				log.Print(err)
				continue
			}
			t := spdy.NewTransport(p)

			receiver, err := t.WaitReceiveChannel()
			if err != nil {
				log.Print(err)
				continue
			}

			mm := &filtering.ManagementMessage{}
			err = receiver.Receive(mm)
			if err != nil {
				log.Print(err)
				continue
			}
			log.Print(mm)
			mc <- *mm
		}
		log.Fatalln("managementListener closed in gosstrak-fc")
	}()

	// receive incoming IDs
	log.Println("setting up an incoming LLRPReadEvent channel")
	var rq = make(chan []*llrp.LLRPReadEvent)
	go func() {
		for {
			res, ok := <-rq
			if !ok {
				break
			}

			//log.Printf("%v tags received", len(res))
			matches := map[string][]byte{}
			for _, re := range res {
				dests := engineFactory.Search(re.ID)
				for _, d := range dests {
					matches[d] = re.ID
				}
			}
			if *enableStat {
				sm.StatMessageChannel <- monitoring.StatMessage{
					Type:  monitoring.Traffic,
					Value: []interface{}{len(res), len(matches)},
				}
			}
			//log.Printf("%v matches", len(matches))
			//notify matches
			/*
				for _, m := range matches {
					log.Printf("match: %s <- %v,%v\n", m, re.PC, re.ID)
				}
			*/
		}
		log.Fatalln("LLRPReadEvent listener exited in gosstrak-fc")
	}()

	// establish a connection to the llrp client
	log.Println("starting a connection to an interrogator")
	conn, err := net.Dial("tcp", *llrpAddr)
	for err != nil {
		log.Print(err)
		log.Println("wait 5 seconds for the interrogator to becom online...")
		time.Sleep(5 * time.Second)
		conn, err = net.Dial("tcp", *llrpAddr)
	}

	// prepare LLRP header storage
	header := make([]byte, 2)
	length := make([]byte, 4)
	messageID := make([]byte, 4)
	for {
		_, err = io.ReadFull(conn, header)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.ReadFull(conn, length)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.ReadFull(conn, messageID)
		if err != nil {
			log.Fatal(err)
		}
		// length containts the size of the entire message in octets
		// starting from bit offset 0, hence, the message size is
		// length - 10 bytes
		var messageValue []byte
		if messageSize := binary.BigEndian.Uint32(length) - 10; messageSize != 0 {
			messageValue = make([]byte, binary.BigEndian.Uint32(length)-10)
			_, err = io.ReadFull(conn, messageValue)
			if err != nil {
				log.Fatal(err)
			}
		}

		h := binary.BigEndian.Uint16(header)
		switch h {
		case llrp.ReaderEventNotificationHeader:
			log.Println(">>> READER_EVENT_NOTIFICATION")
			conn.Write(llrp.SetReaderConfig(currentMessageID))
		case llrp.KeepaliveHeader:
			log.Println(">>> KEEP_ALIVE")
			conn.Write(llrp.KeepaliveAck())
		case llrp.SetReaderConfigResponseHeader:
			log.Println(">>> SET_READER_CONFIG_RESPONSE")
		case llrp.ROAccessReportHeader:
			log.Println(">>> RO_ACCESS_REPORT")
			rq <- llrp.UnmarshalROAccessReportBody(messageValue)
		default:
			log.Fatalf("Unknown LLRP Message Header: %v\n", h)
		}
	}
}

func main() {
	app.Version(version)
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Create cache directory if not exists
	// TODO: set OS specific dataCacheDir
	dataCacheDir := "/var/tmp/gosstrak-fc-cache"
	if _, err := os.Stat(dataCacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(dataCacheDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	switch parse {
	case cmdStart.FullCommand():
		run()
	}
}
