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
	ecspecFile = app.
			Flag("ecspecfile", "A CSV file contains reportURI and urn:epc:pat:<type>:<field1>.<field2>... .").
			Short('f').
			Default("ecspec.csv").
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
	statInterval = app.
			Flag("statInterval", "Measurement interval in seconds for the engine throughput.").
			Default("5").
			Int()
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
		sm = monitoring.NewStatManager("master", *influxAddr, *influxUser, *influxPass, *influxDB)
	}

	// load existing subscriptions from file
	log.Println("loading subscriptions from file")
	sub := filtering.LoadSubscriptionsFromCSVFile(*ecspecFile)

	// receive the engine instance status
	log.Println("setting up a management channel")
	mc := make(chan filtering.ManagementMessage, QueueSize)
	go func() {
		for {
			msg, ok := <-mc
			if !ok {
				break
			}
			switch msg.Type {
			case filtering.TrafficStatus: // TrafficStatus contains EventCount and MatchedCount
				if *enableStat {
					sm.StatMessageChannel <- monitoring.StatMessage{
						Type:  monitoring.Traffic,
						Value: []interface{}{msg.EventCount, msg.MatchedCount},
						Name:  msg.EngineName,
					}
				}
			case filtering.EngineStatus: // EngineStatus contains CurrentThroughput
				if *enableStat {
					sm.StatMessageChannel <- monitoring.StatMessage{
						Type:  monitoring.EngineThroughput,
						Value: []interface{}{msg.CurrentThroughput},
						Name:  msg.EngineName,
					}
				}
			case filtering.SelectedEngine: // SelectedEngine contains the EngineName just selected
				if *enableStat {
					sm.StatMessageChannel <- monitoring.StatMessage{
						Type: monitoring.SelectedEngine,
						Name: msg.EngineName,
					}
				}
			}
		}
		log.Fatalln("management channel closed, dying...")
	}()

	// set up an EngineFactory with a management channel
	log.Println("setting up an engine factory")
	engineFactory := filtering.NewEngineFactory(sub, *statInterval, mc)
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

	// receive incoming IDs and translate them in PureIdentity
	log.Println("setting up an incoming ReadEvent channel")
	var rq = make(chan []*llrp.ReadEvent)
	go func() {
		for {
			res, ok := <-rq
			if !ok {
				break
			}

			reports := map[string][]string{}
			for _, re := range res {
				pureIdentity, reportURIs, err := engineFactory.Search(*re)
				if err != nil { // no much or something went wrong
					continue
				}
				for _, dest := range reportURIs {
					if _, ok := reports[dest]; !ok {
						reports[dest] = []string{}
					}
					reports[dest] = append(reports[dest], pureIdentity)
				}
			}
			// do report
			for _, dest := range reports {
				_ = dest
			}
		}
		log.Fatalln("ReadEvent listener exited in gosstrak-fc")
	}()

	// establish a connection to the llrp client
	log.Println("waiting for the interrogator to becom online...")
	conn, err := net.Dial("tcp", *llrpAddr)
	for err != nil {
		time.Sleep(1 * time.Second)
		conn, err = net.Dial("tcp", *llrpAddr)
	}
	log.Printf("establised an LLRP connection to the interrogator %v", conn.RemoteAddr())

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
		mid := binary.BigEndian.Uint32(messageID)
		switch h {
		case llrp.ReaderEventNotificationHeader:
			log.Printf("[LLRP] %v >>> READER_EVENT_NOTIFICATION[%v]", conn.RemoteAddr(), mid)
			conn.Write(llrp.SetReaderConfig(currentMessageID))
		case llrp.KeepaliveHeader:
			log.Printf("[LLRP] %v >>> KEEP_ALIVE[%v]", conn.RemoteAddr(), mid)
			conn.Write(llrp.KeepaliveAck())
		case llrp.SetReaderConfigResponseHeader:
			log.Printf("[LLRP] %v >>> SET_READER_CONFIG_RESPONSE[%v]", conn.RemoteAddr(), mid)
		case llrp.ROAccessReportHeader:
			log.Printf("[LLRP] %v >>> RO_ACCESS_REPORT[%v]", conn.RemoteAddr(), mid)
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
