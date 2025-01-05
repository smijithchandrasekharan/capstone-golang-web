package handler


import (
	"log"
	"net/http"
	"capstone-golang-web/src/models"
	"capstone-golang-web/src/CommonHandler"

	"github.com/gin-gonic/gin"

)

func SignUpHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "CreateAccount.html", gin.H{
		"title": "Create Account",
	})
}

func CreateUserHandler(c *gin.Context) {
	message := ""
	user := c.PostForm("userName")
	password := c.PostForm("userPassword")
	hashed, err := CommonHandler.HashPassword(password)
	mail := c.PostForm("userEmail")
	phone := c.PostForm("userPhone")
	newUser := models.User{Username: user, Password: hashed, Email: mail, Phone: phone}
	if err != nil {
		message = "Password hashing failed"
		log.Printf("Password hashing failed for the user : %s ", user)
	} else {
		result := CommonHandler.DbSQLConnection.Create(&newUser)
		if result.Error != nil {
			message = "Create User Failed"
			log.Printf("Failed to create the user : %s with error %s", user, result.Error)
		}
	}

	defer func(newUser *models.User) {
		if newUser.ID > 0 {
			message = "Created User Successfully"
			c.HTML(http.StatusOK, "CreateAccount.html", gin.H{
				"message": message,
				"title":   "Create Account",
			})
		} else {
			c.HTML(http.StatusInternalServerError, "CreateAccount.html", gin.H{
				"message": message,
				"title":   "Create Account",
			})
		}
	}(&newUser)

}