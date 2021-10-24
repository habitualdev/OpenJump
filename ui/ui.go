package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"os"
	"time"
)

// CLI UI Setup

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func daemon_start(g *gocui.Gui, terminal chan string) {
	done := make(chan struct{})
	var scroll []string

	for {
		select {
		case <-time.After(500 * time.Millisecond):
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("banner")
				if err != nil {return err}
				select {
				case msg := <-terminal:
					scroll = append(scroll, msg)
					v.Clear()
					for _, i := range scroll {
						fmt.Fprintln(v, i)
					}
				default:
					time.Sleep(1*time.Millisecond)
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
	if v, err := g.SetView("banner", 2, maxY-20, maxX-55, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, "Hello world!")
	}
	if v2, err := g.SetView("loaded_config",maxX-50, 0, maxX-5, maxY-2 ); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v2.Wrap = true



	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	os.Exit(0)
	return gocui.ErrQuit
}
