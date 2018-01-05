package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/gorilla/websocket"
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

var addr = flag.String("addr", "10.0.0.128:8080", "http service address")

func main() {
	f, err := os.OpenFile(
		"/home/berry/mine/claymore/logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
    fmt.Println("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	claymoreLog := make(chan string)
	go tailLogs(claymoreLog)

	// --

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	go func() {
		for {
			select {
			case <-claymoreLog:
			// cl := parseClaymoreLog(<-claymoreLog)
			// if cl != nil {
			// 	log.Println(cl.Error)
				err := c.WriteMessage(websocket.TextMessage, []byte(<-claymoreLog))
				if err != nil {
					log.Println("write:", err)
					return
				}
				// error_array := strings.Split(cl.Error, ", ")
				// if len(error_array) > 1 {
				// 	if error_array[0] == "NVML: cannot get current temperature" {
				// 		log_reboot()
				// 		// reboot()
				// 	} else if error_array[0] == "NVML: cannot get fan speed" {
				// 		log_reboot()
				// 		// reboot()
				// 	}
				// }
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
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
	out, err := exec.Command("/bin/sh", "./mine-setup.sh").Output()
    if err != nil {
        // log.Fatal(err)
    }
    log.Printf(string(out))
}

func startMiner() {
	log.Println("Starting Miner")
	out, err := exec.Command("/bin/sh", "cd claymore && nohup ./start.bash > logs.txt &").Output()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf(string(out))
}

func tailLogs(ch chan string) {
		t, _ := tail.TailFile("/home/berry/mine/claymore/logs.txt", tail.Config{Follow: true})
		for line := range t.Lines {
				ch<- line.Text
		}
}

func testLog() {
	for {
		// log.Println("Hello world!")
		duration := time.Second
		time.Sleep(duration)
	}
}

func log_reboot() {
	log.Println("------ REBOOT / REBOOT / REBOOT / REBOOT ------")
}
