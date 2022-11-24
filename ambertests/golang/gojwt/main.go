package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

const PRIV_KEY = `-----BEGIN PRIVATE KEY-----
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
S9ZKf/2eU0EP7VrnXK0+Nz+B
-----END PRIVATE KEY-----`

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
