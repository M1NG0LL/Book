package main

import (
	"log"

	account "project/Accounts"
	book "project/Book"
	token "project/Token"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	router := gin.Default()

	ERR := godotenv.Load()
    if ERR != nil {
        log.Fatalf("Error loading .env file")
    }
	
	var err error
	db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate 
	if err := db.AutoMigrate(&account.Account{}, &book.Book{}, &book.Like{}); err != nil {
		panic("failed to migrate database")
	}	

	// Initialize account and game packages
	account.Init(db)
	book.Init(db)
	token.Init(db)

	protected := router.Group("/")
	protected.Use(token.AuthMiddleware())

	// Account Part ===================================================
	router.POST("/login", token.Login)
	router.POST("/accounts", account.CreateAccount)

	protected.POST("/send-email", account.SendEmailWithCustomizedMessage)
	
		// Account Activation part
	router.POST("/reactivate", account.ResendActivationLink)
	router.GET("/activate", account.ActivateAccount)
	
		// Pass Reset part
	router.POST("/passreset", account.ForgetPass)
	router.PUT("/update-password", account.UpdatingPassword)

	protected.GET("/accounts/me", account.GetMyAccount)
	protected.GET("/accounts", account.GetAccounts)

	protected.PUT("/accounts/:id", account.UpdateAccountByID)
	protected.DELETE("/accounts/:id", account.DeleteAccountbyid)

	// Book Part =======================================================

	protected.POST("/books", book.CreatingBook)
	protected.PUT("/books/:id", book.UpdateMyBook)
	protected.DELETE("/books/:id", book.DeleteBook)

	protected.GET("/books/me", book.GetMyBooks)
	router.GET("/books", book.GetBooks)

		// Like Part
	protected.POST("/books/:id/like", book.PutLike)

	protected.GET("/books/:id/like", book.GetLikesByBook)
	protected.GET("/books/like", book.GetLikesByAccount)

	protected.DELETE("/books/:id/like", book.DeleteLike)

	
	// Run the server
	router.Run(":8081")
}