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
	//"bytes"
	//"crypto/sha256"
	//"encoding/hex"
	"fmt"
	"log"
	//"math"
	"os/exec"
	//"time"

	"github.com/valyala/fasthttp"
)

/*
const (
	FILESIZE = 5000000
	ROUND    = 100

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
-.p>8DXn,CFP#w^J4?xrx${1@W)xoww9NMwXSbQ:\:^{"{vtVZ9dcV+X/S&N\61
`
)
*/
func main() {
	s := &fasthttp.Server{
		Handler:     handler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	if err := s.ListenAndServe(":8072"); err != nil {
		log.Fatalf("Error in ListenAndServe calc server: %s", err)
	}
}

/*
func handler(ctx *fasthttp.RequestCtx) {
	var buffer bytes.Buffer

	var blockSize = len(RANDOMTEXT)
	var totalSize = blockSize
	buffer.WriteString(RANDOMTEXT)

	for totalSize < FILESIZE {
		totalSize += blockSize
		buffer.WriteString(RANDOMTEXT)
	}
	fmt.Println("Buffer size is ", buffer.Len())

	startTime := time.Now()
	var sha256Bytes [32]byte
	for i := 0; i < ROUND; i++ {
		sha256Bytes = sha256.Sum256(buffer.Bytes())
	}
	fmt.Println(hex.EncodeToString(sha256Bytes[:]))
	fmt.Println(time.Since(startTime).Nanoseconds())

	fmt.Fprintf(ctx, "SHA256 finish successfully\n[%s]\n\n", hex.EncodeToString(sha256Bytes[:]))
}
*/
/*Use 5 times original SQRT calculations to scale up on powerful machines
func handler(ctx *fasthttp.RequestCtx) {
	var x = 0.0001

	for i := 0; i <= 5000000; i++ {
		x += math.Sqrt(x)
	}

	fmt.Fprintf(ctx, "result=%f OK\n\n", x)
}
*/

func handler(ctx *fasthttp.RequestCtx) {
	out, err := exec.Command("./workload.sh").CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(ctx, string(out))
}
