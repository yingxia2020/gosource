package main

import (
	"encoding/base64"
	"syscall/js"
)

var (
	passPhase = "abcdefghijklmnop"
)

func encryptoWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		inputStr := args[0].String()
		//fmt.Printf("input %s\n", inputStr)

		encoded := encryptData([]byte(inputStr), passPhase)
		//fmt.Println(base64.StdEncoding.EncodeToString(encoded))
		return base64.StdEncoding.EncodeToString(encoded)
	})
	return jsonFunc
}

func decryptoWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		inputStr := args[0].String()
		//fmt.Printf("input %s\n", inputStr)

		tmp, _ := base64.StdEncoding.DecodeString(inputStr)
		decoded := decryptData(tmp, passPhase)
		//fmt.Println(decoded)
		return string(decoded)
	})
	return jsonFunc
}

func main() {

	js.Global().Set("encrypt", encryptoWrapper())
	js.Global().Set("decrypt", decryptoWrapper())
	<-make(chan bool)
}
