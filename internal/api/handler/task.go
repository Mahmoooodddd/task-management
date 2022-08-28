package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-management/internal/api/middleware"
	"task-management/internal/di"
	"task-management/internal/task"
	"task-management/internal/user"
)

func Create(c *gin.Context) {
	container := di.GetContainer()
	val, _ := c.Get(middleware.UserKey)
	u := val.(*user.User)
	taskService := container.GetTaskService()
	p := task.CreateTaskParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		errorResponse:=map[string]string{}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	resp,statusCode := taskService.CreateTask(u,p)
	c.JSON(statusCode, resp)
}

func UpdateIsDone(c *gin.Context) {
	container := di.GetContainer()
	val, _ := c.Get(middleware.UserKey)
	u := val.(*user.User)
	taskService := container.GetTaskService()
	p := task.ChangeTaskToDoneParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		errorResponse:=map[string]string{}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	resp,statusCode := taskService.UpdateIsDone(u,p)
	c.JSON(statusCode, resp)
}

func UpdateIsDeleted(c *gin.Context) {
	container := di.GetContainer()
	val, _ := c.Get(middleware.UserKey)
	u := val.(*user.User)
	taskService := container.GetTaskService()
	p := task.ChangeTaskToDeletedParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		errorResponse:=map[string]string{}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	resp,statusCode := taskService.UpdateIsDeleted(u,p)
	c.JSON(statusCode, resp)
}

func GetUserTaskList(c *gin.Context) {
	container := di.GetContainer()
	val, _ := c.Get(middleware.UserKey)
	u := val.(*user.User)
	taskService := container.GetTaskService()
	p := task.GetUserTasksListParams{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		errorResponse:=map[string]string{}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	resp,statusCode := taskService.GetUserTaskList(u,p)
	c.JSON(statusCode, resp)
}

