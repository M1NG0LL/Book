package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST
// Put like on a book by ID
func PutLike(c *gin.Context) {
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
	
	bookID := c.Param("id")

	var book Book
	if err := db.First(&book, "id = ?", bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var like Like
	if err := db.Where("account_id = ? AND book_id = ?", accountID, bookID).First(&like).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Account already liked this book"})
		return
	}

	newLike := Like{
		AccountID: accountID.(string),
		BookID:    bookID,
	}

	if err := db.Create(&newLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	book.Likes++
	if err := db.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Liked successfully"})
}

// GET
// Get All Likes of an account
func GetLikesByAccount(c *gin.Context) {
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

	var likes []Like
	if err := db.Where("account_id = ?", accountID.(string)).Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var likedBooks []Book

	for _, like := range likes {
		var book Book
		if err := db.First(&book, "id = ?", like.BookID).Error; err == nil {
			likedBooks = append(likedBooks, book) 
		}
	}

	c.JSON(http.StatusOK, likedBooks)
}

// GET
// Get Likes of a book
func GetLikesByBook(c *gin.Context) {
	bookID := c.Param("id") 

	var likes []Like
	if err := db.Where("book_id = ?", bookID).Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var accountIDs []string

	for _, like := range likes {
		accountIDs = append(accountIDs, like.AccountID)
	}

	c.JSON(http.StatusOK, accountIDs)
}


// DELETE
// Removes a like from a book
func DeleteLike(c *gin.Context) {
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
	
	bookID := c.Param("id")

	var like Like
	if err := db.Where("account_id = ? AND book_id = ?", accountID, bookID).First(&like).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Delete(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete like"})
		return
	}

	var book Book
	if err := db.First(&book, bookID).Error; err == nil {
		book.Likes--
		db.Save(&book) 
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}
