//+build windows

package main

//go:generate ui2walk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"

	"github.com/lxn/walk"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	hasGraphicalUI = true

	flagGui = kingpin.
		Flag("gui", "Run as GUI (runs automatically as GUI if no arguments given, ignored if command flag used)").
		Short('g').Bool()

	guiInitErr  error
	kernel32    *syscall.DLL
	freeConsole *syscall.Proc

	dlg              *mainDialog
	dlgOriginalTitle string

	history      []string
	historyIndex = 0

	bookmarksFile string
	bookmarks     []bookmark
	bookmarksMenu *walk.Menu
)

type bookmark struct {
	Address  string       `json:"address"`
	Password string       `json:"password"`
	Action   *walk.Action `json:"-"`
}

func init() {
	kernel32, guiInitErr = syscall.LoadDLL("kernel32.dll")
	freeConsole, guiInitErr = kernel32.FindProc("FreeConsole")

	// Creates an folder in user home called "IceCon"
	initStorage()
	loadBookmarks()
}

func initStorage() {
	// If user found create in user home
	if usr, err := user.Current(); err == nil {
		directory := path.Join(usr.HomeDir, "IceCon")
		bookmarksFile = path.Join(usr.HomeDir, "Icecon", "bookmarks.json")

		os.Mkdir(directory, os.ModePerm)
	}
}

func loadBookmarks() {
	// Loads bookmarks.json from IceCon directory
	file, err := ioutil.ReadFile(bookmarksFile)
	if err != nil {
		return
	}

	json.Unmarshal(file, &bookmarks)
}

func saveBookmarks() {
	// Saves current bookmarks to file
	if out, err := json.Marshal(bookmarks); err == nil {
		ioutil.WriteFile(bookmarksFile, out, os.ModePerm)
	}
}

func uiLogError(text string) {
	uiNormalize(&text)
	dlg.Synchronize(func() {
		dlg.ui.rconOutput.AppendText("ERROR: " + text + "\r\n")
		walk.MsgBox(dlg, "Error",
			text,
			walk.MsgBoxIconError)
	})
}

func uiLog(text string) {
	uiNormalize(&text)
	dlg.Synchronize(func() {
		dlg.ui.rconOutput.AppendText(text + "\r\n")
	})
}

func uiNormalize(textRef *string) {
	text := *textRef

	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n", "\r\n", -1)

	*textRef = text
}

func uiUpdateAddress() {
	if len(dlgOriginalTitle) <= 0 {
		dlgOriginalTitle = dlg.Title()
	}
	if len(addressStr) > 0 {
		dlg.SetTitle(dlgOriginalTitle + " - " + addressStr)
	} else {
		dlg.SetTitle(dlgOriginalTitle)
	}
}

func addToHistory(command string) {
	// limit history to 20 items
	if len(history) > 20 {
		history = append(history[:0], history[0+1:]...)
	}

	history = append(history, command)
	historyIndex = len(history)
}

func getBookmark(address string) (bookmark, error) {
	var bookmark bookmark

	// Gets the bookmark by address
	for _, item := range bookmarks {
		if item.Address == address {
			return item, nil
		}
	}

	return bookmark, errors.New("No bookmark found")
}

func createBookmarkItem(address string) *walk.Action {
	// Create bookmark item (used in init UI and add bookmark)
	item := walk.NewAction()
	item.SetText(address)

	// Connect to selected server
	item.Triggered().Attach(func() {
		// Get the right bookmark
		bookmark, err := getBookmark(item.Text())
		if err != nil {
			return
		}
		if err = initSocketAddr(bookmark.Address); err != nil {
			uiLogError(fmt.Sprintf("Couldn't use that address: %s", err))
			return
		}

		password = bookmark.Password
		addressStr = bookmark.Address
		dlg.ui.rconOutput.SetText("")

		uiUpdateAddress()

		// Uncheck other items
		for i := 0; i < bookmarksMenu.Actions().Len(); i++ {
			item := bookmarksMenu.Actions().At(i)
			item.SetChecked(false)
		}
		item.SetChecked(true)
	})

	return item
}

