package internal

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	ID    uuid.UUID
	Email string
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refresjtoken"`
}

var secret = []byte("mydogsnameisrufus")

func generateToken(ID uuid.UUID, email string, secret string, expiry time.Duration) (string, error) {
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
	return token.SignedString([]byte(secret))
}

//validating the token

func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token:%v", err)
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func generateTokenPair(ID uuid.UUID, email, jwtSecret, refreshSecret string, jwtExpiry, refreshExpiry time.Duration) (*TokenPair, error) {

	fmt.Printf("Generating tokens for user:\n")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("JWT Secret: %s\n", jwtSecret)
	fmt.Printf("JWT Expiry: %v\n", jwtExpiry)
	fmt.Printf("Refresh Secret: %s\n", refreshSecret)
	fmt.Printf("Refresh Expiry: %v\n", refreshExpiry)

	accesstoken, err := generateToken(ID, email, jwtSecret, jwtExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Access Token %v", err)
	}

	refreshtoken, err := generateToken(ID, email, refreshSecret, refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Refresh Token %v", err)
	}

	return &TokenPair{
		AccessToken:  accesstoken,
		RefreshToken: refreshtoken,
	}, nil
}
