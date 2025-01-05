package TaskHandler

import (
	"capstone-golang-web/src/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"capstone-golang-web/src/CommonHandler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

)

func getTaskByID(taskID int) (*models.TaskItem, error) {
	var taskItem models.TaskItem
	result := CommonHandler.DbSQLConnection.First(&taskItem, taskID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &taskItem, nil
}

func NavigateCreateTaskHandler(c *gin.Context) {
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
		"EnableButton": "false",
	})
}

func UpdateTaskHandler(c *gin.Context) {
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
		log.Printf("Error parsing due date: %s", err.Error())
	}

	updtTask := models.TaskItem{ID: uint(taskIDUns), Title: title, Description: description, Priority: priority, Project: project, Category: category, DueDate: dateTimeParsed, Status: status}

	result := CommonHandler.DbSQLConnection.Save(&updtTask)
	defer func(updtTask *models.TaskItem, result *gorm.DB) {
		if updtTask.ID > 0 {
			message = "Updated Task Successfully"
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
		} else {
			if result.Error != nil {
				message = "Update Task Failed"
				log.Printf("Failed to update the task : %s with error %s", updtTask.Title, result.Error)
				c.HTML(http.StatusInternalServerError, "CreateTask.html", gin.H{
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
			}
		}

	}(&updtTask, result)
}

func ViewTaskHandler(c *gin.Context) {
	message := ""
	taskId, err := strconv.Atoi(c.Param("ID"))

	if err != nil {
		message = "Task ID is of invalid format :" + err.Error()
		log.Printf("Task ID is of invalid format : %s resulted in error %s", c.Param("ID"), err.Error())
		c.HTML(http.StatusBadRequest, "CreateTask.html", gin.H{
			"message":      message,
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
	} else {
		taskItm, err := getTaskByID(taskId)
		if err != nil {
			log.Printf("Failed to fetch Task with ID  : %s %s", c.Param("ID"), err.Error())
			message = "Failed to fetch Task with ID  :" + c.Param("ID")
			c.HTML(http.StatusOK, "CreateTask.html", gin.H{
				"message":      message,
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
		} else {

			c.HTML(http.StatusOK, "CreateTask.html", gin.H{
				"message":      message,
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
		}
	}
}

func DeleteTaskHandler(c *gin.Context) {

	message := ""
	taskId, err := strconv.Atoi(c.Param("ID"))

	if err != nil {
		message = "Task ID is of invalid format :" + err.Error()
		log.Printf("Task ID is of invalid format : %s resulted in error %s", c.Param("ID"), err.Error())
		c.HTML(http.StatusBadRequest, "CreateTask.html", gin.H{
			"message":      message,
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
	} else {
		taskItm, err := getTaskByID(taskId)
		if err != nil {
			log.Printf("Failed to fetch Task with ID  : %s %s", c.Param("ID"), err.Error())
			message = "Failed to fetch Task with ID  :" + c.Param("ID")
		} else {
			result := CommonHandler.DbSQLConnection.Delete(taskItm)
			if result.Error != nil {
				message = "Delete Task Failed"
				log.Printf("Failed to delete the task : %s with error %s", c.Param("ID"), result.Error)
			} else {
				message = "Removed Task " + taskItm.Title
			}
		}
	}
	resultTasks, errFetchTasks := CommonHandler.GetAllTasks()
	if errFetchTasks != nil {
		log.Printf("Failed to fetch all tasks in the dashboard : %s ", errFetchTasks.Error())
		message = "Failed to fetch all tasks in the dashboard"
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title":      "Dashboard",
			"message":    message,
			"tasksTable": resultTasks,
		})
	} else {
		if len(resultTasks) == 0 {
			message = "There are no tasks for you in dashboard"
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Dashboard",
			"message":    message,
			"tasksTable": resultTasks,
		})
	}

}

func NavigateHomeHandler(c *gin.Context) {
	message := ""
	resultTasks, err := CommonHandler.GetAllTasks()

	if err != nil {
		log.Printf("Failed to fetch all tasks in the dashboard : %s ", err.Error())
		message = "Failed to fetch all tasks in the dashboard"
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title":      "Dashboard",
			"message":    message,
			"tasksTable": resultTasks,
		})
	} else {

		if len(resultTasks) == 0 {
			message = "There are no tasks for you in dashboard"

		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Dashboard",
			"message":    message,
			"tasksTable": resultTasks,
		})
	}

}
func NavigateEditTask(c *gin.Context) {
	message := ""
	taskId, err := strconv.Atoi(c.Param("ID"))

	if err != nil {
		log.Printf("Task ID is of invalid format : %s resulted in error %s", c.Param("ID"), err.Error())
	} else {
		taskItm, err := getTaskByID(taskId)
		if err != nil {
			log.Printf("Failed to fetch Task with ID  : %s %s", c.Param("ID"), err.Error())
			message = "Failed to fetch Task with ID  :" + c.Param("ID")
			c.HTML(http.StatusInternalServerError, "CreateTask.html", gin.H{
				"message":      message,
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
		} else {
			c.HTML(http.StatusOK, "CreateTask.html", gin.H{
				"message":      message,
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
			})
		}
	}
}

func CreateTaskHandler(c *gin.Context) {
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
		log.Printf("Error parsing due date: %s", err.Error())
		message = "Error parsing due date: " + err.Error()
		c.HTML(http.StatusBadRequest, "CreateTask.html", gin.H{
			"message":      message,
			"title":        "Create Task",
			"TaskTitle":    title,
			"Priority":     priority,
			"Description":  description,
			"DueDate":      dueDate,
			"Status":       status,
			"Category":     category,
			"Project":      project,
			"TaskAction":   "Create",
			"frmAction":    "/CreateTask",
			"EnableButton": "false",
		})
	}

	newTask := models.TaskItem{Title: title, Description: description, Priority: priority, Project: project, Category: category, DueDate: dateTimeParsed, Status: status}

	result := CommonHandler.DbSQLConnection.Create(&newTask)
	defer func(newTask *models.TaskItem) {
		if newTask.ID > 0 {
			message = "Created Task Successfully"
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
		} else {
			if result.Error != nil {
				message = "Create Task Failed"
				log.Printf("Failed to create the task : %s with error %s", newTask.Title, result.Error)
				c.HTML(http.StatusInternalServerError, "CreateTask.html", gin.H{
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
			}
		}

	}(&newTask)

}
