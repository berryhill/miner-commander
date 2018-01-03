package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

type ClaymoreLog struct {
	Timestamp 			string
	Source 					string
	Error						string
	// Level 				string
	// Code 				string
}

func parseClaymoreLog(error string) *ClaymoreLog {
	cl := new(ClaymoreLog)
  cl_array := strings.Split(error, "	")
	timestamp_array := strings.Split(cl_array[0], ":")
	if len(timestamp_array) == 4 && len(cl_array) >= 3 {
		cl.Timestamp = cl_array[0]
		cl.Source = cl_array[1]
		cl.Error = cl_array[2]

		return cl
	}

	return nil
}

func main() {
	f, err := os.OpenFile(
		"miner_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
    fmt.Println("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	setupMiner()
	startMiner()

	claymoreLog := make(chan string)
	go tailLogs(claymoreLog)

	for {
		cl := parseClaymoreLog(<-claymoreLog)
		if cl != nil {
			// fmt.Println(cl.Error)
			log.Println(cl.Error)
			error_array := strings.Split(cl.Error, ", ")
			if len(error_array) > 1 {
				if error_array[0] == "NVML: cannot get current temperature" {
					log_reboot()
					// reboot()
				} else if error_array[0] == "NVML: cannot get fan speed" {
					log_reboot()
					// reboot()
				}
			}
		}
	}
}

func logToSlack() {
	//todo: IMPLEMENT
}

func reboot() {
	log.Println("Rebooting Miner")
	out, err := exec.Command("/bin/sh", "./reboot.sh").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf(string(out))
}

func setupMiner() {
	log.Println("Setting up Miner")
	// out, err := exec.Command("/bin/sh", "../../mine-setup.sh").Output()
  //   if err != nil {
  //       log.Fatal(err)
  //   }
  //   log.Printf(string(out))
}

func startMiner() {
	log.Println("Starting Miner")
	// out, err := exec.Command("/bin/sh", "../start.bash").Output()
  //   if err != nil {
  //       log.Fatal(err)
  //   }
  //   log.Printf(string(out))
}

func tailLogs(ch chan string) {
		t, _ := tail.TailFile("./claymore/log.txt", tail.Config{Follow: true})
		for line := range t.Lines {
				ch<- line.Text
		}
}

func testLog() {
	for {
		log.Println("Hello world!")
		duration := time.Second
		time.Sleep(duration)
	}
}

func log_reboot() {
	log.Println("------ REBOOT / REBOOT / REBOOT / REBOOT ------")
	// fmt.Println("------ REBOOT / REBOOT / REBOOT / REBOOT ------")
}
