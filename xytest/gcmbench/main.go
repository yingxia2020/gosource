package main

import (
        "crypto/aes"
        "crypto/cipher"
        "flag"
        "fmt"
	"runtime"
        "time"
)

var benchIterations = flag.Int("i", 1, "Benchmark iterations")
var ptLen = flag.Int("b", 8192, "Plaintext length, in bytes")
var verbose = flag.Int("v", 1, "Verbose output")

func gcm() (float64, float64) {
        var key [16]byte
        var nonce [12]byte
        var ad [13]byte
        aes, _ := aes.NewCipher(key[:])
        aesgcm, _ := cipher.NewGCM(aes)
        var out []byte
        var buf = make([]byte, *ptLen)
        cipherIterations := *benchIterations * 10000
        tStart := time.Now()
        for i := 0; i < cipherIterations; i++ {
                out = aesgcm.Seal(out[:0], nonce[:], buf, ad[:])
        }
        t := time.Now().Sub(tStart).Seconds()
        bw := (float64(*ptLen*cipherIterations) / t) / float64(1024.0*1024.0)
        return bw, t
}

func main() {
	runtime.GOMAXPROCS(1)
        flag.Parse()
        bw, t := gcm()
        if *verbose == 1 {
                fmt.Printf("GCMBench completed, throughput is %4.2f MB/s (%2.4f s)\n", bw, t)
        }
}

