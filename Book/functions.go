package book

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
// Creating New Book
func CreatingBook(c *gin.Context) {
	accountID, ID_exists := c.Get("accountID")
	IsActive, _ := c.Get("isActive")

	// Validations

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

	if !account.Author {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not Author"})
		return
	}

	// Creating Book
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = uuid.New().String()
	book.AuthorID = account.ID
	book.AuthorName = account.Username
	book.Likes = 0
	book.UploadTime = time.Now()

	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book Created."})
}

// GET
// Get all books of an author
func GetMyBooks(c *gin.Context) {
	accountID, ID_exists := c.Get("accountID")
	IsActive, _ := c.Get("isActive")

	// Validations

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

	if !account.Author {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not Author"})
		return
	}

	var books []Book

	if err := db.Where("author_id = ?", accountID).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Account": account, "Books": books})
}

// GET
// Get All books
func GetBooks(c *gin.Context) {
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// PUT
// Update a book by ID
func UpdateMyBook(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	
	var book Book
	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var updatedBook Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.Name = updatedBook.Name
	book.ProfilePhoto = updatedBook.ProfilePhoto
	book.PDF = updatedBook.PDF

	if err := db.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// DELETE
// Delete a Book
func DeleteBook(c *gin.Context) {
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
	
	if !account.Author {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not Author"})
		return
	}

	bookID := c.Param("id")

	var book Book
	if err := db.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if book.AuthorID != account.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "This Book isn't yours"})
	}

	if err := db.Where("book_id = ?", bookID).Delete(&Like{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete likes"})
		return
	}

	if err := db.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}