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
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Plot struct {
	Legend string
	Data   string
}

type Config struct {
	Title  string
	Output string
	Plots  []Plot
}

func main() {
	viper.SetConfigFile("plot.json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	conf := &Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("Unable to decode into config struct, %s", err.Error())
	}

	fmt.Println(conf.Title)
	fmt.Println(conf.Plots[1].Legend)
	fmt.Println(len(conf.Plots))
}
