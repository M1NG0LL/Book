package accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(database *gorm.DB) {
	db = database
}

// POST
// Creating Account
func CreateAccount(c *gin.Context) {
	var account Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidEmail(account.Email) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
        return
    }

	var existingAccount Account
	if err := db.Where("username = ? OR email = ?", account.Username, account.Email).First(&existingAccount).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or Email already exists"})
		return
	}

	if err, error := ValidatePassword(account.Password); err {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.MinCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	account.Password = string(hashedPassword)

	account.ID = uuid.New().String()
	
	account.IsActive = false
	account.IsAdmin = false

	baseURL := "http://localhost:8081"
	if err := SendActivationEmail(&account, baseURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Account created. Please check your email to activate your account."})
}

// POST
// Sending Email with a customized message
func SendEmailWithCustomizedMessage(c *gin.Context) {
	IsActive, _ := c.Get("isActive")
	isAdmin, Admin_exists := c.Get("isAdmin")

	if !Admin_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if !isAdmin.(bool) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This url is for ADMIN ONLY."})
		return
	}

	if !IsActive.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
		return
	}

	type Info struct {
		Email 	string	`binding:"required"`
		Subject	string	`binding:"required"`
		message	string 	`binding:"required"`
	}

	var input Info
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var account Account
	if err := db.Where("email = ?", input.Email).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := SendEmail(input.Email,input.Subject,input.message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email has been sent"})
} 

// GET
// get Account info using the Token
func GetMyAccount(c *gin.Context) {
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

	var account Account
	if err := db.First(&account, "id = ?", accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       account.ID,
		"username": account.Username,
		"profile_photo": account.ProfilePhoto,
		"email":    account.Email,
		
		"is_active": account.IsActive,
		"is_admin": account.IsAdmin,
	})
}

// GET
// ADMIN's FUNCTION
// Func to Get All Accounts
func GetAccounts(c *gin.Context) {
	IsActive, _ := c.Get("isActive")
	isAdmin, Admin_exists := c.Get("isAdmin")

	if !isAdmin.(bool) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This url is for ADMIN ONLY."})
		return
	}

	if !IsActive.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
		return
	}

	if !Admin_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var accounts []Account
	if err := db.Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// PUT
// This func is for ADMINS ONLY
// Update any account by putting id in url 
func UpdateAccountByID(c *gin.Context) {
	isAdmin, Admin_exists := c.Get("isAdmin")
	IsActive, _ := c.Get("isActive")

	if  !Admin_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if !IsActive.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
		return
	}
	
	if !isAdmin.(bool) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This url is for ADMIN ONLY."})
		return
	} 

	paramID := c.Param("id")
	
	var preaccount Account
	if err := db.First(&preaccount, "id = ?", paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	
	var account Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if account.Password != preaccount.Password {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		account.Password = string(hashedPassword)
	}

	if err := db.Model(&Account{}).Where("id = ?", paramID).Updates(account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account updated successfully"})
}

// DELETE
// This func is for ADMINS ONLY
// Delete any account by putting id in url 
func DeleteAccountbyid(c *gin.Context)  {
	isAdmin, Admin_exists := c.Get("isAdmin")

	if !Admin_exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	paramID := c.Param("id")

	if !isAdmin.(bool) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This url is for ADMIN ONLY."})
		return
	}

	var account Account
	if err := db.First(&account, "id = ?", paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := db.Delete(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}