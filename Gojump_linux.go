//go:build linux
// +build linux

package main

import (
	"OpenJump/ui"
	"github.com/jroimartin/gocui"
	"log"
	"os"
	"strconv"
)

var (
	done     = make(chan struct{})
	terminal = make(chan string, 10000)
)

func main() {

	waittime := 300

	if len(os.Args) == 2{
		if os.Args[1] == "-h"{
			println("Utility to Rotate across multiple openvpn tunnels in an automated fashion.")
			println("Usage: ./OpenJump <Rotate time in seconds>")
			println("If no time is supplied, defaults to 300 seconds")
			os.Exit(0)
		}
		if arg_int, err := strconv.Atoi(os.Args[1]); err == nil {
			waittime = arg_int
		} else {println("Error: Value Supplied is not an integer")}
	}

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
	go ui.Daemon_start(g, terminal, waittime)
	g.MainLoop()

}
