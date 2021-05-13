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
)

var article = [5]string{"the", "a", "one", "some", "this"}

var Article = [5]string{"The", "A", "One", "Some", "This"}

var adjective = [10]string{"new", "newest", "best", "brightest", "10x",
	"fast", "slimline", "lowcost", "highspeed", "reliable"}

var Adjective = [10]string{"1200MHz", "36GB", "100Mbs", "8x", "4.0GHz",
	"1GB", "Red", "White", "Blue", "Open-box"}

var adverb = [10]string{"fast", "quickly", "now", "today", "always",
	"here", "too", "fine", "great", "really"}

var CatGroups = [5]string{"Personal", "Office", "Corporate", "Academic", "Value"}

var Currencies = [20]string{"USD", "CAD", "EUR", "JPY", "AUD",
	"MXN", "ARS", "BMD", "BOB", "CLP",
	"EGP", "HKD", "INR", "JMD", "NOK",
	"RUR", "SSL", "THB", "VEB", "YUM"}

var directVerb = [20]string{"click", "download", "Select", "double-click", "extract",
	"copy", "load", "save", "remove", "send",
	"send", "test", "press", "set", "close",
	"power-on", "power-off", "shutdown", "get", "install"}

var DirectVerb = [20]string{"Click", "Download", "Select", "Double-Click", "Extract",
	"Copy", "Load", "Save", "Remove", "Send",
	"Send", "Test", "Press", "Set", "Close",
	"Power-On", "Power-Off", "Shutdown", "Get", "Install"}

var DnldCatList = [20]string{"All", "Application", "Audio Drivers", "BIOS", "Chipset",
	"Communication Drivers", "Diagnostics", "IDE/SCSI", "Input Drivers",
	"Keyboard  Drivers",
	"Monitors", "Network Drivers", "Patches", "Removable Media Drivers",
	"Security Patches",
	"Software Dev. Tools", "System Utilities", "System Management",
	"Video Drivers", "Virus Protection"}

var LanguageList = [30]string{"Arabic", "Bulgarian", "Chinese-S", "Chinese-T", "Czech",
	"Danish", "Dutch", "English", "Estonian", "Finnish",
	"French", "German", "Greek", "Hebrew", "Hungarian",
	"Indonesian", "Italian", "Japanese", "Korean", "Norwegian",
	"Pan-Euro", "Polish", "Portuguese", "Russian", "Slovak",
	"Slovenian", "Spanish", "Swedish", "Thai", "Turkish"}

var ModelType = [10]string{"PRO", "XL-", "ES-", "SD-", "IS-",
	"SCH", "UX-", "VAL", "EZ-", "RX-"}

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

var OSList = [10]string{"RNZX", "RN3K", "RN2045", "Enterprise Xilin V3.01", "Desktop Xilin V3.0",
	"OS 743 for Architecture N7", "NP-OS V13.41", "CafeOS 2.3.1 for HA82",
	"CafeOS 2.3.1 for NA90", "FreeBinOS 7.5"}

var preposition = [10]string{"to", "from", "with", "like", "on",
	"by", "before", "through", "under", "over"}

var Software = [10]string{"RNZX_BIOS", "RN3K_DIAGNOSTIC", "RN2045_Utilities",
	"Xilin_V3.01", "Driver_Xilin_V3.0",
	"Patch_Set_OS743N7", "NP-OS_V13.41", "CafeOS2.3.1_HA82_Drivers",
	"CafeOS2.3.1_NA90_SecurityPatch", "FreeBinOS7.2to7.3_Updater"}

var Suffix = [10]string{"exe", "bin", "zip", "class", "doc",
	"tgz", "txt", "sh", "kit", "upd"}

var SuffixDesc = [10]string{"executable binary", "self-installing binary", "compressed file",
	"JVM binary", "document file",
	"compressed tar file", "ascii text file", "shell script",
	"distibution", "auto-updater"}

