package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"capstone-golang-web/src/authentication/login/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

)

var router *gin.Engine

// Setup function that runs before each test
func setupRouter() {
	router = gin.Default()
	router.LoadHTMLFiles("../../../../src/authentication/login/template/loginForm.html",
		"../../../../src/authentication/login/template/CreateAccount.html",
		"../../../../src/dashboard/home/template/index.html",
		"../../../../src/dashboard/home/template/CreateTask.html",
	)
	router.Static("/css", "./src/authentication/login/template/css")

	// Mock the routes

	router.POST("/LoginUser", handler.LoginHandler)
	router.POST("/CreateAccount", handler.CreateUserHandler)
}

func TestCreateUserHandler_Success(t *testing.T) {
	setupRouter()
	reqBody := map[string]string{
		"userName":     "newuser2",
		"userPassword": "newpassword2",
		"userEmail":    "user@example.com",
		"userPhone":    "1234567890",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/CreateAccount", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Created User Successfully")
}

