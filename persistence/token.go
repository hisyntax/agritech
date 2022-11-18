package persistence

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type SignedDetails struct {
	Uid string
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_kEY")

// GenerateAllTokens generates both teh detailed token and refresh token
func GenerateAllTokens(uid string) (signedToken string, signedRefreshToken string, err error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
	claims := &SignedDetails{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println(err)
		return "", "", err

	}

	return token, refreshToken, err
}

func GenerateAdminTokens(uid string) (signedToken string, err error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
	claims := &SignedDetails{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return token, err
}

// ValidateToken validates the jwt token
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "invalid token"
		return
	}
	// claims.Uid
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		msg = err.Error()
		return
	}

	return claims, msg
}
