package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"strings"
)

// verifyToken entry method to handle all token formats
func verifyToken(tokenstring string) (bool, error) {
	parts := strings.Split(tokenstring, ".")
	if len(parts) <= 2 {
		return false, errors.New("not a valid jwt format token")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	var header policyHeader
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return false, err
	}
	if header.Algorithm == NONE_ALG {
		//unsigned case
		return verifyUnsignedToken(tokenstring, parts[1])
	}
	if header.X5C != nil && len(header.X5C) > 0 {
		// signed case
		return verifySignedToken(tokenstring, header.X5C[0])
	}
	return false, errors.New("token is not supported")
}

// verifySignedToken example how to verify the signed policy token
func verifySignedToken(tokenstring, certstring string) (bool, error) {
	cert, err := parseCertificate([]byte(addHeaderTrailer(certstring)))
	if err != nil {
		return false, err
	}
	token, err := jwt.ParseWithClaims(tokenstring, &policyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cert.PublicKey.(*rsa.PublicKey), nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(*policyClaims); ok && token.Valid {
		policy, err := base64.RawURLEncoding.DecodeString(claims.AttestationPolicy)
		if err != nil {
			return false, err
		}
		fmt.Println("\nPolicy payload: \n" + string(policy))
		return token.Valid, nil
	}
	return token.Valid, errors.New("failed to decode policy token")
}

// verifyUnsignedToken example how to verify the unsigned policy token
func verifyUnsignedToken(tokenstring, payloadstring string) (bool, error) {
	payload, err := base64.RawURLEncoding.DecodeString(payloadstring)
	if err != nil {
		return false, err
	}

	var claims policyClaims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return false, err
	}
	policy, err := base64.RawURLEncoding.DecodeString(claims.AttestationPolicy)
	if err != nil {
		return false, err
	}
	fmt.Println("\nPolicy payload: \n" + string(policy))

	// Not secure mode but verify with best efforts
	parts := strings.Split(tokenstring, ".")
	method := jwt.GetSigningMethod(NONE_ALG)
	err = method.Verify(strings.Join(parts[0:2], "."), parts[2], jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		return false, err
	}
	return true, nil
}

// publicKeyToBytes public key to bytes
func publicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		fmt.Println(err.Error())
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  PUBLIC_KEY,
		Bytes: pubASN1,
	})

	return pubBytes
}

// removeHeaderTrailer remove header/trailer of cert file
func removeHeaderTrailer(cert []byte) string {
	certString := string(cert)
	var buf bytes.Buffer
	lines := strings.Split(certString, "\n")
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if strings.Contains(line, CERT_HEADER) || strings.Contains(line, CERT_TRAILER) {
			continue
		}
		buf.WriteString(line + "\n")
	}
	return strings.TrimSpace(buf.String())
}

// addHeaderTrailer add header/trailer of cert file
func addHeaderTrailer(cert string) string {
	return CERT_HEADER + "\n" + cert + "\n" + CERT_TRAILER
}

// parseCertificate parse certificate from unencrypted string format
func parseCertificate(certBytes []byte) (*x509.Certificate, error) {
	certBlock, _ := pem.Decode(certBytes)
	return x509.ParseCertificate(certBlock.Bytes)
}

// checkCertFiles check input private key and certificate files are valid
func checkCertFiles() (*rsa.PrivateKey, string, error) {
	privKeyBytes, err := ioutil.ReadFile(*privKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	privKeyFinal, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	certBytes, err := ioutil.ReadFile(*certificate)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}

	cert, err := parseCertificate(certBytes)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}

	pubKeyBytesFromCert := publicKeyToBytes(cert.PublicKey.(*rsa.PublicKey))
	pubKeyBytesFromPriv := publicKeyToBytes(&privKeyFinal.PublicKey)

	if bytes.Compare(pubKeyBytesFromCert, pubKeyBytesFromPriv) != 0 {
		fmt.Println("Provided private key and certificate do not match")
		return nil, "", errors.New("provided private key and certificate do not match")
	}

	certContents := removeHeaderTrailer(certBytes)
	return privKeyFinal, certContents, nil
}

// checkSigningAlgorithm check if provided algorithm makes sense
func checkSigningAlgorithm(privKeyFinal *rsa.PrivateKey) jwt.SigningMethod {
	if privKeyFinal.N.BitLen() == 2048 && !strings.Contains(*algorithm, SIZE_256) {
		fmt.Println("Input private key file and algorithm do not match")
		return nil
	}
	if privKeyFinal.N.BitLen() == 3072 && !strings.Contains(*algorithm, SIZE_384) {
		fmt.Println("Input private key file and algorithm do not match")
		return nil
	}
	if privKeyFinal.N.BitLen() == 4096 && !strings.Contains(*algorithm, SIZE_512) {
		fmt.Println("Input private key file and algorithm do not match")
		return nil
	}
	signMethod := jwt.GetSigningMethod(*algorithm)
	if signMethod == nil {
		fmt.Println("Input signing algorithm not found")
	}
	return signMethod
}
