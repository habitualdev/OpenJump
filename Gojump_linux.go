//go:build linux
// +build linux

package main

import (
	"OpenJump/ui"
	"github.com/jroimartin/gocui"
	"log"
)

var (
	done     = make(chan struct{})
	terminal = make(chan string, 1000)
)

func main() {

	// Init the gui
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(ui.Layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		log.Panicln(err)
	}
	go ui.Daemon_start(g, terminal)
	g.MainLoop()

}
