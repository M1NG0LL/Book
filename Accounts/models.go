package accounts

import (
	"time"
)

type Account struct {
	ID		  string 			`gorm:"primaryKey"`
	Username  string 			`gorm:"unique;not null"`
	ProfilePhoto     string    	`gorm:"type:varchar(255)"`
	Email     string 			`gorm:"unique;not null"`
	Password  string			`gorm:"not null"`

	Code 	  string			`gorm:"default:' '"`

	ActivationToken string    	`json:"activation_token"`
	TokenExpiresAt  time.Time 	`json:"token_expires_at"`

	IsActive  bool 				`gorm:"default:false"`

	IsAdmin   bool 				`gorm:"default:false"`
}