package book

import "time"

type Book struct {
	ID           string `gorm:"primaryKey"`
	ProfilePhoto string `gorm:"type:varchar(255)"`
	Name         string `gorm:"unique;not null"`

	AuthorID   string
	AuthorName string

	UploadTime time.Time

	Likes	int

	PDF string `gorm:"type:varchar(255)"`
}

type Like struct {
	AccountID string `gorm:"primaryKey"` 
	BookID    string `gorm:"primaryKey"` 
}