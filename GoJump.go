package main

import (
	"OpenJump/openvpn"
	"OpenJump/ui"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"os"
	"time"
)

var (
	done = make(chan struct{})
	terminal = make(chan string)
	term = make(chan bool)
	)

func main(){

	log_file, _ := os.OpenFile("openvpn.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer log_file.Close()
	wrt := io.MultiWriter(log_file)
	log.SetOutput(wrt)
	log.Println("Opened log file:", time.Now())

	// Init the gui
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {log.Panicln(err)}
	defer g.Close()
	g.SetManagerFunc(ui.Layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {log.Panicln(err)}
	openvpn.Tunnel("euro-hop.ovpn", wrt, term)
	exit_test := <- term
	// Start the background services
	if exit_test{g.Close()}

	//Launch the gui
	g.MainLoop()

}