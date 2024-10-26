// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"log"
	"net"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/moby/spdystream"
)

// Notification is the struct to send/receive captured ID
type Notification struct {
	ID []byte
}

var (
	// kingpin app
	app  = kingpin.New("noticap", "Capture notification from F&C.")
	port = app.Flag("port", "Listening port.").Short('p').Default("9323").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.Println("Starting event capture...")
	var listener net.Listener
	var err error
	listener, err = net.Listen("tcp", "0.0.0.0:"+*port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening at 0.0.0.0:" + *port)

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Print(err)
			break
		}

		spdyConn, err := spdystream.NewConnection(c, true)
		if err != nil {
			log.Print(err)
			continue
		}
		go spdyConn.Serve(spdystream.MirrorStreamHandler)

		// go func() {
		// 	for {
		// 		receiver, err := t.WaitReceiveChannel()
		// 		if err != nil {
		// 			log.Print(err)
		// 			break
		// 		}

		// 		go func() {
		// 			for {
		// 				noti := &Notification{}
		// 				err := receiver.Receive(noti)
		// 				if err != nil {
		// 					if err != io.EOF {
		// 						log.Print(err)
		// 					}
		// 					break
		// 				}
		// 				log.Println(noti.ID)
		// 			}
		// 		}()
		// 	}
		// }()
	}
}
