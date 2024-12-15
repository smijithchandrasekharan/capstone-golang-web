package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"default:NULL"`
	Phone     string    `gorm:"default:NULL"`
	DOB       time.Time // Date of birth of the user
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"` // When the user was created in the system
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"` // When the user was last updated in the system
}

func ConnectToPostgreSQL() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 dbname=TaskManagement user=postgres password=password sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := ConnectToPostgreSQL()

	if err != nil {
		log.Fatal(err)
	}

	// Perform database migration
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.LoadHTMLFiles("./src/authentication/login/template/loginForm.html",
		"./src/authentication/login/template/CreateAccount.html",
	"./src/dashboard/home/template/index.html")
	router.Static("/css", "./src/authentication/login/template/css")
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This server is ready to serve requests at PORT#",
		})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "loginForm.html", gin.H{
			"title": "Login to Task Management System",
		})
	})

	router.POST("/LoginUser", func(c *gin.Context) {
		//user := c.PostForm("userName")
		//password := c.PostForm("userPassword")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Dashboard",
		})
	})

	router.GET("/signUp", func(c *gin.Context) {
		c.HTML(http.StatusOK, "CreateAccount.html", gin.H{
			"title": "Create Account",
		})
	})

	router.POST("/CreateAccount", func(c *gin.Context) {

		user := c.PostForm("userName")
		password := c.PostForm("userPassword")
		mail := c.PostForm("userEmail")
		phone := c.PostForm("userPhone")
		newUser := User{Username: user, Password: password, Email: mail, Phone: phone}
		db.Create(&newUser)
		defer func(newUser *User) {
			message := ""
			if newUser.ID>0{
				message="Created User Successfully"
			}
			c.HTML(http.StatusOK, "CreateAccount.html", gin.H{
				"message": message,
				"title" : "Create Account",
			})
		}(&newUser)
	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
