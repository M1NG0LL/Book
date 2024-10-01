package token

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	ID       string `gorm:"primaryKey"`
	IsActive bool   `gorm:"default:false"`
	IsAdmin  bool   `gorm:"default:false"`

	jwt.StandardClaims
}