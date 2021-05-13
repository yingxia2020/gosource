/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2019
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TOTAL = 20
)

func main() {
	fmt.Println("Welcome! Are you ready for the math quiz now (y/n)?")

	reader := bufio.NewReader(os.Stdin)

	answer := getInput(reader)

	if answer != "y" {
		fmt.Println("See you next time!")
		return
	}

	var count, failcount = 0, 0
	var log bytes.Buffer

	start := time.Now()
	rand.Seed(start.UnixNano())

	for count < TOTAL {
		x := rand.Intn(89) + 10
		y := rand.Intn(89) + 10
		fmt.Printf("%d X %d =  ", x, y)
		answer := getInput(reader)
		if answer != strconv.Itoa(x*y) {
			fmt.Println("The answer is WRONG!")
			failcount++
			log.WriteString(fmt.Sprintf("%d X %d =  \n", x, y))
		} else {
			fmt.Println("You are RIGHT and gain one point!")
		}
		count++
	}

	if failcount <= int(TOTAL/20) {
		fmt.Println("The quiz is done. You did a good job! ")
	} else if failcount <= int(TOTAL/10) {
		fmt.Println("The quiz is done. You should practise more and do better next time! ")
	} else {
		fmt.Println("The quiz is done. You failed the quiz!")
	}

	log.WriteString("\n")
	log.WriteString(fmt.Sprintf("Your quiz score is %d\n", count-failcount))
	log.WriteString(fmt.Sprintf("You finished this quiz in %d seconds\n", int(time.Since(start).Seconds())))

	timeNow := time.Now().Format("20060102150405")
	ioutil.WriteFile("quizresult_"+timeNow+".txt", log.Bytes(), 0644)
}

func getInput(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	return strings.TrimSpace(text)
}
