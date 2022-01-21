package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// IJWTService interface
type IJWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)
}

// AuthCustomClaims ...
type AuthCustomClaims struct {
	UserName string `json:"username"`
	User     bool   `json:"user"`
	jwt.StandardClaims
}

// JWTService ...
type jwtService struct {
	secretKey string
	issuer    string
}

func getSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

// GenerateToken ...
func (service *jwtService) GenerateToken(email string, isUser bool) string {
	claims := &AuthCustomClaims{
		email,
		isUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

// ValidateToken ...
func (service *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {

	return jwt.ParseWithClaims(
		encodedToken,
		&AuthCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {

			if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
				return nil, fmt.Errorf("Invalid token %s", token.Header["alg"])

			}

			return []byte(service.secretKey), nil
		})

	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {

		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token %s", token.Header["alg"])

		}

		return []byte(service.secretKey), nil
	})
}

// JWTService returns pointer to a jwtService
func JWTService() *jwtService {
	return &jwtService{
		issuer:    "FitApp",
		secretKey: getSecretKey(),
	}
}
