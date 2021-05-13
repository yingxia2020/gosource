package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var node = "node-a4231-1"
	var cmd = "kubectl describe node " + node + " | grep 'cpu:'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		tokens := strings.Fields(line)
		temp, err := strconv.Atoi(tokens[1])
		if err == nil && temp > 0 {
			fmt.Println(temp)
			break
		}
	}
}
