package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/fatih/set"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
)

type policyClaims struct {
	AttestationPolicy string `json:"AttestationPolicy"`
	jwt.RegisteredClaims
}

type policyHeader struct {
	Algorithm string   `json:"alg"`
	Type      string   `json:"typ"`
	X5C       []string `json:"x5c",omitempty`
}

const (
	CERT_HEADER  = "-----BEGIN CERTIFICATE-----"
	CERT_TRAILER = "-----END CERTIFICATE-----"
	PUBLIC_KEY   = "PUBLIC KEY"
	KEY_HEADER   = "x5c"
	NONE_ALG     = "none"
	SIZE_256     = "256"
	SIZE_384     = "384"
	SIZE_512     = "512"
)

var (
	policy      = flag.String("policyfile", "", "Input policy file to be signed")
	privKey     = flag.String("privkeyfile", "", "Input private key file to sign the policy")
	certificate = flag.String("certfile", "", "Input certificate file to verify the policy")
	algorithm   = flag.String("algorithm", "PS384", "Supported algorithm of RSA key pair")
)

func main() {
	flag.Parse()

	// Create permitted algorithm set
	algorithms := set.New(set.NonThreadSafe)
	algorithms.Add("RS256", "PS256", "RS384", "PS384", "RS512", "PS512")

	// Check input parameters
	if len(*policy) == 0 {
		fmt.Println("Input policy file could not be empty")
		return
	}
	policyBytes, err := ioutil.ReadFile(*policy)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !algorithms.Has(*algorithm) {
		fmt.Println("Input algorithm is not supported")
		return
	}

	claims := policyClaims{
		AttestationPolicy: base64.RawURLEncoding.EncodeToString(policyBytes),
	}

	// Assume it is unsigned case
	if len(*privKey) == 0 || len(*certificate) == 0 {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)

		tokenString, err := token.SigningString()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		tokenString = tokenString + "."

		// Output to screen and file
		fmt.Println(tokenString)
		err = ioutil.WriteFile(*policy+".signed", []byte(tokenString), 0664)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Verify the token
		verified, err := verifyToken(tokenString)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Token is verified: ", verified)
		}
		return
	}

	// Otherwise it is signed case
	// Check the input private key and certificate files
	privKeyFinal, certContents, err := checkCertFiles()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Check if provided algorithm makes sense
	signMethod := checkSigningAlgorithm(privKeyFinal)
	if signMethod == nil {
		return
	}

	signedToken := jwt.NewWithClaims(signMethod, claims)
	signedToken.Header[KEY_HEADER] = []string{certContents}
	signedTokenString, err := signedToken.SignedString(privKeyFinal)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Output to screen and file
	fmt.Println(signedTokenString)
	err = ioutil.WriteFile(*policy+".signed", []byte(signedTokenString), 0664)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Verify the token
	valid, err := verifyToken(signedTokenString)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Token is verified: ", valid)
	}
}
