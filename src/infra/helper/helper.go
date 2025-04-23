package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	dto "github.com/mrdolles/wallet-app/src/app/dto/user"
)

var jwtKey = []byte("secret_key")

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrUserNotFound        = errors.New("destination user not found")
)

// mempersiapkan Tokenclaims untuk penyimpanan klaim jwt
type TokenClaims struct {
	UserID   int64  `json:"user_id"`
	WalletID int64  `json:"wallet_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// mengGenerate Token untuk membuat token jwt
func GenerateToken(data *dto.RegisterModel) (string, error) {
	expirationTime := time.Now().Add(120 * time.Minute)
	claims := &TokenClaims{
		UserID:   data.ID,
		WalletID: data.WalletID,
		UserName: data.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// verifikasi token JWT
func VerifyToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
