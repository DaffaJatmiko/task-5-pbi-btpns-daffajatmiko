package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("API_SECRET"))

// ClaimJWT defines the structure for JWT claims.
type ClaimJWT struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token with the given email and username.
func GenerateJWT(email string, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour) // Initialize expiration time
	claims := &ClaimJWT{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Initialize token
	tokenString, err = token.SignedString(jwtKey)             // Generate token string
	return
}

// ValidateToken validates a JWT token.
func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken, // Token string
		&ClaimJWT{},
		func(token *jwt.Token) (interface{}, error) { // Validate token
			return []byte(jwtKey), nil // Return an error if the token is invalid
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*ClaimJWT) // Get claims
	if !ok {
		err = errors.New("couldn't parse claims token") // Return an error if claims are invalid
		return
	}
	if claims.ExpiresAt < time.Now().Unix() { // Return an error if the token is expired
		err = errors.New("token has expired")
		return
	}
	return
}

// GetEmail retrieves the email data from a JWT token.
func GetEmail(signedToken string) (email string, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&ClaimJWT{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*ClaimJWT) // Get claims
	if !ok {
		err = errors.New("couldn't parse claims token")
		return
	}
	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token has expired")
		return
	}

	return claims.Email, nil // Return the email
}
