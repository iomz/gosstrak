package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Current Version
	version = "0.1.0"

	// kingpin app
	app = kingpin.New("gosstrak-fc", "An RFID middleware to replace Fosstrak F&C.")
	// kingpin verbose mode flag
	verbose = app.Flag("debug", "Enable verbose mode.").Short('v').Default("false").Bool()

	// kingpin server command
	server = app.Command("run", "Run as an F&C middleware.")
)

func runMW() int {
	fmt.Println(">> running the gosstrak-fc")
	return 0
}

func main() {
	app.Version(version)
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parse {
	case server.FullCommand():
		os.Exit(runMW())
	}
}
