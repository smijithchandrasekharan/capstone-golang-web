package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

type TaskItem struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Priority    string    `gorm:"not null"`
	Description string    `gorm:"default:NULL"`
	DueDate     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status      string    `gorm:"status"`
	Category    string    `gorm:"default:NULL"`
	Project     string    `gorm:"default:NULL"`
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

func GetAllTasks(db *gorm.DB) ([]TaskItem, error) {
	var taskItems = []TaskItem{}
	rows, err := db.Model(&TaskItem{}).Rows()

	for rows.Next() {
		var taskItem TaskItem
		// ScanRows scans a row into a struct
		db.ScanRows(rows, &taskItem)
		taskItems = append(taskItems, taskItem)
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

func getTaskByID(db *gorm.DB, taskID int) (*TaskItem, error) {
	var taskItem TaskItem
	result := db.First(&taskItem, taskID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &taskItem, nil
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
	err = db.AutoMigrate(&User{},&TaskItem{})
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
		if result.RowsAffected > 0 {
			fmt.Printf("password %s", (userDB).Password)
			if !CheckPasswordHash(password, (userDB).Password) {
				log.Printf("Login failed for the user %s as password is wrong",user)
					c.HTML(http.StatusOK, "loginForm.html", gin.H{
					"title":   "Login to Task Management System",
					"message": "Login failed",
					
				})
			} else {

				resultTasks, _ := GetAllTasks(db)

				c.HTML(http.StatusOK, "index.html", gin.H{
					"title":      "Dashboard",
					"message":    strings.Join([]string{"Welcome user  ", user}, ""),
					"tasksTable": resultTasks,
				})
				fmt.Println(len(resultTasks))
			}
		}else{
			log.Printf("Login failed for the user %s as user doesn't exist Go ahead and create a new account",user)
			c.HTML(http.StatusOK, "loginForm.html", gin.H{
				"title":   "Login to Task Management System",
				"message": "Login failed",
			})
		}
	})
	router.GET("/signUp", func(c *gin.Context) {
		c.HTML(http.StatusOK, "CreateAccount.html", gin.H{
			"title": "Create Account",
		})
	})

	router.GET("/NavigateCreateTask", func(c *gin.Context) {
		c.HTML(http.StatusOK, "CreateTask.html", gin.H{
			"title":        "Create Task",
			"TaskTitle":    "",
			"Priority":     "",
			"Description":  "",
			"DueDate":      "",
			"Status":       "",
			"Category":     "",
			"Project":      "",
			"TaskAction":   "Create",
			"frmAction":    "/CreateTask",
			"EnableButton": "",
		})
	})

	router.POST("/UpdateTask/:ID", func(c *gin.Context) {
		message := ""
		//taskID,_:= strconv.Atoi(c.Param("ID"))
		taskIDUns, _ := strconv.ParseUint(c.Param("ID"), 10, 32)
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
			log.Printf("Error parsing due date: %s",err.Error())
		}

		updtTask := TaskItem{ID: uint(taskIDUns), Title: title, Description: description, Priority: priority, Project: project, Category: category, DueDate: dateTimeParsed, Status: status}

		result:=db.Save(&updtTask)
		defer func(updtTask *TaskItem, result *gorm.DB) {
			if updtTask.ID > 0 {
				message = "Updated Task Successfully"
			}else{
				if result.Error != nil {
					message = "Update Task Failed"
					log.Printf("Failed to update the task : %s with error %s",updtTask.Title,  result.Error)
				}
			}
			c.HTML(http.StatusOK, "CreateTask.html", gin.H{
				"message":      message,
				"title":        "Update Task",
				"TaskTitle":    updtTask.Title,
				"Priority":     updtTask.Priority,
				"Description":  updtTask.Description,
				"DueDate":      updtTask.DueDate,
				"Status":       updtTask.Status,
				"Category":     updtTask.Category,
				"Project":      updtTask.Project,
				"TaskAction":   "Update",
				"frmAction":    "/UpdateTask/" + c.Param("ID"),
				"EnableButton": "false",
			})
		}(&updtTask,result)
	})

	router.GET("/NavigateViewTask/:ID", func(c *gin.Context) {
		message := ""
		taskId, _ := strconv.Atoi(c.Param("ID"))

		if err != nil {
			log.Printf("Task ID is of invalid format : %s resulted in error %s",c.Param("ID"),err.Error())
		}else{
			taskItm, err := getTaskByID(db, taskId)
			if 	err!=nil{
				log.Printf("Failed to fetch Task with ID  : %s %s",c.Param("ID"),err.Error())
				message="Failed to fetch Task with ID  :"+c.Param("ID")
				c.HTML(http.StatusOK, "CreateTask.html", gin.H{
					"message":message,
					"title":        "View Task",
					"TaskTitle":    "",
					"Priority":     "",
					"Description":  "",
					"DueDate":      "",
					"Status":       "",
					"Category":     "",
					"Project":      "",
					"TaskAction":   "View",
					"frmAction":    "#",
					"EnableButton": "true",
				})
			}else{

		c.HTML(http.StatusOK, "CreateTask.html", gin.H{
			"message":message,
			"title":        "View Task",
			"TaskTitle":    taskItm.Title,
			"Priority":     taskItm.Priority,
			"Description":  taskItm.Description,
			"DueDate":      taskItm.DueDate,
			"Status":       taskItm.Status,
			"Category":     taskItm.Category,
			"Project":      taskItm.Project,
			"TaskAction":   "View",
			"frmAction":    "#",
			"EnableButton": "disabled",
		})
	}}
	})

	router.GET("/Delete/:ID", func(c *gin.Context) {

		message:=""
		taskId, err := strconv.Atoi(c.Param("ID"))

		if err != nil {
			log.Printf("Task ID is of invalid format : %s resulted in error %s",c.Param("ID"),err.Error())
		}else{
			taskItm, err := getTaskByID(db, taskId)
			if 	err!=nil{
				log.Printf("Failed to fetch Task with ID  : %s %s",c.Param("ID"),err.Error())
				message="Failed to fetch Task with ID  :"+c.Param("ID")
			}else{
				result:=db.Delete(taskItm)
				if result.Error != nil {
					message = "Delete Task Failed"
					log.Printf("Failed to delete the task : %s with error %s",c.Param("ID"),  result.Error)
				}else{
					message = "Removed Task "+taskItm.Title
				}
			}
			}
			resultTasks, errFetchTasks := GetAllTasks(db)
			if errFetchTasks!=nil{
				log.Printf("Failed to fetch all tasks in the dashboard : %s ",errFetchTasks.Error())
				message="Failed to fetch all tasks in the dashboard"
			}else{
					if len(resultTasks)==0{
						message="There are no tasks for you in dashboard"
				}}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Dashboard",
			"message":    message,
			"tasksTable": resultTasks,
		})
	})
	router.GET("/Home", func(c *gin.Context) {
		message:=""
		resultTasks, err := GetAllTasks(db)

		if err!=nil{
			log.Printf("Failed to fetch all tasks in the dashboard : %s ",err.Error())
			message="Failed to fetch all tasks in the dashboard"
		}else{

			if len(resultTasks)==0{
				message="There are no tasks for you in dashboard"
			}
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Dashboard",
			"message":    message,
			"tasksTable": resultTasks,
		})
	})

	router.GET("/NavigateEditTask/:ID", func(c *gin.Context) {
		message:=""
		taskId, err := strconv.Atoi(c.Param("ID"))

		if err != nil {
			log.Printf("Task ID is of invalid format : %s resulted in error %s",c.Param("ID"),err.Error())
		}else{
		taskItm, err := getTaskByID(db, taskId)
		if 	err!=nil{
			log.Printf("Failed to fetch Task with ID  : %s %s",c.Param("ID"),err.Error())
			message="Failed to fetch Task with ID  :"+c.Param("ID")
			c.HTML(http.StatusOK, "CreateTask.html", gin.H{
				"message":message
				"title":        "Edit Task",
				"TaskTitle":    "",
				"Priority":     "",
				"Description":  "",
				"DueDate":      "",
				"Status":       "",
				"Category":     "",
				"Project":      "",
				"TaskAction":   "Update",
				"frmAction":    "/UpdateTask/" + c.Param("ID"),
				"EnableButton": "true",
			
			})
		}else{
		c.HTML(http.StatusOK, "CreateTask.html", gin.H{
			"message":message,
			"title":        "Edit Task",
			"TaskTitle":    taskItm.Title,
			"Priority":     taskItm.Priority,
			"Description":  taskItm.Description,
			"DueDate":      taskItm.DueDate,
			"Status":       taskItm.Status,
			"Category":     taskItm.Category,
			"Project":      taskItm.Project,
			"TaskAction":   "Update",
			"frmAction":    "/UpdateTask/" + c.Param("ID"),
			"EnableButton": "false",
		
		})}
	}
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
			log.Printf("Password hashing failed for the user : %s ",user)
		} else {
			result := db.Create(&newUser)
			if result.Error != nil {
				message = "Create User Failed"
				log.Printf("Failed to create the user : %s with error %s",user,  result.Error)
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
			log.Printf("Error parsing due date: %s",err.Error())
		}

		newTask := TaskItem{Title: title, Description: description, Priority: priority, Project: project, Category: category, DueDate: dateTimeParsed, Status: status}

		result:=db.Create(&newTask)
		defer func(newTask *TaskItem) {
			if newTask.ID > 0 {
				message = "Created Task Successfully"
			}else{
				if result.Error != nil {
					message = "Create Task Failed"
					log.Printf("Failed to create the task : %s with error %s",newTask.Title,  result.Error)
				}
			}
			c.HTML(http.StatusOK, "CreateTask.html", gin.H{
				"message":      message,
				"title":        "Create Task",
				"TaskTitle":    newTask.Title,
				"Priority":     newTask.Priority,
				"Description":  newTask.Description,
				"DueDate":      newTask.DueDate,
				"Status":       newTask.Status,
				"Category":     newTask.Category,
				"Project":      newTask.Project,
				"TaskAction":   "Create",
				"frmAction":    "/CreateTask",
				"EnableButton": "false",
			})
		}(&newTask)

	})

	router.Run("localhost:8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