var Superlative = [10]string{"Best", "Leading", "State-of-the-Art", "High-Performance", "Stylish",
	"Cool", "Industry-Standard", "Award-Winning", "New", "Value-Priced"}

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
	lanSize  = 30
	osSize   = 10

	SCALEDLOAD = 15

	FEATURELEN    = 31
	LONGNAMELEN   = 31
	PRODUCTLEN    = 51
	HIGHLIGHTSLEN = 127
	OVERVIEWLEN   = 255
	ADDINFO       = 2047
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

func getLongname() string {
	temp := fmt.Sprintf("%s %s", Adjective[rand.Intn(10)], Nouns[rand.Intn(20)])
	if len(temp) > LONGNAMELEN {
		return temp[0:LONGNAMELEN]
	}
	return temp
}

func getProduct() string {
	temp := fmt.Sprintf("%s %s%04d", Nouns[rand.Intn(20)], ModelType[rand.Intn(5)],
		rand.Intn(1000))
	if len(temp) > PRODUCTLEN {
		return temp[0:PRODUCTLEN]
	}
	return temp
}

func getFilename() string {
	temp := fmt.Sprintf("%s.%s", Software[rand.Intn(10)],
		Suffix[rand.Intn(10)])
	if len(temp) > PRODUCTLEN {
		return temp[0:PRODUCTLEN]
	}
	return temp
}

func getFiledesc() string {
	i := rand.Intn(5)
	temp := fmt.Sprintf("This is the %s for the %s%04d %s %s using %s",
		SuffixDesc[i], ModelType[i], i*10, CatGroups[i], Nouns[rand.Intn(25)],
		Software[i])
	if len(temp) > OVERVIEWLEN {
		return temp[0:OVERVIEWLEN]
	}
	return temp
}

func getAdditionalInfo() string {
	var buffer bytes.Buffer
	textlength := ADDINFO/2 + (ADDINFO/20)*(rand.Intn(10)) + (ADDINFO/20)*(rand.Intn(2))
	rlen := 0
	lnctr := 0
	slen := 0

	for rlen < textlength {
		if (textlength - rlen) > 120 {
			lnctr++
			temp := fmt.Sprintf("%d. %s %s %s %s %s %s %s %s %s.<BR>", lnctr,
				DirectVerb[rand.Intn(20)], article[rand.Intn(5)], adjective[rand.Intn(10)],
				noun[rand.Intn(20)], directVerb[rand.Intn(20)], preposition[rand.Intn(10)],
				article[rand.Intn(5)], adjective[rand.Intn(10)], noun[rand.Intn(20)])
			slen = len(temp)
			buffer.WriteString(temp)
		} else if (textlength - rlen) > 80 {
			temp := fmt.Sprintf("	%s %s %s %s %s.",
				Article[rand.Intn(5)], adjective[rand.Intn(10)], noun[rand.Intn(20)],
				directVerb[rand.Intn(20)], adverb[rand.Intn(10)])
			slen = len(temp)
			buffer.WriteString(temp)
		} else if (textlength - rlen) > 40 {
			temp := fmt.Sprintf(" %s %s %s %s.",
				Article[rand.Intn(5)], noun[rand.Intn(20)], verb[rand.Intn(10)],
				adverb[rand.Intn(10)])
			slen = len(temp)
			buffer.WriteString(temp)
		} else if (textlength - rlen) > 20 {
			slen = len(" Please Reboot Now.")
			buffer.WriteString(" Please Reboot Now.")
		} else {
			slen++
			buffer.WriteString(".")
		}
		rlen += slen
	}

	temp := buffer.String()
	if len(temp) > ADDINFO {
		return temp[0:ADDINFO]
	}
	return temp
}

// TODO, point to where file server files are stored, it has the format like:
// http://file-server-ip-address:8080/support###.data, ### from 0-99
func getURL() string {
	temp := fmt.Sprintf("http://143.183.198.59:8080/support%03d.data",
		rand.Intn(100))
	return temp
}

func errorCase(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("1\n")
	buffer.WriteString("Unknown action type!\n")
	buffer.WriteString("</pre>\n")
}
