package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"task-management/internal/di"
)

const (
	UserKey = "user"
)

func AuthMiddleware(c *gin.Context) {
	h := authHeader{}
	err := c.ShouldBindHeader(&h)
	if err != nil {
		abort(c)
		return
	}
	parts := strings.Split(h.Token, " ")
	if len(parts) < 2 {
		abort(c)
		return
	}

	token := parts[1]
	container := di.NewContainer()
	authService := container.GetAuthService()
	user, err := authService.GetUser(token)
	if err != nil {
		abort(c)
		return
	}
	c.Set(UserKey, user)
	c.Next()
}

type authHeader struct {
	Token string `header:"Authorization"`
}

func abort(c *gin.Context) {
	res := map[string]string{}
	c.AbortWithStatusJSON(http.StatusUnauthorized, res)
}
