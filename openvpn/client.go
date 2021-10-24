package openvpn

import (
	"fmt"
	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)
var (
	Mesg = make(chan string, 10000))

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
		//fmt.Println(Mesg,"Openvpn log >>", line)
		Mesg <- line
	}
}

func (lc *loggingCallbacks) OnEvent(event openvpn3.Event) {
	//fmt.Printf("Openvpn event >> %+v\n", event)
	Mesg <- event.Info
}

func (lc *loggingCallbacks) OnStats(stats openvpn3.Statistics) {
	//fmt.Printf("Openvpn stats >> %+v\n", stats)
	stats_line := "Bytes in: " + strconv.FormatUint(stats.BytesIn, 10) + " | Bytes out: " + strconv.FormatUint(stats.BytesOut, 10)
	Mesg <- stats_line
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

func Tunnel(config_file string, waittime int, wrt io.Writer, terminate chan bool, terminal chan string) {
	var i int
	fmt.Fprintln(wrt, "Constructing Tunnel")
	//sleep_time, _ := time.ParseDuration(string(waittime)+"s")
	var logger StdoutLogger = func(text string) {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Fprintln(wrt, "Library check >>", line)
			Mesg <- line
		}
	}

	openvpn3.SelfCheck(logger)

	bytes, err := ioutil.ReadFile(config_file)

	if err != nil {
		fmt.Fprintln(wrt,err.Error())
		os.Exit(1)
	}


	config := openvpn3.NewConfig(string(bytes))
	session := openvpn3.NewSession(config, openvpn3.UserCredentials{}, &loggingCallbacks{})
	session.Start()
	if err != nil {
		Mesg <- err.Error()
		//fmt.Println("Openvpn3 error: ", err)
	} else {
		Mesg <- "Graceful Exit"
		//fmt.Println("Graceful exit")
	}
	waitfor := time.Duration(waittime)*time.Millisecond
	for start := time.Now(); ;{
		if i % 10 == 0{if time.Since(start) > waitfor {break}}
		select {
		case <-time.After(100 * time.Millisecond):

		case mesg_term := <-Mesg:
			terminal <- mesg_term

		}
	}
	session.Stop()
}

func WheelInTheSky(waittime int, wrt io.Writer, terminate chan bool, terminal chan string, config_list []fs.FileInfo) {
	fmt.Fprintln(wrt,"wheel started")
	fmt.Fprintln(wrt,"Configs Retrieved")
	for{
		for _, config_file := range config_list {
			filename := config_file.Name()
			fmt.Fprintln(wrt, filename)
			pathname := "hops/" + filename
			Tunnel(pathname, waittime, wrt, terminate, terminal)
		}
	}
}
