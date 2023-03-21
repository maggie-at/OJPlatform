package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var myKey = []byte("OJPlatform Key")

// GetMd5 生成md5
func GetMd5(s string) string {
	// %x: 转化为16进制
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GenerateToken 生成token
func GenerateToken(identity string, name string, isAdmin int) (string, error) {
	userClaim := UserClaims{
		Identity: identity,
		Name:     name,
		IsAdmin:  isAdmin,
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &userClaim)
	tokenStr, err := claims.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenStr, err
}

// AnalyseToken 解析token
func AnalyseToken(token string) (*UserClaims, error) {
	var userClaim UserClaims
	claims, err := jwt.ParseWithClaims(token, &userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims.Valid == false {
		return nil, fmt.Errorf("token analyse error")
	}
	return &userClaim, nil
}
