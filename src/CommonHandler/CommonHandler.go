package CommonHandler

import (
	"capstone-golang-web/src/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DbSQLConnection *gorm.DB

func ConnectToPostgreSQL() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 dbname=TaskManagement user=postgres password=password sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "This server is ready to serve requests at PORT#",
	})
}

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func init() {
	db, err := ConnectToPostgreSQL()
	DbSQLConnection = db

	if err != nil {
		log.Fatal(err)
	}

	// Perform database migration
	err = db.AutoMigrate(&models.User{}, &models.TaskItem{})
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllTasks() ([]models.TaskItem, error) {
	var taskItems = []models.TaskItem{}
	rows, err := DbSQLConnection.Model(&models.TaskItem{}).Rows()

	for rows.Next() {
		var taskItem models.TaskItem
		// ScanRows scans a row into a struct
		DbSQLConnection.ScanRows(rows, &taskItem)
		taskItems = append(taskItems, taskItem)
		fmt.Println(taskItem.Title)
		// Perform operations on each user
	}
	return taskItems, err

}
