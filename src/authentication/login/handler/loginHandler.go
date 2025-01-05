package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"capstone-golang-web/src/models"
	"capstone-golang-web/src/CommonHandler"

	"github.com/gin-gonic/gin"

)

func LoginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "loginForm.html", gin.H{
		"title":   "Login to Task Management System",
		"message": "",
	})
}

func LoginHandler(c *gin.Context) {
	//var isPasswordMatch bool
	user := c.PostForm("userName")
	password := c.PostForm("userPassword")

	fmt.Printf("user %s password %s", user, password)
	var userDB models.User
	result := CommonHandler.DbSQLConnection.First(&userDB, "username = ?", user)
	if result.RowsAffected > 0 {
		fmt.Printf("password %s", (userDB).Password)
		if !CommonHandler.CheckPasswordHash(password, (userDB).Password) {
			log.Printf("Login failed for the user %s as password is wrong", user)
			c.HTML(http.StatusUnauthorized, "loginForm.html", gin.H{
				"title":   "Login to Task Management System",
				"message": "Login failed as password is wrong",
			})
		} else {

			resultTasks, _ := CommonHandler.GetAllTasks()

			c.HTML(http.StatusOK, "index.html", gin.H{
				"title":      "Dashboard",
				"message":    strings.Join([]string{"Welcome user  ", user}, ""),
				"tasksTable": resultTasks,
			})
			fmt.Println(len(resultTasks))
		}
	} else {
		log.Printf("Login failed for the user %s as user doesn't exist Go ahead and create a new account", user)
		c.HTML(http.StatusUnauthorized, "loginForm.html", gin.H{
			"title":   "Login to Task Management System",
			"message": "Login failed as user doesn't exist",
		})
	}
}
