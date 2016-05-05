//+build !windows

package main

var flagGuiIncompatible = false
var flagGui = &flagGuiIncompatible

var hasGraphicalUI = false

func runGraphicalUi() error {
	return errNotSupported
}
