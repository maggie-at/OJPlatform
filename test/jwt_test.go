package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

type UserClaim struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var myKey = []byte("OJPlatform Key")

// TestGenerateToken 测试生成token
func TestGenerateToken(t *testing.T) {
	userClaim := UserClaim{
		Identity:       "alan25",
		Name:           "alan",
		StandardClaims: jwt.StandardClaims{},
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &userClaim)
	tokenStr, err := claims.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenStr)
}

// TestParseToken 测试解析token
func TestParseToken(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXIgMSIsIm5hbWUiOiJVU0VSMTExIn0.48cllcr67LNu6F3H63b_aUFBFnigNPNGKlohK35pcmg"

	var userClaim UserClaim
	claims, err := jwt.ParseWithClaims(tokenStr, &userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid == true {
		fmt.Println(userClaim)
	}
}
