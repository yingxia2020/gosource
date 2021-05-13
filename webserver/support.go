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
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
)

func supportSim(buffer *bytes.Buffer, action string, result []string) {
	switch action {
	case "1":
		categoryList(buffer)
	case "2":
		productList(buffer)
	case "3":
		searchProduct(buffer)
	case "4":
		downloadCategories(buffer)
	case "5":
		languages(buffer)
	case "6":
		opSys(buffer)
	case "7":
		fileCatalog(buffer, result[4])
	case "8":
		fileInfo(buffer, result[4])
	default:
		errorCase(buffer)
	}
}

func categoryList(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")

	prefix := rand4Digits()
	for i := minCurr; i <= maxCurr; i++ {
		buffer.WriteString(fmt.Sprintf("%04d&%s\n", i+prefix, getLongname()))
	}
	buffer.WriteString("</pre>\n")
}

func productList(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")

	prefix := rand4Digits()
	for i := minCurr; i <= maxCurr; i++ {
		buffer.WriteString(fmt.Sprintf("%04d&%s\n", i+prefix, getProduct()))
	}
	buffer.WriteString("</pre>\n")
}

func searchProduct(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")

	prefix := rand4Digits()
	for i := minCurr; i <= SCALEDLOAD; i++ {
		buffer.WriteString(fmt.Sprintf("%04d&%s\n", i+prefix, getProduct()))
	}
	buffer.WriteString("</pre>\n")
}

func downloadCategories(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")

	for i := 0; i < currSize; i++ {
		buffer.WriteString(fmt.Sprintf("%s\n", DnldCatList[i]))
	}
	buffer.WriteString("</pre>\n")
}

func languages(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")

	for i := 0; i < lanSize; i++ {
		buffer.WriteString(fmt.Sprintf("%s\n", LanguageList[i]))
	}
	buffer.WriteString("</pre>\n")
}

func opSys(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")

	for i := 0; i < osSize; i++ {
		buffer.WriteString(fmt.Sprintf("%s\n", OSList[i]))
	}
	buffer.WriteString("</pre>\n")
}

func fileCatalog(buffer *bytes.Buffer, param string) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")

	key := rand.Intn(100)
	nprodid, _ := strconv.Atoi(param)
	for i := 0; i <= lanSize; i++ {
		buffer.WriteString(fmt.Sprintf("%06d&%s&2004-6-1 10:15am&%d&%s\n",
			i+(nprodid*10), getFilename(), (key+i)*nprodid,
			getFiledesc()))
	}
	buffer.WriteString("</pre>\n")
}

func fileInfo(buffer *bytes.Buffer, param string) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")

	nprodid, _ := strconv.Atoi(param)
	key := rand.Intn(100)
	buffer.WriteString(fmt.Sprintf("%06d&%s&2004-6-1 10:15am&%d&%s\n%s\n%s\n",
		(nprodid * 10), getFilename(), key*nprodid, getURL(),
		getFiledesc(), getAdditionalInfo()))

	buffer.WriteString("</pre>\n")
}

func rand4Digits() int {
	return (rand.Intn(9) + 1) * 1000
}
