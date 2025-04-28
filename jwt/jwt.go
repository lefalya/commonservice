package jwt

import (
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/lefalya/commonuser"
	"os"
	"time"
)

type Claims struct {
	commonuser.Account
	jwt.RegisteredClaims
}

func JWTEncode(issuer string, duration time.Duration, subject string, data commonuser.Account) (string, error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	claims.Sub = data.Sub
	claims.UUID = data.UUID
	claims.Name = data.Name
	claims.Email = data.Email
	claims.Username = data.Username
	claims.CreatedAt = data.CreatedAt
	claims.UpdatedAt = data.UpdatedAt
	claims.AssociatedAccount = data.AssociatedAccount

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func JWTDecode(token string) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to decode token claims")
	}

	claims.UUID = claims.Sub

	return claims, nil
}

func JWTVerify(token string) error {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok && !parsedToken.Valid {
		return errors.New("invalid token")
	}

	return nil
}
