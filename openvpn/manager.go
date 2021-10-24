package openvpn

import (
	"bufio"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func GetConfigs() []fs.FileInfo {
	dir_check, err := ioutil.ReadDir("hops")
	if err != nil {
		os.Mkdir("hops", 0777)
		dir_check, err = ioutil.ReadDir("hops")
	}
	return dir_check
}

func checkError(err error, terminal chan string) {
	if err != nil {
		terminal <- err.Error()
	}
}

func KillProcess(pid int32) error {

	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		id := p.Pid
		if err != nil {
			return err
		}
		if id == pid {
			return p.Kill()
		}
	}
	return fmt.Errorf("process not found")
}

func comms(input io.Reader, config_name string, waittime int, terminal chan string) {
	waitfor := time.Duration(waittime) * time.Second
	reader := bufio.NewReader(input)
	line, _ := reader.ReadString('\n')
	for start := time.Now(); ; {
		if time.Since(start) > waitfor {
			break
		}
		terminal <- line
		line, _ = reader.ReadString('\n')
	}
	terminal <- "Exited Tunnel - " + config_name
}

func Start_process(config_file string, waittime int, terminal chan string) {

	cmd := exec.Command("./openvpn_bin", config_file)
	stdout, err := cmd.StdoutPipe()
	checkError(err, terminal)
	err = cmd.Start()
	checkError(err, terminal)
	pid := cmd.Process.Pid
	terminal <- "Tunnel for " + config_file + " created."
	go comms(stdout, config_file, waittime, terminal)
	waitfor := time.Duration(waittime) * time.Second
	time.Sleep(waitfor)
	err = KillProcess(int32(pid))
	if err != nil {
		println(err.Error())
		return
	}

}

func Wheel(config_list []os.FileInfo, waittime int, terminal chan string, config_chan chan string) {
	for {
		for _, config_file := range config_list {
			filename := config_file.Name()
			config_chan <- filename
			pathname := "hops/" + filename
			Start_process(pathname, waittime, terminal)
		}
	}
}
