package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-management/internal/auth"
	"task-management/internal/di"
)


func Login(c *gin.Context) {
	container := di.NewContainer()
	authService := container.GetAuthService()
	p := auth.LoginParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		errorResponse:=map[string]string{}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	resp,statusCode := authService.Login(p)
	c.JSON(statusCode, resp)
}

func Register(c *gin.Context) {
	container := di.NewContainer()
	authService := container.GetAuthService()
	p := auth.RegisterParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		errorResponse:=map[string]string{}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	resp,statusCode := authService.Register(p)
	c.JSON(statusCode, resp)
}