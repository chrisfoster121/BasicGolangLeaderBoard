package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthHelper struct {
	secretKey      []byte
	adminSecretKey []byte
}

type AuthToken struct {
	Token string
}

//The following functions were pulled from :
// https://medium.com/@cheickzida/golang-implementing-jwt-token-authentication-bba9bfd84d60

func CreateAuthHelper(inSecretKey string, inAdminKey string) AuthHelper {
	return AuthHelper{
		secretKey:      []byte(inSecretKey),
		adminSecretKey: []byte(inAdminKey),
	}
}

func (authHelper *AuthHelper) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 30).Unix(),
		})

	tokenString, err := token.SignedString(authHelper.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (authHelper *AuthHelper) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authHelper.secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (authHelper *AuthHelper) VerifyAdminToken(tokenString string) error {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return authHelper.adminSecretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
