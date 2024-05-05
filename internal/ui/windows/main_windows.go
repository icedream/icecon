//go:build windows
// +build windows

package windows

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/icedream/icecon/internal/rcon"
	"github.com/icedream/icecon/internal/ui"
	walk "github.com/lxn/walk"
)

var (
	kernel32    *syscall.DLL
	freeConsole *syscall.Proc

	initErr error
)

func init() {
	kernel32, initErr = syscall.LoadDLL("kernel32")
	if initErr != nil {
		return
	}
	freeConsole, initErr = kernel32.FindProc("FreeConsole")
	if initErr != nil {
		return
	}

	ui.RegisterUserInterface(ui.UserInterfaceProvider{
		IsGraphical: true,
		New:         NewWindowsUserInterface,
	})
}

type windowsUserInterface struct {
	rcon *rcon.Client

	mainDialog          *mainDialog
	originalDialogTitle string

	history      []string
	historyIndex int
}

func (ui *windowsUserInterface) logError(text string) {
	text = normalizeTextForUI(text)
	ui.mainDialog.Synchronize(func() {
		ui.mainDialog.ui.rconOutput.AppendText("ERROR: " + text + "\r\n")
		walk.MsgBox(ui.mainDialog, "Error",
			text,
			walk.MsgBoxIconError)
	})
}

func (ui *windowsUserInterface) log(text string) {
	text = normalizeTextForUI(text)
	ui.mainDialog.Synchronize(func() {
		ui.mainDialog.ui.rconOutput.AppendText(text + "\r\n")
	})
}

func normalizeTextForUI(text string) string {
	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n", "\r\n", -1)

	return text
}

func (ui *windowsUserInterface) updateAddress() {
	if len(ui.originalDialogTitle) <= 0 {
		ui.originalDialogTitle = ui.mainDialog.Title()
	}
	if len(ui.rcon.AddressString()) > 0 {
		ui.mainDialog.SetTitle(ui.originalDialogTitle + " - " + ui.rcon.AddressString())
	} else {
		ui.mainDialog.SetTitle(ui.originalDialogTitle)
	}
}

func (ui *windowsUserInterface) addToHistory(command string) {
	// limit history to 20 items
	if len(ui.history) > 20 {
		ui.history = append(ui.history[:0], ui.history[0+1:]...)
	}

	ui.history = append(ui.history, command)
	ui.historyIndex = len(ui.history)
}

func (ui *windowsUserInterface) Dispose() {
	if ui.mainDialog != nil {
		ui.mainDialog.Dispose()
		ui.mainDialog = nil
	}
}

func NewWindowsUserInterface(rconClient *rcon.Client) (ui.UserInterface, error) {
	if initErr != nil {
		return nil, initErr
	}

	return &windowsUserInterface{
		rcon: rconClient,
	}, nil
}

func (ui *windowsUserInterface) Run() error {
	var err error
	defer func() {
		if err != nil {
			ui.Dispose()
		}
	}()

	ui.mainDialog = new(mainDialog)
	if err := ui.mainDialog.init(); err != nil {
		return err
	}

	// Window icon
	// TODO - Do this more intelligently
	for i := 0; i < 128; i++ {
		if icon, err := walk.NewIconFromResourceId(i); err == nil {
			ui.mainDialog.SetIcon(icon)
			break
		}
	}

	// Quit button
	quitAction := walk.NewAction()
	if err = quitAction.SetText("&Quit"); err != nil {
		return err
	}
	quitAction.Triggered().Attach(func() { ui.mainDialog.Close() })
	if err = ui.mainDialog.Menu().Actions().Add(quitAction); err != nil {
		return err
	}

	// Connect button
	connectAction := walk.NewAction()
	if err = connectAction.SetText("&Connect"); err != nil {
		return err
	}
	connectAction.Triggered().Attach(func() {
		result, addr, pw, err := runConnectDialog(
			ui.rcon.AddressString(),
			ui.rcon.Password(),
			ui.mainDialog)
		if err != nil {
			ui.logError(fmt.Sprintf("Failed to run connect dialog: %s", err))
			return
		}
		if result {
			if err = ui.rcon.SetSocketAddr(addr); err != nil {
				ui.logError(fmt.Sprintf("Couldn't use that address: %s", err))
				return
			}
			ui.rcon.SetPassword(pw)
			ui.mainDialog.ui.rconOutput.SetText("")
			ui.updateAddress()
		}
	})
	if err = ui.mainDialog.Menu().Actions().Add(connectAction); err != nil {
		return err
	}

	// Handle input
	ui.mainDialog.ui.rconInput.KeyPress().Attach(func(key walk.Key) {
		// handle history (arrow up/down)
		if key == walk.KeyUp || key == walk.KeyDown {
			if len(ui.history) == 0 {
				return
			}

			if key == walk.KeyUp {
				if ui.historyIndex == 0 {
					return
				}

				ui.historyIndex -= 1
				ui.mainDialog.ui.rconInput.SetText(ui.history[ui.historyIndex])
			} else {
				if (ui.historyIndex + 1) >= len(ui.history) {
					return
				}

				ui.historyIndex += 1
				ui.mainDialog.ui.rconInput.SetText(ui.history[ui.historyIndex])
			}

			return
		}

		if key != walk.KeyReturn {
			return
		}

		if ui.rcon.Address() == nil {
			ui.logError("No server configured.")
			return
		}

		cmd := ui.mainDialog.ui.rconInput.Text()
		ui.mainDialog.ui.rconInput.SetText("")

		ui.log(ui.rcon.Address().String() + "> " + cmd)
		ui.rcon.Send(cmd)

		// add to history
		ui.addToHistory(cmd)
	})

	// When window is initialized we can let a secondary routine print all
	// output received
	ui.mainDialog.Synchronize(func() {
		ui.updateAddress()

		go func() {
			for {
				msg, err := ui.rcon.Receive()
				if err != nil {
					ui.logError(err.Error())
					continue
				}
				switch strings.ToLower(msg.Name) {
				case "print":
					ui.log(string(msg.Data))
				default:
					log.Println(msg.Name)
				}
			}
		}()
	})

	// Get rid of the console window
	// 	freeConsole.Call()

	ui.mainDialog.Show()

	// Message loop starts here and will block the main goroutine!
	if retval := ui.mainDialog.Run(); retval != 0 {
		err = syscall.Errno(retval)
	}

	return err
}
