package openvpn

import (
	"fmt"
	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"time"
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

func GetConfigs() []fs.FileInfo {
	dir_check, err := ioutil.ReadDir("hops")
	if err != nil {
		os.Mkdir("hops", 0777)
		dir_check, err = ioutil.ReadDir("hops")
	}
	return dir_check

}

func Tunnel(config_file string, waittime int, wrt io.Writer, terminate chan bool) {

	var logger StdoutLogger = func(text string) {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Fprintln(wrt, "Library check >>", line)
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
	if err != nil {
		fmt.Println("Openvpn3 error: ", err)
	} else {
		fmt.Println("Graceful exit")
	}
	for {
		select {
		case <-terminate:
			session.Stop()
		case <-time.After(time.Duration(waittime) * time.Millisecond):
			session.Stop()

		}
	}
}

func WheelInTheSky(waittime int, wrt io.Writer, terminate chan bool) {
	config_list := GetConfigs()
	for _, config_file := range config_list {
		filename := config_file.Name()
		pathname := "hops/" + filename
		Tunnel(pathname, waittime, wrt, terminate)

	}
}
