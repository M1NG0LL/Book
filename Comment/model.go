package comment

import "time"

type Comment struct {
	ID          string    `gorm:"primaryKey"`
	AccountID   string    
	AccountName string    

	CommentTo   string     
	Description string     
	CommentDate time.Time `gorm:"autoCreateTime"`

	Likes       int       
}

type CommentLike struct {
	CommentID  string
	AccountID  string
	CreatedAt  time.Time
}