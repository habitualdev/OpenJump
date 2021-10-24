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
	//exit_loop := true
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
	go g.MainLoop()
	for true {
		select {
			case <- done:
				//exit_loop = false
			default:
				openvpn.WheelInTheSky(1000, wrt, term)
				time.Sleep(500*time.Millisecond)
		}
	}
		// Start the background services
		g.Close()

		//Launch the gui

	}

