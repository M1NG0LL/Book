package comment

import (
	"net/http"
	Account "project/Accounts"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(database *gorm.DB) {
	db = database
}

// POST
// Make comment on something(comment to)
func MakeComment(c *gin.Context) {
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
	
	var account Account.Account
	if err := db.First(&account, "id = ?", accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	commentTo := c.Param("id")

	var input struct {
        Description string `json:"description" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

	newComment := Comment{
        ID:          uuid.New().String(),
        AccountID:   account.ID,
        AccountName: account.Username,
        CommentTo:   commentTo,
        Description: input.Description,
        CommentDate: time.Now(),
    }

	if err := db.Create(&newComment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully", "comment": newComment})
}

// DELETE
// Delete comment if you were admin or the one who wrote it
func DeleteComment(c *gin.Context) {
	accountID, ID_exists := c.Get("accountID")
	IsActive, _ := c.Get("isActive")
	isAdmin, _ := c.Get("isAdmin")

	if !ID_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	} 

	if !IsActive.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
		return
	}

    commentID := c.Param("id")

    var comment Comment
    if err := db.Where("id = ?", commentID).First(&comment).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }

	if accountID == comment.AccountID || isAdmin.(bool) {
		if err := db.Delete(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't Delete Comment"})
        return
	}
}

// GET
// Get All comments on something
func GetCommentsByCommentTo(c *gin.Context) {
    commentTo := c.Param("id")

    var comments []Comment
    if err := db.Where("comment_to = ?", commentTo).Find(&comments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, comments)
}

// GET
// Get All comments of account
func GetCommentsByAccount(c *gin.Context) {
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

    var comments []Comment
    if err := db.Where("account_id = ?", accountID).Find(&comments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, comments)
}