package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

const PRIV_KEY = `XXXXXXXX`

func main() {
	//mySigningKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(PRIV_KEY))
	//mySigningKey := []byte("AllYourBase")
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode([]byte(PRIV_KEY)); block == nil {
		fmt.Println("Decode failed")
		return
	}
	mySigningKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}

	// Create the claims
	claims := MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	// Create claims while leaving out some of the optional fields
	claims = MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS384, claims)
	ss, err := token.SignedString(mySigningKey)
	switch key := mySigningKey.(type) {
	case *rsa.PrivateKey:
		bitLen := key.N.BitLen()
		fmt.Println(bitLen)
	}
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(ss)
	}
	//fmt.Printf("%v %v", ss, err)
}
