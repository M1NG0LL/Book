package accounts

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// POST
// Func if you forget the password
func ForgetPass(c *gin.Context) {
	type Info struct {
		Email     string 		`binding:"required"`
	}

	var input Info
	var account Account

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Where("email = ?",input.Email).First(&account).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email doesn't exist"})
		return
	}

	Reset_Code := GenerateCode(6)
	account.Code = Reset_Code

	if err := db.Model(&account).Update("Code", Reset_Code).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset code"})
		return
	}

	passresetLink := fmt.Sprintf("http://localhost:8081/update-password?id=%s&code=%s",account.ID, account.Code)
	
	message := fmt.Sprintf("Welcome to our app!\n\nIf you requested password reset click on the following link: %s", passresetLink)

	if err := SendEmail(account.Email, "Password Reset", message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send activation email"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Please check your email to change your password."})
}

// PUT
// Func to update password from URL in email
func UpdatingPassword(c *gin.Context) {
	accountID := c.Query("id")
	code := c.Query("code")

	type PassReset struct {
		Password string `json:"password" binding:"required"`
	}

	var account Account
	var input PassReset

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := db.Where("id = ? AND code = ?", accountID, code).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found or invalid code"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	account.Password = string(hashedPassword)

	if err := db.Model(&account).Update("password", account.Password).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
