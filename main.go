package main

import (
	"capstone-golang-web/src/authentication/login/handler"
	"capstone-golang-web/src/CommonHandler"
	"capstone-golang-web/src/dashboard/TaskHandler"

	"github.com/gin-gonic/gin"

)

func main() {

	router := gin.Default()

	router.LoadHTMLFiles("./src/authentication/login/template/loginForm.html",
		"./src/authentication/login/template/CreateAccount.html",
		"./src/dashboard/home/template/index.html",
		"./src/dashboard/home/template/CreateTask.html",
	)
	router.Static("/css", "./src/authentication/login/template/css")

	router.GET("/ping", CommonHandler.PingHandler)
	router.GET("/health", CommonHandler.HealthCheckHandler)
	router.GET("/login", handler.LoginPageHandler)
	router.POST("/LoginUser", handler.LoginHandler)
	router.GET("/signUp", handler.SignUpHandler)
	router.GET("/NavigateCreateTask", TaskHandler.NavigateCreateTaskHandler)
	router.POST("/UpdateTask/:ID", TaskHandler.UpdateTaskHandler)
	router.GET("/NavigateViewTask/:ID", TaskHandler.ViewTaskHandler)
	router.GET("/Delete/:ID", TaskHandler.DeleteTaskHandler)
	router.GET("/Home", TaskHandler.NavigateHomeHandler)
	router.GET("/NavigateEditTask/:ID", TaskHandler.NavigateEditTask)
	router.POST("/CreateAccount", handler.CreateUserHandler)
	router.POST("/CreateTask", TaskHandler.CreateTaskHandler)

	router.Run("localhost:8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
