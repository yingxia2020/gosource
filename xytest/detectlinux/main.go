package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("hostnamectl").Output()
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(strings.ToLower(string(out)), "centos") {
		fmt.Println("This is CentOS")
	} else if strings.Contains(strings.ToLower(string(out)), "ubuntu") {
		fmt.Println("This is Ubuntu")
	} else {
		fmt.Println("This is unknow OS")
	}
}