func runGraphicalUi() (err error) {
	dlg = new(mainDialog)
	if err := dlg.init(); err != nil {
		panic(err)
	}
	defer dlg.Dispose()

	// Window icon
	// TODO - Do this more intelligently
	for i := 0; i < 128; i++ {
		if icon, err := walk.NewIconFromResourceId(i); err == nil {
			dlg.SetIcon(icon)
			break
		}
	}

	// Quit button
	quitAction := walk.NewAction()
	if err = quitAction.SetText("&Quit"); err != nil {
		return
	}
	quitAction.Triggered().Attach(func() { dlg.Close() })
	if err = dlg.Menu().Actions().Add(quitAction); err != nil {
		return
	}

	// Connect button
	connectAction := walk.NewAction()
	if err = connectAction.SetText("&Connect"); err != nil {
		return
	}
	connectAction.Triggered().Attach(func() {
		result, addr, pw, err := runConnectDialog(addressStr, password, dlg, "")
		if err != nil {
			uiLogError(fmt.Sprintf("Failed to run connect dialog: %s", err))
			return
		}
		if result {
			if err = initSocketAddr(addr); err != nil {
				uiLogError(fmt.Sprintf("Couldn't use that address: %s", err))
				return
			}
			password = pw
			dlg.ui.rconOutput.SetText("")
			uiUpdateAddress()
		}
	})
	if err = dlg.Menu().Actions().Add(connectAction); err != nil {
		return
	}

	// Bookmarks list
	bookmarksMenu, _ = walk.NewMenu()

	// Add for every bookmark a menu action
	for i, bookmark := range bookmarks {
		item := createBookmarkItem(bookmark.Address)

		bookmarks[i].Action = item
		bookmarksMenu.Actions().Add(item)
	}

	// Add bookmark item
	addBookmarkAction := walk.NewAction()
	if err = addBookmarkAction.SetText("New bookmark"); err != nil {
		return
	}
	addBookmarkAction.Triggered().Attach(func() {
		// Re-use connect dialog as "add bookmark" dialog
		result, addr, pw, err := runConnectDialog("", "", dlg, "Add bookmark")
		if err != nil {
			uiLogError(fmt.Sprintf("Failed to run connect dialog: %s", err))
			return
		}

		// Add bookmark if not empty addr
		if result && len(addr) > 0 {
			item := createBookmarkItem(addr)
			bookmarks = append(bookmarks, bookmark{addr, pw, item})

			// Add before "add bookmark" action
			bookmarksMenu.Actions().Insert((bookmarksMenu.Actions().Len() - 2), item)

			// Save bookmarks to file
			saveBookmarks()
		}
	})

	removeBookmarksAction := walk.NewAction()
	if err = removeBookmarksAction.SetText("Remove all"); err != nil {
		return
	}
	removeBookmarksAction.Triggered().Attach(func() {
		for _, bookmark := range bookmarks {
			bookmarksMenu.Actions().Remove(bookmark.Action)
		}

		// Clear bookmarks in memory
		bookmarks = bookmarks[:0]

		// Update in file
		saveBookmarks()
	})

	bookmarksMenu.Actions().Add(addBookmarkAction)
	bookmarksMenu.Actions().Add(removeBookmarksAction)

	// Bookmarks menu
	bookmarksAction, _ := dlg.Menu().Actions().AddMenu(bookmarksMenu)
	if err = bookmarksAction.SetText("Bookmarks"); err != nil {
		return
	}

	// Handle input
	dlg.ui.rconInput.KeyPress().Attach(func(key walk.Key) {
		// handle history (arrow up/down)
		if key == walk.KeyUp || key == walk.KeyDown {
			if len(history) == 0 {
				return
			}

			if key == walk.KeyUp {
				if historyIndex == 0 {
					return
				}

				historyIndex -= 1
				dlg.ui.rconInput.SetText(history[historyIndex])
			} else {
				if (historyIndex + 1) >= len(history) {
					return
				}

				historyIndex += 1
				dlg.ui.rconInput.SetText(history[historyIndex])
			}

			return
		}

		if key != walk.KeyReturn {
			return
		}

		if address == nil {
			uiLogError("No server configured.")
			return
		}

		cmd := dlg.ui.rconInput.Text()
		dlg.ui.rconInput.SetText("")

		uiLog(address.String() + "> " + cmd)
		sendRcon(cmd)

		// add to history
		addToHistory(cmd)
	})

	// When window is initialized we can let a secondary routine print all
	// output received
	dlg.Synchronize(func() {
		uiUpdateAddress()

		go func() {
			for {
				msg, err := receiveRcon()
				if err != nil {
					uiLogError(err.Error())
					continue
				}
				switch strings.ToLower(msg.Name) {
				case "print":
					uiLog(string(msg.Data))
				default:
					log.Println(msg.Name)
				}
			}
		}()
	})

	// Get rid of the console window
	freeConsole.Call()

	dlg.Show()

	// Message loop starts here and will block the main goroutine!
	dlg.Run()

	return
}
