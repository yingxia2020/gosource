package main

import (
        "github.com/kayac/parallel-benchmark/benchmark"
        "log"
        "time"

        ecies "github.com/ecies/go"
)

func main() {
        result := benchmark.RunFunc(
                func() (subscore int) {
                        encryptDecrypt()
                        return 1
                },
                time.Duration(3) * time.Second,
                10,
        )
        log.Printf("%#v", result)
}

func fib(n int) int {
        if n == 0 {
                return 0
        }
        if n == 1 {
                return 1
        }
        return (fib(n-1) + fib(n-2))
}
func encryptDecrypt() {
        k, _ := ecies.GenerateKey()
        ciphertext, _ := ecies.Encrypt(k.PublicKey, []byte("THIS IS THE TEST"))
        _, _ = ecies.Decrypt(k, ciphertext)
}

