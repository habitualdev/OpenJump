package ui

import (
	"OpenJump/libopenvpn"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// CLI UI Setup

func Quit(g *gocui.Gui, v *gocui.View) error {
	g.Close()
	os.Exit(0)
	return gocui.ErrQuit
}

func Daemon_start(g *gocui.Gui, terminal chan string, waittime int) {
	//done := make(chan struct{})
	config_chan := make(chan string, 1000)
	var scroll []string
	count := 0

	log_file, _ := os.OpenFile("libopenvpn.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer log_file.Close()
	wrt := io.MultiWriter(log_file)
	log.SetOutput(wrt)
	log.Println("OpenJump started.")
	config_list := libopenvpn.GetConfigs()
	go libopenvpn.Wheel(config_list, waittime, terminal, config_chan)
	for {
		select {
		case config_msg := <-config_chan:
			UpdateConfigTitle(g, config_msg)
		case msg := <-terminal:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("banner")
				if err != nil {
					return err
				}
					scroll = append(scroll, msg)
				v.Clear()
				for _, i := range scroll {
					fmt.Fprintln(v, i)
				}
				return nil
			})
		case <- time.After(time.Duration(1)*time.Second):
			count = count + 1
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("uptime")
				if err != nil {
					return err
				}
				v.Clear()
				write_string := "Uptime: " + strconv.Itoa(count) + " seconds"
				v.Write([]byte(write_string))
				return nil
			})
		//case <-done:
		//	return

		}
	}
}

func UpdateConfigTitle(g *gocui.Gui, configname string) error {
	v, err := g.View("config")
	buf := v.Buffer()
	buf = buf + string('\n') + configname
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintln(v,buf)
	return err
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("banner", 2, 0, maxX-2, maxY-22); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Log"
		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprintln(v, "Hello world!")
	}

	if v2, err := g.SetView("config", 2, maxY-20, 20, maxY-16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v2.Wrap = true
		v2.Autoscroll = true
		v2.Title = "Current Config"
	}

	if v3, err := g.SetView("uptime", maxX-24, maxY-20, maxX-2, maxY-16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v3.Title = "Uptime"
		v3.Wrap = true
		v3.Autoscroll = true
		fmt.Fprintln(v3,"0")
	}

	if v4, err := g.SetView("config_list", 2, maxY-14, maxX-2, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v4.Title = "Config List"
		v4.Wrap = true
		v4.Autoscroll = true
		var list []string
		config_list := libopenvpn.GetConfigs()
		for _, config_file := range config_list {
			list = append(list, config_file.Name())
			}
		for i, list_entry := range list {
			entry_string := strconv.Itoa(i + 1) + ".) " + list_entry
			fmt.Fprintln(v4, entry_string)
		}
	}

	if v5, err := g.SetView("vanity", 22, maxY-20, maxX-26, maxY-16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v5.Wrap = true
		v5.Autoscroll = true
		sizex, _ := v5.Size()
		offset := (sizex/2)-10
		space := strings.Repeat(" ", offset)

		fmt.Fprintln(v5,space + "OpenJump by Habitual")
	}


	return nil
}
