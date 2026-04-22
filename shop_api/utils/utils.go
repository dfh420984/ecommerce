package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint64, expire time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func MD5(str string) string {
	data := []byte(str)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func SHA256(str string) string {
	data := []byte(str)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func GenerateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("%s%s%s%s%s",
		now.Format("2006010215"),
		randString(4, "0123456789"),
		randString(4, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	)
}

func GenerateTradeNo() string {
	return fmt.Sprintf("%s%s",
		time.Now().Format("20060102150405"),
		randString(8, "0123456789"),
	)
}

func randString(length int, chars string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func InSlice(v string, slice []string) bool {
	for _, s := range slice {
		if v == s {
			return true
		}
	}
	return false
}
