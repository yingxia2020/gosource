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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	FILESIZE = 100
	INITSIZE = 1000000
	SUFFIX   = "DATA"

	RANDOMTEXT = `bfeWQobe:9oY+vO.DjRV:V@JBZ)Fvj5UqWg?BrppM:'u0/y[cgi9_<,L4I8mQ>\
13fnp|@zH{'R7_VFuZ4M3Om@hMCh@suVsxA(3msi2oVS{kc)hTyc[#tRBZ"isOH
(d%nWAR48*rqFKEQH&pWIgl*DLMpbXah3C]T4[|Rq\@{6w"b<Q'<i\D,||W"ECx
!sCmnp,qeNK^46L1tE8n12"1Z,<gGh-2!R\*y,|<Mcl-01v$U\w*09qW9o>Bn+q
WWw+yW0jp7SEN8HE2Dr.*yv_yFC'P3B'$f6&)8j%G7lqR[Fn0dKw^I'-'CH\gIB
<8:zK0t<U&bzh2E)2T"|fR*:*qPjYy|!u@nq]Ch;r,=ddh?t@-,f,p&|g,@C'W1
q*G7X!pBci_NEfR)7NMqmkj=Ai8K1zgTH/IA9JgesNNXe81>^{)[Yf)[Y>+5;43
9B\':)5p3^{;6R,-vxo9dkH$I2bNI#62,'+,<OgG*)iSJ3/pv^xF2,l2tz!!cx{
:.V$SD?*ePM';x7).p0,7q>S#07{9J6pn0cP3P_-nv__Jyf;Kzz[GhXI8Ci\2?K
4h4|;>!PD'^H?F^2,QgX&k.6SgCJTM>KMnqS!hC?8s+N"O.EB'*2'$I=O|p.qmg
-.p>8DXn,CFP#w^J4?xrx${1@W)xoww9NMwXSbQ:\:^{"{vtVZ9dcV+X/S&N\61`
)

var (
	filedir string
	prefix  string
)

func init() {
	flag.StringVar(&filedir, "dir", "/home/n/temp/download/", "Directory to put download files")
	flag.StringVar(&prefix, "pre", "SUPPORT", "Download file prefix")
}

func main() {
	flag.Parse()

	var randomLen = len(RANDOMTEXT)
	// create the directory is not exist
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		os.Mkdir(filedir, 0777)
	} else {
		os.RemoveAll(filedir)
		os.Mkdir(filedir, 0777)
	}

	// generate files
	for i := 1; i <= FILESIZE; i++ {
		var buffer bytes.Buffer
		var filesize = 0

		if !strings.HasSuffix(filedir, "/") {
			filedir += "/"
		}

		filename := fmt.Sprintf("%s%s%03d.%s", filedir, prefix, i, SUFFIX)
		buffer.WriteString(filename + "\n")
		filesize += len(filename)

		for filesize < INITSIZE*i {
			buffer.WriteString(RANDOMTEXT)
			filesize += randomLen
		}
		ioutil.WriteFile(filename, buffer.Bytes(), 0666)
	}
}
