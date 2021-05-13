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
	"math/rand"
)

var noun = [20]string{"computer", "memory", "processor", "printer", "cartridge",
	"disk", "cable", "camera", "PDA", "DVD",
	"box", "CD-ROM", "RAM", "paper", "software",
	"price", "LCD", "monitor", "program", "projector"}

var Nouns = [25]string{"Computers", "DIMMs", "Processors", "Printers", "Cartridges",
	"Disks", "Networks", "Cameras", "PDAs", "DVDs",
	"Boxes", "CD-ROMs", "Laptops", "Accessories", "Software",
	"Servers", "LCD-TVs", "Monitors", "Programs", "Projectors",
	"Desktops", "Operating Systems", "Video Centers",
	"Digital Cameras", "MP3 Players"}

var Currencies = [20]string{"USD", "CAD", "EUR", "JPY", "AUD",
	"MXN", "ARS", "BMD", "BOB", "CLP",
	"EGP", "HKD", "INR", "JMD", "NOK",
	"RUR", "SSL", "THB", "VEB", "YUM"}

var Article = [5]string{"The", "A", "One", "Some", "This"}

var Adjective = [10]string{"1200MHz", "36GB", "100Mbs", "8x", "4.0GHz",
	"1GB", "Red", "White", "Blue", "Open-box"}

var adjective = [10]string{"new", "newest", "best", "brightest", "10x",
	"fast", "slimline", "lowcost", "highspeed", "reliable"}

var Superlative = [10]string{"Best", "Leading", "State-of-the-Art", "High-Performance", "Stylish",
	"Cool", "Industry-Standard", "Award-Winning", "New", "Value-Priced"}

var adverb = [10]string{"fast", "quickly", "now", "today", "always",
	"here", "too", "fine", "great", "really"}

var preposition = [10]string{"to", "from", "with", "like", "on",
	"by", "before", "through", "under", "over"}

var verb = [20]string{"is", "has", "runs", "includes", "meets",
	"works", "sets", "makes", "looks", "builds",
	"drives", "sends", "changes", "receives", "takes",
	"upgrades", "applies", "delivers", "features", "adjusts"}

var Verb = [20]string{"Is", "Has", "Runs", "Includes", "Meets",
	"Works", "Sets", "Makes", "Looks", "Builds",
	"Drives", "Sends", "Changes", "Receives", "Takes",
	"Upgrades", "Applies", "Delivers", "Features", "Adjusts"}

const (
	minCurr  = 1
	maxCurr  = 75
	currSize = 20

	SCALEDLOAD    = 15
	FEATURELEN    = 31
	HIGHLIGHTSLEN = 127
	OVERVIEWLEN   = 255
)

func getOverview() string {
	temp := fmt.Sprintf("%s %s %s %s %s %s", Article[rand.Intn(5)],
		adjective[rand.Intn(10)], adjective[rand.Intn(10)], noun[rand.Intn(20)],
		verb[rand.Intn(20)], adverb[rand.Intn(10)])
	if len(temp) > OVERVIEWLEN {
		return temp[0:OVERVIEWLEN]
	}
	return temp
}

func getHighlights() string {
	temp := fmt.Sprintf("%s %s %s",
		Verb[rand.Intn(20)], preposition[rand.Intn(10)], adjective[rand.Intn(10)])
	if len(temp) > HIGHLIGHTSLEN {
		return temp[0:HIGHLIGHTSLEN]
	}
	return temp
}

func getFeature() string {
	temp := fmt.Sprintf("%s %s %s",
		Superlative[rand.Intn(10)], Adjective[rand.Intn(10)], Nouns[rand.Intn(20)])
	if len(temp) > FEATURELEN {
		return temp[0:FEATURELEN]
	}
	return temp
}
