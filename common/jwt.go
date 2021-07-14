package common

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"xlk/ginessential/model"
)

var jwtKey = []byte("www.topgoer.com")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) //时间戳，一个小时
	Claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //token的有效期
			IssuedAt:  time.Now().Unix(),     //token的发放时间
			Issuer:    "127.0.0.1",           //token发放者
			Subject:   "user token",          //token的主题
		}, //注意这个逗号
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	//return token.SignedString(jwtKey)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
