package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/icedream/icecon/internal/rcon"
	"github.com/icedream/icecon/internal/ui"

	_ "github.com/icedream/icecon/internal/ui/console"
)

var (
	flagCommand = kingpin.Flag("command",
		"Run a one-off command and then exit.").
		Short('c').String()
	argAddress = kingpin.Arg("address",
		"Server IP/hostname and port, written as \"server:port\".")
	argPassword = kingpin.Arg("password", "The RCON password.")

	password string
)

func usage() {
	kingpin.Usage()
}

var (
	hasGraphicalUI = ui.HasGraphicalUI()
	flagGui        *bool
)

func init() {
	// only provide -gui/-g flag if there is a graphical user interface available
	if hasGraphicalUI {
		flagGui = kingpin.
			Flag("gui", "Run as GUI (runs automatically as GUI if no arguments given, ignored if command flag used)").
			Short('g').Bool()
	}
}

func main() {
	fmt.Println("IceCon - Icedream's RCON Client")
	fmt.Println("\t\u00A9 2016-2025 Carl Kittelberger/Icedream")
	fmt.Println()

	argAddressTCP := argAddress.TCP()
	argPasswordStr := argPassword.String()

	kingpin.Parse()

	// If no arguments, fall back to running the shell
	wantGui := (*argAddressTCP == nil && *flagCommand == "") || *flagGui

	// Command-line shell doesn't support starting up without arguments
	// but graphical Windows UI does
	if !(hasGraphicalUI && wantGui) {
		argAddress = argAddress.Required()
		argPassword = argPassword.Required()
		kingpin.Parse()
	}

	// Initialize socket
	rconClient := rcon.NewRconClient()
	rconClient.InitSocket()
	defer rconClient.Release()

	// Set target address if given
	if *argAddressTCP != nil {
		rconClient.SetSocketAddr((*argAddressTCP).String())
	}

	// Get password
	password = *argPasswordStr

	// Run one-off command?
	if *flagCommand != "" {
		// Send
		err := rconClient.Send(*flagCommand)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Receive
		msg, err := rconClient.Receive()
		if err != nil {
			log.Fatal(err)
			return
		}
		switch strings.ToLower(msg.Name) {
		case "print":
			fmt.Println(string(msg.Data))
		}
		return
	}

	// Which UI should be run?
	if err := ui.Run(rconClient, wantGui); err != nil {
		log.Fatal("User interface failed:", err)
		return
	}
}
