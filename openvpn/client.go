package openvpn

import (
	"fmt"
	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type callbacks interface {
	openvpn3.Logger
	openvpn3.EventConsumer
	openvpn3.StatsConsumer
}

type loggingCallbacks struct {
}

func (lc *loggingCallbacks) Log(text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fmt.Println("Openvpn log >>", line)
	}
}

func (lc *loggingCallbacks) OnEvent(event openvpn3.Event) {
	fmt.Printf("Openvpn event >> %+v\n", event)
}

func (lc *loggingCallbacks) OnStats(stats openvpn3.Statistics) {
	fmt.Printf("Openvpn stats >> %+v\n", stats)
}

var _ callbacks = &loggingCallbacks{}

// StdoutLogger represents the stdout logger callback
type StdoutLogger func(text string)

// Log logs the given string to stdout logger
func (lc StdoutLogger) Log(text string) {
	lc(text)
}

func Tunnel(config_file string, wrt io.Writer, term chan bool) {
	if len(os.Args) < 2 {
		fmt.Println("Missing profile file")
		fmt.Fprintln(wrt, "Missing profile file")
		go func() { term <- true }()
	}


	var logger StdoutLogger = func(text string) {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Println("Library check >>", line)
		}
	}

	openvpn3.SelfCheck(logger)

	bytes, err := ioutil.ReadFile(config_file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	config := openvpn3.NewConfig(string(bytes))

	session := openvpn3.NewSession(config, openvpn3.UserCredentials{}, &loggingCallbacks{})
	session.Start()
	err = session.Wait()
	if err != nil {
		fmt.Println("Openvpn3 error: ", err)
	} else {
		fmt.Println("Graceful exit")
	}

}