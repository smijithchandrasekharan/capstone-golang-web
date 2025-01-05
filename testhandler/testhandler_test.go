package testhandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"capstone-golang-web/src/authentication/login/handler"
	"capstone-golang-web/src/dashboard/TaskHandler"
	"capstone-golang-web/src/CommonHandler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

)

var router *gin.Engine

// Setup function that runs before each test
func setupRouter() {
	router = gin.Default()
	router.LoadHTMLFiles("./src/authentication/login/template/loginForm.html",
		"./src/authentication/login/template/CreateAccount.html",
		"./src/dashboard/home/template/index.html",
		"./src/dashboard/home/template/CreateTask.html",
	)
	router.Static("/css", "./src/authentication/login/template/css")

	// Mock the routes
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
}

func TestPingHandler(t *testing.T) {
	setupRouter()
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "pong"}`, w.Body.String())
}

func TestHealthCheckHandler(t *testing.T) {
	setupRouter()
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "This server is ready to serve requests at PORT#")
}

func TestLoginPageHandler(t *testing.T) {
	setupRouter()
	req, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login to Task Management System")
}

func TestLoginHandler_Success(t *testing.T) {
	// Mock successful login (you can mock db calls using a library like mockery or mock-gorm)
	setupRouter()
	req, _ := http.NewRequest("POST", "/LoginUser", nil)
	req.PostForm = map[string][]string{
		"userName":     {"testuser"},
		"userPassword": {"testpassword"},
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Dashboard")
}

func TestLoginHandler_Failure_WrongPassword(t *testing.T) {
	// Test case where password is incorrect
	setupRouter()
	req, _ := http.NewRequest("POST", "/LoginUser", nil)
	req.PostForm = map[string][]string{
		"userName":     {"testuser"},
		"userPassword": {"wrongpassword"},
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Login failed as password is wrong")
}

func TestCreateUserHandler_Success(t *testing.T) {
	setupRouter()
	reqBody := map[string]string{
		"userName":     "newuser",
		"userPassword": "newpassword",
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

func TestCreateTaskHandler_Success(t *testing.T) {
	// Mock successful task creation
	setupRouter()
	reqBody := map[string]string{
		"title":       "New Task",
		"priority":    "High",
		"description": "Task Description",
		"dueDate":     time.Now().Format(time.RFC3339),
		"status":      "Pending",
		"category":    "Work",
		"project":     "Project A",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/CreateTask", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Created Task Successfully")
}

func TestUpdateTaskHandler_Success(t *testing.T) {
	// Mock task update logic
	setupRouter()
	reqBody := map[string]string{
		"title":       "Updated Task",
		"priority":    "Low",
		"description": "Updated Description",
		"dueDate":     time.Now().Format(time.RFC3339),
		"status":      "In Progress",
		"category":    "Personal",
		"project":     "Project B",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/UpdateTask/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated Task Successfully")
}

func TestViewTaskHandler(t *testing.T) {
	// Mock successful task retrieval
	setupRouter()
	req, _ := http.NewRequest("GET", "/NavigateViewTask/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "View Task")
}

func TestDeleteTaskHandler_Success(t *testing.T) {
	// Mock successful task deletion
	setupRouter()
	req, _ := http.NewRequest("GET", "/Delete/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Removed Task")
}

func TestNavigateHomeHandler(t *testing.T) {
	// Mock home dashboard rendering
	setupRouter()
	req, _ := http.NewRequest("GET", "/Home", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Dashboard")
}
