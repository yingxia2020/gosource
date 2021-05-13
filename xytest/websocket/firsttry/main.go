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
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
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
	nodeNum    int
	cpuNode    int
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
	maxPod     string
)

const (
	HTMLTEMPLATE = "websockets.html"
	DEFAULTADDR  = "localhost"
	CPUSTR       = "000m"
	CONFIGFILE   = "config.json"
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
	CPU     string `json:"cpu"`
	MaxPod  string `json:"maxpodnumber"`
}

type Config struct {
	RunOption   string     `json:"runoption"`
	RunTimes    int        `json:"runtimes"`
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

			// Print the message to the console
			// fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			errors := precheckInput(string(msg))
			if len(errors) > 0 {
				// Write errors back to browser
				if err = conn.WriteMessage(msgType, []byte(errors)); err != nil {
					log.Fatalln(err)
				}
				continue
			}

			writeToJSON()

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				log.Fatalln(err)
			}
		}
	})

	// Need preprocess HTML file
	processHTMLFile()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, HTMLTEMPLATE)
	})

	http.ListenAndServe(":8080", nil)
}

func writeToJSON() {
	data := Config{
		RunOption:   option,
		RunTimes:    1, // TODO
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
			MaxPod:  maxPod,
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

	nodeNum, err = strconv.Atoi(inputs[0])
	if err != nil {
		return "Number of Nodes should be integer: " + err.Error()
	}

	cpuNode, err = strconv.Atoi(inputs[1])
	if err != nil {
		return "CPU per Node should be integer: " + err.Error()
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

	cpudiv, err := strconv.Atoi(inputs[5])
	if err != nil {
		return "CPU per POD should be integer: " + err.Error()
	}

	cpuPod = inputs[5] + CPUSTR

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

	if eclient <= sclient {
		return "Last client number should be bigger than initial client number"
	}

	iclients = inputs[6]
	clientStep = inputs[7]
	lclients = inputs[8]

	interval, err := strconv.Atoi(inputs[9])
	if err != nil || interval < 1 {
		return "Time interval should be integer greater or equal to 1: " + err.Error()
	}
	time = inputs[9]

	_, err = strconv.ParseFloat(inputs[10], 32)
	if err != nil {
		return "SLA should be a number: " + err.Error()
	}
	sla = inputs[10]

	var temp = int((nodeNum*cpuNode*9)/(cpudiv*10)) - nodeNum
	if temp < 1 {
		return "Not enough CPU to create OCR POD"
	}

	maxPod = strconv.Itoa(temp)

	return ""
}

func processHTMLFile() {
	if len(serverAddr) == 0 {
		return
	}

	input, err := ioutil.ReadFile(HTMLTEMPLATE)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, DEFAULTADDR) {
			lines[i] = strings.Replace(line, DEFAULTADDR, serverAddr, 1)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(HTMLTEMPLATE, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
