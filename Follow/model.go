package follow

import "time"

type Follow struct {
	FirstID      string `gorm:"primaryKey"` 
	SecondID     string `gorm:"primaryKey"`
	Relationship string 

	CreatedAt time.Time 
}