package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

type TaskItem struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Priority    string    `gorm:"not null"`
	Description string    `gorm:"default:NULL"`
	DueDate     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status      string    `gorm:"status"`
	Category    string    `gorm:"default:NULL"`
	Project    string    `gorm:"default:NULL"`
}

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

func GetAllTasks(db *gorm.DB )([]TaskItem, error) {
	var taskItems = []TaskItem {}
	rows, err := db.Model(&TaskItem{}).Rows()
	
	for rows.Next() {
		var taskItem TaskItem
	// ScanRows scans a row into a struct
		db.ScanRows(rows, &taskItem)
		taskItems = append(taskItems,taskItem)
		fmt.Println(taskItem.Title)
	// Perform operations on each user
	}
	return taskItems, err
			
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByID(db *gorm.DB, Username string, Password string) (*User, error) {
	var user User
	result := db.First(&user, Username, Password)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
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
	err = db.AutoMigrate(&TaskItem{})
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.LoadHTMLFiles("./src/authentication/login/template/loginForm.html",
		"./src/authentication/login/template/CreateAccount.html",
		"./src/dashboard/home/template/index.html",
		"./src/dashboard/home/template/CreateTask.html",
	)
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
			"title":   "Login to Task Management System",
			"message": "",
		})
	})

	router.POST("/LoginUser", func(c *gin.Context) {
		//var isPasswordMatch bool
		user := c.PostForm("userName")
		password := c.PostForm("userPassword")

		fmt.Printf("user %s password %s", user, password)
		var userDB User
		result := db.First(&userDB, "username = ?", user)
		if result.RowsAffected >= 0 {
			fmt.Printf("password %s", (userDB).Password)
			if !CheckPasswordHash(password, (userDB).Password) {
				c.HTML(http.StatusOK, "loginForm.html", gin.H{
					"title":   "Login to Task Management System",
					"message": "Login failed",
				})
			} else {
				
				resultTasks,_ :=  GetAllTasks(db)
				replacedString:=""
				if len(resultTasks)>0 {
					replacedString+="<table>"
				htmlRowTemplate := "<td><a href='#'>ID</a></td><td>Title</td><td>Description</td><td>Priority</td><td>Project</td><td>Category</td><td>Status</td><td>DueDate</td></td>"
				for _,tsk:=range(resultTasks){
					replacedString+="tr"
					replacedString += strings.ReplaceAll(htmlRowTemplate,"ID","1")
					replacedString = strings.ReplaceAll(replacedString,"Title",tsk.Title)
					replacedString = strings.ReplaceAll(replacedString,"Description",tsk.Description)
					replacedString = strings.ReplaceAll(replacedString,"Priority",tsk.Priority)
					replacedString = strings.ReplaceAll(replacedString,"Project",tsk.Project)
					replacedString = strings.ReplaceAll(replacedString,"Category",tsk.Category)
					replacedString = strings.ReplaceAll(replacedString,"DueDate",tsk.DueDate.GoString())
					replacedString = strings.ReplaceAll(replacedString,"Status",tsk.Status)
					replacedString+="</tr>"
				}
				replacedString+="</table>"
				}
				
				c.HTML(http.StatusOK, "index.html", gin.H{
					"title":   "Dashboard",
					"message": strings.Join([]string{"Welcome user  ", user}, ""),
					"tasksTable":resultTasks , 
				})
				fmt.Println(len(resultTasks))
			}
		}
	})
	router.GET("/signUp", func(c *gin.Context) {
		c.HTML(http.StatusOK, "CreateAccount.html", gin.H{
			"title": "Create Account",
		})
	})
	router.GET("/NavigateCreateTask", func(c *gin.Context) {
		c.HTML(http.StatusOK, "CreateTask.html", gin.H{
			"title": "Create Task",
		})
	})

	router.POST("/CreateAccount", func(c *gin.Context) {
		message := ""
		user := c.PostForm("userName")
		password := c.PostForm("userPassword")
		hashed, err := HashPassword(password)
		mail := c.PostForm("userEmail")
		phone := c.PostForm("userPhone")
		newUser := User{Username: user, Password: hashed, Email: mail, Phone: phone}
		if err != nil {
			message = "Password hashing failed"
		} else {
			result := db.Create(&newUser)
			if result.Error != nil {
				message = "Create User Failed"
			}
		}
		

		defer func(newUser *User) {
			if newUser.ID > 0 {
				message = "Created User Successfully"
			}
			c.HTML(http.StatusOK, "CreateAccount.html", gin.H{
				"message": message,
				"title":   "Create Account",
			})
		}(&newUser)

	})

	router.POST("/CreateTask", func(c *gin.Context) {
		message := ""
		title := c.PostForm("title")
		description := c.PostForm("description")
		priority := c.PostForm("priority")
		project := c.PostForm("project")
		category := c.PostForm("category")
		dueDate := c.PostForm("dueDate")
		status := c.PostForm("status")
		
		// Define the layout of the date string
		layout := time.RFC3339

		// Parse the date string into a time.Time object
		dateTimeParsed, err := time.Parse(layout, dueDate)
		if err != nil {
			fmt.Println("Error parsing date:", err)
		}
		
		newTask := TaskItem{Title: title, Description: description, Priority: priority, Project: project,Category:category,DueDate: dateTimeParsed,Status:  status }
		
		db.Create(&newTask)
		defer func(newTask *TaskItem) {
			if newTask.ID > 0 {
				message = "Created Task Successfully"
			}
			c.HTML(http.StatusOK, "CreateTask.html", gin.H{
				"message": message,
				"title":   "Create Task",
			})
		}(&newTask)

	})

	router.Run("localhost:8085") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
