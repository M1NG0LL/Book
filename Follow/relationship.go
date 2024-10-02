package follow

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChangeRelationship(c *gin.Context) {
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
	newRelationship := c.Query("relationship")

	validRelationships := map[string]bool{"friends": true, "besties": true, "lovers": true}
	if !validRelationships[newRelationship] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid relationship"})
		return
	}

	var follow Follow
	if err := db.Where("first_id = ? AND second_id = ?", firstID, secondID).First(&follow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow relationship not found"})
		return
	}

	follow.Relationship = newRelationship
	if err := db.Save(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relationship updated successfully"})
}