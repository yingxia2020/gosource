package main

import (
	"encoding/json"
        "fmt"
	"os"
)

type config struct {
	Nodes []node `json:"nodes"`
}

type node struct {
	Addr string `json:"ip_address"`
	Hostname string `json:"hostname"`
}

func main() {
	var nodes []node
	args := os.Args[1:]

	for i:=0; i<len(args); i++ {
		node := node{args[i], ""}
		nodes = append(nodes, node)	
	}

	conf := &config{
		Nodes:  nodes}

	result, _ := json.Marshal(conf)
	fmt.Println(string(result))
}	
