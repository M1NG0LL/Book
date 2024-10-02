package follow

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(database *gorm.DB) {
	db = database
}

// POST
// Follow someone 
func MakeFollow(c *gin.Context) {
	firstID, ID_exists := c.Get("accountID")
	IsActive, _ := c.Get("isActive")

	if !ID_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	} 

	if !IsActive.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
		return
	}

	secondID := c.Param("id")

	if firstID == secondID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	var follow Follow
	err := db.Where("first_id = ? AND second_id = ?", firstID, secondID).First(&follow).Error
	if err == nil {
		var reverseFollow Follow
		err = db.Where("first_id = ? AND second_id = ?", secondID, firstID).First(&reverseFollow).Error
		if err == nil {
			follow.Relationship = "friends"
			db.Save(&follow)
			reverseFollow.Relationship = "friends"
			db.Save(&reverseFollow)

			c.JSON(http.StatusOK, gin.H{"message": "Relationship upgraded to friends"})
			return
		}

		c.JSON(http.StatusConflict, gin.H{"error": "Already following this person"})
		return
	}

	newFollow := Follow{
		FirstID:      firstID.(string),
		SecondID:     secondID,
		Relationship: "follower",
		CreatedAt: time.Now(),
	}

	if err := db.Create(&newFollow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

// DELETE
// Unfollow Someone
func UnFollow(c *gin.Context) {
	firstID, ID_exists := c.Get("accountID")
	IsActive, _ := c.Get("isActive")

	if !ID_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	} 

	if !IsActive.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
		return
	}

	secondID := c.Param("id")

	if firstID == secondID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}        

	var follow Follow
	if err := db.Where("first_id = ? AND second_id = ?", firstID, secondID).First(&follow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow relationship not found"})
		return
	}

	if err := db.Delete(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

// GET
// Get Followers
func GetFollowers(c *gin.Context) {
	accountID, ID_exists := c.Get("accountID")
	IsActive, _ := c.Get("isActive")

	if !ID_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	} 

	if !IsActive.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
		return
	}

	var followers []Follow
	if err := db.Where("second_id = ?", accountID).Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, followers)
}

// GET
// Get Number of followers 
func GetNumberOfFollowersAndFriends(c *gin.Context) {
	accountID, ID_exists := c.Get("accountID")
	IsActive, _ := c.Get("isActive")

	if !ID_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	} 

	if !IsActive.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
		return
	} 

	var counts = map[string]int64{}
	relationships := []string{"follower", "friends", "besties", "lovers"}

	for _, relation := range relationships {
		var count int64
		db.Model(&Follow{}).Where("second_id = ? AND relationship = ?", accountID, relation).Count(&count)
		counts[relation] = count
	}

	c.JSON(http.StatusOK, counts)
}