/*
 *  Copyright (C) 2022 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"fmt"
	"strconv"
	"strings"
)

/*
	input = {"aa", "bb", "cc"}
	output = "["aa", "bb", "cc"]"
*/
func convertStringArray(input []string) string {
	if len(input) == 0 {
		return ""
	}
	var output strings.Builder
	output.WriteString("[")
	for idx, item := range input {
		if idx != len(input)-1 {
			output.WriteString(fmt.Sprintf("%q", item) + ", ")
		} else {
			output.WriteString(fmt.Sprintf("%q", item) + "]\n")
		}
	}
	return output.String()
}

func intToInterfaceSlice(intSlice []int) []interface{} {
	var interfaceSlice []interface{} = make([]interface{}, len(intSlice))
	for i, d := range intSlice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

func stringToInterfaceSlice(stringSlice []string) []interface{} {
	var interfaceSlice []interface{} = make([]interface{}, len(stringSlice))
	for i, d := range stringSlice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

/*
String type:
	input = {"aa", "bb", "cc"}
	output = "["aa", "bb", "cc"]"
Other types:
	input = {1, 2, 3}
	output = "[1, 2, 3]"
*/
func convertArray(input []interface{}) string {
	if len(input) == 0 {
		return ""
	}
	var output strings.Builder
	output.WriteString("[")
	for idx, item := range input {
		if idx != len(input)-1 {
			switch input[0].(type) {
			case string:
				output.WriteString(fmt.Sprintf("%q", item) + ", ")
			default:
				output.WriteString(fmt.Sprintf("%v", item) + ", ")
			}
		} else {
			switch input[0].(type) {
			case string:
				output.WriteString(fmt.Sprintf("%q", item) + "]\n")
			default:
				output.WriteString(fmt.Sprintf("%v", item) + "]\n")
			}
		}
	}
	return output.String()
}

func convertIncludeValue(includeValue interface{}, title string) string {
	var output strings.Builder
	var tmp string
	switch includeValue.(type) {
	case int, bool, float32, float64:
		tmp = fmt.Sprintf(title+"Value := [%v]\n", includeValue)
	case string:
		tmp = fmt.Sprintf(title+"Value := [%q]\n", includeValue)
	default:
		panic("Not supported value format")
	}
	output.WriteString(tmp)
	output.WriteString("includes_value(" + title + "Value, quote." + title + ")\n\n")
	return output.String()
}

func convertIncludeArray(includeValues []interface{}, title string) string {
	var output strings.Builder
	output.WriteString(title + "Values := ")
	output.WriteString(convertArray(includeValues))
	output.WriteString("includes_value(" + title + "Values, quote." + title + ")\n\n")
	return output.String()
}

/* Right now only support string and int type data */
func convertLimitValue(limitField string, title string, stringValue bool) string {
	var output strings.Builder
	if strings.HasPrefix(limitField, "<=") {
		score := strings.TrimSpace(limitField[2:])
		if stringValue {
			output.WriteString(title + "_limit := " + fmt.Sprintf("%q", score) + "\n")
		} else {
			// assume it is int type value
			intScore, _ := strconv.Atoi(score)
			output.WriteString(title + "_limit := " + fmt.Sprintf("%v", intScore) + "\n")
		}
		output.WriteString("smaller_equal_than(quote." + title + ", " + title + "_limit)\n\n")
		usedFuncFlag = usedFuncFlag | SMALLE
	} else if strings.HasPrefix(limitField, ">=") {
		score := strings.TrimSpace(limitField[2:])
		if stringValue {
			output.WriteString(title + "_limit := " + fmt.Sprintf("%q", score) + "\n")
		} else {
			// assume it is int type value
			intScore, _ := strconv.Atoi(score)
			output.WriteString(title + "_limit := " + fmt.Sprintf("%v", intScore) + "\n")
		}
		output.WriteString("bigger_equal_than(quote." + title + ", " + title + "_limit)\n\n")
		usedFuncFlag = usedFuncFlag | BIGE
	} else if strings.HasPrefix(limitField, "<") {
		score := strings.TrimSpace(limitField[1:])
		if stringValue {
			output.WriteString(title + "_limit := " + fmt.Sprintf("%q", score) + "\n")
		} else {
			// assume it is int type value
			intScore, _ := strconv.Atoi(score)
			output.WriteString(title + "_limit := " + fmt.Sprintf("%v", intScore) + "\n")
		}
		output.WriteString("smaller_than(quote." + title + ", " + title + "_limit)\n\n")
		usedFuncFlag = usedFuncFlag | SMALL
	} else if strings.HasPrefix(limitField, ">") {
		score := strings.TrimSpace(limitField[1:])
		if stringValue {
			output.WriteString(title + "_limit := " + fmt.Sprintf("%q", score) + "\n")
		} else {
			// assume it is int type value
			intScore, _ := strconv.Atoi(score)
			output.WriteString(title + "_limit := " + fmt.Sprintf("%v", intScore) + "\n")
		}
		output.WriteString("bigger_than(quote." + title + ", " + title + "_limit)\n\n")
		usedFuncFlag = usedFuncFlag | BIG
	} else {
		score := strings.TrimSpace(limitField)
		if stringValue {
			output.WriteString(title + "_value := [" + fmt.Sprintf("%q", score) + "]\n")
		} else {
			// assume it is int type value
			intScore, _ := strconv.Atoi(score)
			output.WriteString(title + "_value := [" + fmt.Sprintf("%v", intScore) + "]\n")
		}
		output.WriteString("includes_value(" + title + "_value, quote." + title + ")\n\n")
	}
	return output.String()
}
