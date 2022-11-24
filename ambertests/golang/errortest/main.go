package main

import (
	"errors"
	"fmt"
)

var ERR1 = errors.New("error1")
var ERR2 = errors.New("error2")

func main() {
	err := ERR1

	if err == ERR1 {
		if err := ERR2; err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(err.Error())
	}
}
