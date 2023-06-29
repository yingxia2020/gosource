package main


import (
	"github.com/skythen/bertlv"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

/*
func main() {
	var privateKeyInHex = "e9d1f05b6b6a4632095dc7d5e5808b7394257b207bd64426b630369791eaec49"
	//var publicKey = "0407ec0e6002002d7aab16f0a3eb3f728baba5aaf383f51ae370710f10565f54e0f49e292cca8dc15b90800da5ad4b2a21303e2046a3f35763bb9aba1c9d03001ca50ba28e04a59f45036789dac541bcbe97a479cedf6fcbd8c5dbe1d3bfd4b727"
	//var subjectIdentifier = "260c059d1561788036cb637662"
	//// expected output
	//var expected = "3f2181c05f200d260c059d1561788036cb6376627f4966b0610407ec0e6002002d7aab16f0a3eb3f728baba5aaf383f51ae370710f10565f54e0f49e292cca8dc15b90800da5ad4b2a21303e2046a3f35763bb9aba1c9d03001ca50ba28e04a59f45036789dac541bcbe97a479cedf6fcbd8c5dbe1d3bfd4b727f00100950200805f3740b6aec05160d67de55ada5377db18f2da61df16e3fafa1d5b279d5968df382d0dab7c46c92547572137a4bfad0438590790b9b2bacc3622788e678669dc37a3cb"

	privatekey, err := getPrivateKeyFromHex("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err == nil {
		fmt.Println("OK")
	}
	// Convert private key to raw bytes
	privateKeyBytes := privatekey.D.Bytes()

	// Convert raw bytes to hex string
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	if privateKeyHex == privateKeyInHex {
		fmt.Println("Matched Private Key (hex):", privateKeyHex)
	}
}
*/

func generateCsr(privatekey string, publickey string, subject string) {
	builder := bertlv.Builder{}
	val := []byte(publickey)
	val1 := []byte("00")
	pubkeyBytes := builder.AddBytes(bertlv.NewOneByteTag(0x10), val).AddBytes(bertlv.NewOneByteTag(0x50), val1).Bytes()

	builder1 := bertlv.Builder{}
	dataToBeSigned := builder1.AddBytes(bertlv.NewTwoByteTag(0x5F, 0x20), []byte(subject)).AddBytes(bertlv.NewTwoByteTag(0x7F, 0x49), pubkeyBytes).AddBytes(bertlv.NewOneByteTag(0x95),[]byte("0080")).Bytes()

	fmt.Println(hex.EncodeToString(dataToBeSigned))
}

func generate384KeyPair() (string, string){
	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating ECDSA private key:", err)
		return "", ""
	}

	// Convert private key to raw bytes
	privateKeyBytes := privateKey.D.Bytes()

	// Convert raw bytes to hex string
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	fmt.Println("Private Key (hex):", privateKeyHex)

	// Get the raw bytes of the public key
	publicKeyBytes := elliptic.Marshal(privateKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)

	// Convert raw bytes to hex string
	publicKeyHex := hex.EncodeToString(publicKeyBytes)

	fmt.Println("Public Key (hex):", publicKeyHex)

	return privateKeyHex, publicKeyHex
}

func getPrivateKeyFromHex(privateKeyHex string) (*ecdsa.PrivateKey, error){
	// Decode hex string to bytes
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		fmt.Println("Error decoding private key hex:", err)
		return nil, err
	}

	// Generate private key from bytes
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = elliptic.P384() // Replace with the appropriate curve if needed
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)

	// Validate the private key
	if !privateKey.Curve.IsOnCurve(privateKey.PublicKey.X, privateKey.PublicKey.Y) {
		fmt.Println("Invalid private key")
		return nil, err
	}

	return privateKey, nil
}