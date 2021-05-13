/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2019
 */

// websockets.go
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10 * 1024,
	WriteBufferSize: 10 * 1024,
}

var (
	serverAddr string
	started    bool
	hpaMode    bool
	option     string
	version    string
	cpuPod     string
	cpudiv     int
	iclients   string
	clientStep string
	lclients   string
	time       string
	sla        string
	runtimes   int
)

const (
	HTMLTEMPLATE = "cnbweb.html"
	DEFAULTADDR  = "localhost"
	CONFIGFILE   = "config.json"
	ERRORPRE     = "Error found: "
)

type Autoloader struct {
	InitialC string `json:"initialclients"`
	ClientsS string `json:"clientstep"`
	LastC    string `json:"lastclients"`
	SLA      string `json:"SLA"`
	Interval string `json:"timeinterval"`
}

type Workload struct {
	Version string `json:"version"`
	CPU     string `json:"cpurequests"`
}

type Config struct {
	RunOption   string     `json:"runoption"`
	Iterations  int        `json:"iterations"`
	HPAMode     bool       `json:"hpamode"`
	PostProcess bool       `json:"postprocess"`
	PPFile      string     `json:"ppoutputfile"`
	Loader      Autoloader `json:"autoloader"`
	Load        Workload   `json:"workload"`
}

func init() {
	flag.StringVar(&serverAddr, "server", "", "CNB web server IP address")
}

func main() {
	flag.Parse()

	http.HandleFunc("/cnb", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatalln(err)
			}

			if started {
				continue
			}

			errors := precheckInput(string(msg))
			if len(errors) > 0 {
				// Write errors back to browser
				if err = conn.WriteMessage(msgType, []byte(ERRORPRE+errors)); err != nil {
					log.Fatalln(err)
				}
				continue
			}

			// write to config.json file which will be used by cnbrun
			writeToJSON()

			// Mark CNB stared to run
			started = true
			cmd := exec.Command("./cnbrun")

			// create a pipe for the output of the script
			cmdReader, err := cmd.StdoutPipe()
			if err != nil {
				log.Fatalln("Error creating StdoutPipe for Cmd", err)
			}

			scanner := bufio.NewScanner(cmdReader)
			go func() {
				for scanner.Scan() {
					// Write message back to browser
					if err = conn.WriteMessage(msgType, scanner.Bytes()); err != nil {
						log.Fatalln(err)
					}
				}
			}()

			err = cmd.Start()
			if err != nil {
				log.Fatalln("Error starting Cmd", err)
			}

			err = cmd.Wait()
			if err != nil {
				log.Fatalln("Error waiting for Cmd", err)
			}

			started = false
		}
	})

	// Need preprocess HTML file
	processHTMLFile()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, HTMLTEMPLATE)
	})

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func writeToJSON() {
	data := Config{
		RunOption:   option,
		Iterations:  runtimes,
		HPAMode:     hpaMode,
		PostProcess: false,
		PPFile:      "",
		Loader: Autoloader{
			InitialC: iclients,
			ClientsS: clientStep,
			LastC:    lclients,
			SLA:      sla,
			Interval: time,
		},
		Load: Workload{
			Version: version,
			CPU:     cpuPod,
		},
	}
	file, _ := json.MarshalIndent(data, "", "    ")

	_ = ioutil.WriteFile(CONFIGFILE, file, 0644)
}

func precheckInput(input string) string {
	var err error
	var inputs = strings.Split(input, "#")
	for i, item := range inputs {
		inputs[i] = strings.TrimSpace(item)
		if len(inputs[i]) == 0 {
			return "Empty input field found."
		}
	}

	if inputs[2] == "true" {
		hpaMode = true
	} else {
		hpaMode = false
	}

	if inputs[3] == "user" {
		return "User mode is not supported now"
	} else {
		option = inputs[3]
	}

	version = inputs[4]

	_, err = strconv.Atoi(inputs[5])
	if err != nil {
		return "CPU per POD should be integer: " + err.Error()
	}

	cpuPod = inputs[5]

	sclient, err := strconv.Atoi(inputs[6])
	if err != nil {
		return "Initial client number should be integer: " + err.Error()
	}

	step, err := strconv.Atoi(inputs[7])
	if err != nil || step < 1 {
		return "Client step number should be integer greater or equal to 1: " + err.Error()
	}

	eclient, err := strconv.Atoi(inputs[8])
	if err != nil {
		return "Last client number should be integer: " + err.Error()
	}

	if eclient <= sclient && eclient != -1 {
		return "Last client number should be bigger than initial client number"
	}

	iclients = inputs[6]
	clientStep = inputs[7]
	lclients = inputs[8]

	interval, err := strconv.Atoi(inputs[9])
	if err != nil {
		return "Time interval should be an integer: " + err.Error()
	} else if interval < 1 {
		return "Time interval should be integer greater or equal to 1"
	}
	time = inputs[9]

	_, err = strconv.ParseFloat(inputs[10], 32)
	if err != nil {
		return "SLA should be a number: " + err.Error()
	}
	sla = inputs[10]

	runtimes, err = strconv.Atoi(inputs[11])
	if err != nil {
		return "Run times should be an integer: " + err.Error()
	} else if runtimes < 1 || runtimes > 9 {
		return "Run times should be an integer between 1 to 9"
	}

	return ""
}

func processHTMLFile() {
	var replaced = false

	input, err := ioutil.ReadFile(HTMLTEMPLATE)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, DEFAULTADDR) {
			if len(serverAddr) == 0 {
				log.Fatalln("Please run cnbweb with -server option")
			}
			lines[i] = strings.Replace(line, DEFAULTADDR, serverAddr, 1)
			replaced = true
		}
	}

	if replaced {
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(HTMLTEMPLATE, []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
