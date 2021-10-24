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
	done := make(chan struct{})
	config_chan := make(chan string, 1000)
	var scroll []string

	log_file, _ := os.OpenFile("openvpn.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer log_file.Close()
	wrt := io.MultiWriter(log_file)
	log.SetOutput(wrt)
	log.Println("OpenJump started.")
	config_list := openvpn.GetConfigs()
	go openvpn.Wheel(config_list, 300, terminal, config_chan)
	for {
		select {
		case config_msg := <-config_chan:
			UpdateConfigTitle(g, config_msg)

		case <-time.After(100 * time.Millisecond):
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("banner")
				if err != nil {
					return err
				}
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

func UpdateConfigTitle(g *gocui.Gui, configname string) error {
	v, err := g.View("config")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintln(v, configname)
	return err
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("banner", 2, 0, maxX-2, maxY-22); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprintln(v, "Hello world!")
	}

	if v2, err := g.SetView("config", 2, maxY-20, 20, maxY-18); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v2.Wrap = true
		v2.Autoscroll = true

	}

	return nil
}
