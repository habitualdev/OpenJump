package ui

import (
	"OpenJump/openvpn"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"os"
	"time"
)

// CLI UI Setup

func Quit(g *gocui.Gui, v *gocui.View) error {
	g.Close()
	os.Exit(0)
	return gocui.ErrQuit
}

func Daemon_start(g *gocui.Gui, terminal chan string) {
	term := make(chan bool)
	done := make(chan struct{})
	var scroll []string
	log_file, _ := os.OpenFile("openvpn.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer log_file.Close()
	wrt := io.MultiWriter(log_file)
	log.SetOutput(wrt)
	log.Println("Opened log file:", time.Now())
	config_list := openvpn.GetConfigs()
	go openvpn.WheelInTheSky(10000, wrt, term, terminal, config_list)
	//go openvpn.Tunnel("hops/japan.ovpn",10000,wrt,term,terminal)
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("banner")
				if err != nil {return err}
				select {
				case msg := <-terminal:

					scroll = append(scroll, msg)
						}
					v.Clear()
					for _, i := range scroll {
						fmt.Fprintln(v, i)
					}
				return nil
			})
		case <-done:
			return

		}
	}
}
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("banner", 2, 0, maxX-55, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprintln(v, "Hello world!")
	}

	return nil
}
