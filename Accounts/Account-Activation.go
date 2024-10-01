package accounts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// POST
// Sending activating email
func ResendActivationLink(c *gin.Context) {
	var input struct {
		Email string `binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var account Account
	if err := db.Where("email = ?", input.Email).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if account.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account is already active"})
		return
	}

	// Resend activation email
	baseURL := "http://localhost:8081"
	if err := SendActivationEmail(&account, baseURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activation email has been resent"})
}

// GET
// Activate account
func ActivateAccount(c *gin.Context) {
	// Get the token from the query parameters
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activation token"})
		return
	}

	var account Account

	if err := db.Where("activation_token = ?", token).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired activation token"})
		return
	}

	if time.Now().After(account.TokenExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Activation token has expired"})
		return
	}

	account.IsActive = true
	account.ActivationToken = ""

	if err := db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account successfully activated"})
}

func SendActivationEmail(account *Account, baseURL string) error {
	activationToken := uuid.New().String()
	tokenExpiresAt := time.Now().Add(24 * time.Hour)

	account.ActivationToken = activationToken
	account.TokenExpiresAt = tokenExpiresAt

	activationLink := fmt.Sprintf("%s/activate?token=%s", baseURL, activationToken)
	message := fmt.Sprintf("Welcome to our app!\n\nPlease activate your account by clicking the following link: %s", activationLink)

	if err := SendEmail(account.Email, "Account Activation", message); err != nil {
		return fmt.Errorf("failed to send activation email: %w", err)
	}

	if err := db.Model(account).Updates(map[string]interface{}{
		"ActivationToken": activationToken,
		"TokenExpiresAt":  tokenExpiresAt,
	}).Error; err != nil {
		return fmt.Errorf("failed to update account with activation token: %w", err)
	}

	return nil
}