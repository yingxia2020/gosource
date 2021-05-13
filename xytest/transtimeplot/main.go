package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var podMap = make(map[string]int)
	var outputBuffer bytes.Buffer

	// initialize the map
	fp, err := os.Open("podnames.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		tokens := strings.Fields(scanner.Text())
		podMap[strings.TrimSpace(tokens[0])] = 0
	}

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		sum := 0.0
		var podName string
		tokens := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		for _, token := range tokens {
			if strings.Contains(token, "mc-service") {
				podName = token
				podMap[podName]++
			} else {
				f, _ := strconv.ParseFloat(token, 64)
				outputBuffer.WriteString(fmt.Sprintf("%.2f", f) + " ")
				sum += f
			}
		}
		outputBuffer.WriteString(fmt.Sprintf("%.2f", sum) + " " + podName + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for k, v := range podMap {
		fmt.Println(k+": ", v)
	}
	ioutil.WriteFile("output.txt", outputBuffer.Bytes(), 0666)
}
