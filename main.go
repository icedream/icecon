package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/icedream/go-q3net"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	flagCommand = kingpin.Flag("command",
		"Run a one-off command and then exit.").
		Short('c').String()
	argAddress = kingpin.Arg("address",
		"Server IP/hostname and port, written as \"server:port\".")
	argPassword = kingpin.Arg("password", "The RCON password.")

	address    *net.UDPAddr
	addressStr string
	password   string

	socket       *net.UDPConn
	socketBuffer = make([]byte, 64*1024)

	bufferedStdin *bufio.Reader

	errNotSupported = errors.New("Not supported")
)

func initSocketAddr(addr string) (err error) {
	newAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}

	address, addressStr = newAddr, addr

	return
}

func initSocket() (err error) {
	socket, err = net.ListenUDP("udp", nil)
	if err != nil {
		return
	}

	return
}

func receive() (msg *quake.Message, err error) {
	length, _, err := socket.ReadFromUDP(socketBuffer)
	if err != nil {
		return
	}

	msg, err = quake.UnmarshalMessage(socketBuffer[0:length])
	if err != nil {
		return
	}
	return
}

func receiveRcon() (msg *quake.Message, err error) {
	msg, err = receive()
	if err != nil {
		return
	}
	if !strings.EqualFold(msg.Name, "print") {
		err = errors.New("rcon: Unexpected response from server: " + msg.Name)
	}
	return
}

func sendRcon(input string) (err error) {
	buf := new(bytes.Buffer)
	msg := &quake.Message{
		Header: quake.OOBHeader,
		Name:   "rcon",
		Data:   []byte(fmt.Sprintf("%s %s", password, input)),
	}
	if err = msg.Marshal(buf); err != nil {
		return
	}
	if _, err = socket.WriteToUDP(buf.Bytes(), address); err != nil {
		return
	}
	return
}

func readLineFromInput() (input string, err error) {
	if bufferedStdin == nil {
		bufferedStdin = bufio.NewReader(os.Stdin)
	}
	for {
		if line, hasMoreInLine, err := bufferedStdin.ReadLine(); err != nil {
			return input, err
		} else {
			input += string(line)
			if !hasMoreInLine {
				break
			}
		}
	}
	return
}

func usage() {
	kingpin.Usage()
}

func main() {
	fmt.Println("IceCon - Icedream's RCON Client")
	fmt.Println("\t\u00A9 2016-2017 Carl Kittelberger/Icedream")
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
	initSocket()

	// Set target address if given
	if *argAddressTCP != nil {
		initSocketAddr((*argAddressTCP).String())
	}

	// Get password
	password = *argPasswordStr

	// Run one-off command?
	if *flagCommand != "" {
		// Send
		err := sendRcon(*flagCommand)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Receive
		msg, err := receiveRcon()
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
	if wantGui {
		if err := runGraphicalUi(); err != nil {
			log.Fatal(err)
			return
		}
	} else {
		runConsoleShell()
	}

	if socket != nil {
		socket.Close()
	}

}

func runConsoleShell() {
	for {
		input, err := readLineFromInput()
		if err != nil {
			log.Fatal(err)
			continue
		}

		// "quit" => exit shell
		if strings.EqualFold(strings.TrimSpace(input), "quit") {
			break
		}

		err = sendRcon(input)
		if err != nil {
			log.Println(err)
			continue
		}
		msg, err := receiveRcon()
		if err != nil {
			log.Println(err)
			continue
		}
		switch strings.ToLower(msg.Name) {
		case "print":
			log.Println(string(msg.Data))
		}
	}
}
