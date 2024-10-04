package main

import (
	"log"

	account "project/Accounts"
	book "project/Book"
	comment "project/Comment"
	follow "project/Follow"
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
	if err := db.AutoMigrate(
		&account.Account{},
		 &book.Book{}, &book.Like{},
		  &follow.Follow{},
		   &comment.Comment{}, &comment.CommentLike{},
		   ); err != nil {
		panic("failed to migrate database")
	}	

	// Initialize account and game packages
	account.Init(db)
	book.Init(db)
	token.Init(db)
	follow.Init(db)
	comment.Init(db)

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

		// Follow part =============================================
	protected.POST("accounts/follow/:id", follow.MakeFollow)
	protected.DELETE("accounts/follow/:id", follow.UnFollow)
	protected.GET("accounts/follow", follow.GetFollowers)
	protected.GET("accounts/follow/num", follow.GetNumberOfFollowersAndFriends)
	protected.PUT("accounts/follow/relationship/:id", follow.ChangeRelationship)

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

		// Comment part =============================================
	protected.POST("/books/:id/comments", comment.MakeComment)
	protected.DELETE("/books/:id/comments", comment.DeleteComment)
	router.GET("/books/comments/:id", comment.GetCommentsByCommentTo)
	protected.GET("/books/comments/me", comment.GetCommentsByAccount)

			// Comment Likes Part
	protected.POST("/books/comments/:id/like", comment.AddLikeToComment)
	protected.GET("/books/comments/:id/like", comment.GetLikesOfComment)
	protected.DELETE("/books/comments/:id/like", comment.RemoveLikeFromComment)
	
	
	// Run the server
	router.Run(":8081")
}