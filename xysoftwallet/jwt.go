package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)
/*
func main() {
	// This token is expired
	var tokenString = "eyJraWQiOiIxMDEiLCJhbGciOiJFUzM4NCJ9.eyJuY2UiOiJhR2gxYXdsbGZCVEFrOWJwNEpkajVoM0Z0Vm5Jb0JsdV9FM1kwNUtVSmdVIiwiYXVkIjoiREVWIiwic2NwIjoiREVWLUFQSSBJTkRTIiwic3ViIjoiY25idGVzdDIwMjMrMUBnbWFpbC5jb20iLCJ2ZXIiOiIwOTEwNTU1NiIsIm5iZiI6MTY4NjMzMzM5OSwiaXNzIjoiaWRlbnRpdHkuaW50ZWwuY29tIiwidHlwIjoiREVWIiwiZXhwIjoxNjg2MzM3MDA0LCJpYXQiOjE2ODYzMzMzOTksImp0aSI6IjE2NjcyMjkwNTE2MDcwNjQ1OTAifQ.dl3xsKB-JxzhYSAQoTLhXiWHUmcxcqlSfy799LBZ-o-eU7Tb6Iw_qpjl-jyvZpqlhf9QkfC9Jz4cg3sfp_dyTVpGp0w6d8YDNk9GcGod10sYLvFRc6rzxhCswQxFcasu"
	fmt.Println(getHashedUsernameFromJwtToken(tokenString))

	// expected results:
	// username is:  cnbtest2023+1@gmail.com
	// 81d6a0ba793fa01ffd04f55e1bab20c117c2027eac2d35a3d5a845a153c873b1
}
*/
func getHashedUsernameFromJwtToken(tokenString string) string {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("failed to get JWT token claims")
		return ""
	}

	userName := claims["sub"].(string)
	if len(userName) == 0 {
		fmt.Println("could not find sub field in JWT token claims")
		return ""
	} else {
		fmt.Println("username is: ", userName)
	}

	hash := sha256.Sum256([]byte(userName))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}