package ui

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/icedream/icecon/internal/rcon"
	"github.com/icedream/icecon/internal/ui"
)

func init() {
	ui.RegisterUserInterface(ui.UserInterfaceProvider{
		New: NewConsoleUserInterface,
	})
}

type consoleUserInterface struct {
	rcon          *rcon.Client
	bufferedStdin *bufio.Reader
}

func NewConsoleUserInterface(rconClient *rcon.Client) (ui.UserInterface, error) {
	return &consoleUserInterface{
		rcon:          rconClient,
		bufferedStdin: bufio.NewReader(os.Stdin),
	}, nil
}

func (ui *consoleUserInterface) readLineFromInput() (input string, err error) {
	for {
		if line, hasMoreInLine, err := ui.bufferedStdin.ReadLine(); err != nil {
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

func (ui *consoleUserInterface) Run() error {
	for {
		input, err := ui.readLineFromInput()
		if err != nil {
			log.Fatal(err)
			continue
		}

		// "quit" => exit shell
		if strings.EqualFold(strings.TrimSpace(input), "quit") {
			break
		}

		err = ui.rcon.Send(input)
		if err != nil {
			log.Println(err)
			continue
		}
		msg, err := ui.rcon.Receive()
		if err != nil {
			log.Println(err)
			continue
		}
		switch strings.ToLower(msg.Name) {
		case "print":
			log.Println(string(msg.Data))
		}
	}
	return nil
}
