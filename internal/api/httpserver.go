package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-management/internal/api/handler"
	"task-management/internal/api/middleware"
)

type HTTPServer interface {
	ListenAndServe(address string) error
	GetEngine() http.Handler
}

type httpServer struct {
	server      *http.Server
	engine      *gin.Engine
}

func (s *httpServer) ListenAndServe(address string) error  {
	s.server.Addr = address
	return s.server.ListenAndServe()
}

func (s *httpServer) GetEngine() http.Handler  {
	return s.engine
}

func NewHttpServer() HTTPServer  {
	apiRouter := gin.New()
	server := &http.Server{
		Addr:           "0.0.0.0:8000",
		Handler:        apiRouter,
	}
	s:=&httpServer{
		server: server,
		engine:   apiRouter,
	}
	s.registerRoutes()
	return s
}

func (s *httpServer) registerRoutes()  {
	r := s.engine
	v1 := r.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", handler.Login)
			authRoutes.POST("/register", handler.Register)
		}

		taskRoutes := v1.Group("/task")
		taskRoutes.Use(middleware.AuthMiddleware)
		{
			taskRoutes.POST("/create", handler.Create)
			taskRoutes.POST("/update-is-done", handler.UpdateIsDone)
			taskRoutes.POST("/update-is-deleted", handler.UpdateIsDeleted)
			taskRoutes.GET("/list", handler.GetUserTaskList)
		}
	}
}