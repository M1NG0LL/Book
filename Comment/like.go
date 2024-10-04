package comment

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// POST
// Add a like to a comment
func AddLikeToComment(c *gin.Context) {
	accountID, ID_exists := c.Get("accountID")
	if !ID_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	commentID := c.Param("comment_id")

	var existingLike CommentLike
	if err := db.Where("comment_id = ? AND account_id = ?", commentID, accountID).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already liked this comment"})
		return
	}

	newLike := CommentLike{
		CommentID: commentID,
		AccountID: accountID.(string),
		CreatedAt: time.Now(),
	}

	if err := db.Create(&newLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var comment Comment
	if err := db.First(&comment, "id = ?", commentID).Error; err == nil {
		comment.Likes++
		db.Save(&comment)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment liked successfully"})
}

// GET
// Get likes of a comment
func GetLikesOfComment(c *gin.Context) {
	commentID := c.Param("comment_id")

	var likes []CommentLike
	if err := db.Where("comment_id = ?", commentID).Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, likes)
}

// DELETE
// Remove a like from a comment
func RemoveLikeFromComment(c *gin.Context) {
	accountID, ID_exists := c.Get("accountID")
	if !ID_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	commentID := c.Param("comment_id")

	if err := db.Where("comment_id = ? AND account_id = ?", commentID, accountID).Delete(&CommentLike{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var comment Comment
	if err := db.First(&comment, "id = ?", commentID).Error; err == nil {
		if comment.Likes > 0 {
			comment.Likes--
			db.Save(&comment)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}