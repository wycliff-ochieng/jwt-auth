package internal

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	//"github.com/google/uuid"
)

type Claims struct {
	ID    int64
	Email string
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refresjtoken"`
}

var JwtSecret = []byte("mydogsnameisrufus")
var RefreshSecret = []byte("myotherdogiscalledbuckeye")

func GenerateToken(ID int64, email string, JwtSsecret string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		ID:    ID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(token)
	return token.SignedString([]byte(JwtSecret))
}

//validating the token

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token:%v", err)
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func GenerateTokenPair(ID int64, email, JwtSecret, RefreshSecret string, jwtExpiry, refreshExpiry time.Duration) (*TokenPair, error) {

	accesstoken, err := GenerateToken(ID, email, JwtSecret, jwtExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Access Token %v", err)
	}

	refreshtoken, err := GenerateToken(ID, email, RefreshSecret, refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Refresh Token %v", err)
	}

	return &TokenPair{
		AccessToken:  accesstoken,
		RefreshToken: refreshtoken,
	}, nil
}

func ExtractBearerToken(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer" {
		return authHeader[7:]
	}
	return ""
}
