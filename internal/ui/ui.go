package ui

import (
	"errors"
	"fmt"
	"slices"

	"github.com/icedream/icecon/internal/rcon"
)

// ErrNotSupported is returned on calls to graphical user interface routines
// when none provide one.
var ErrNotSupported = errors.New("not supported")

// All registered user interface providers.
var userInterfaces = []*UserInterfaceProvider{}

// UserInterface describes the methods implemented by user interfaces.
type UserInterface interface {
	Run() error
}

// UserInterfaceProvider contains metadata and an entrypoint for a provided
// user interface.
type UserInterfaceProvider struct {
	// Whether the user interface renders as a graphical user interface.
	IsGraphical bool
	New         func(
		rconClient *rcon.Client,
	) (UserInterface, error)
}

// RegisterUserInterface must be called on init by any user interface provider
// loaded in to be discoverable by other methods in this package.
func RegisterUserInterface(uiDesc UserInterfaceProvider) {
	userInterfaces = append(userInterfaces, &uiDesc)
}

// HasGraphicalUI returns true if at least one user interface provider provides
// a graphical user interface.
func HasGraphicalUI() bool {
	return slices.IndexFunc(
		userInterfaces,
		func(ui *UserInterfaceProvider) bool {
			return ui.IsGraphical
		}) >= 0
}

// Run scans for a matchin user interface provider. If it finds one, it will run
// the user interface through it, otherwise it will return nil.
func Run(rconClient *rcon.Client, wantGraphical bool) error {
	for _, uiDesc := range userInterfaces {
		if uiDesc.IsGraphical != wantGraphical {
			continue
		}
		ui, err := uiDesc.New(rconClient)
		if err != nil {
			return fmt.Errorf("failed to instantiate user interface: %w", err)
		}
		return ui.Run()
	}
	return ErrNotSupported
}
