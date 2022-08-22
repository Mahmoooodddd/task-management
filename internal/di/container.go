package di

import (
	"github.com/jmoiron/sqlx"
	"task-management/config"
	"task-management/internal/auth"
	"task-management/internal/platform"
	"task-management/internal/task"
	"task-management/internal/user"
)

type container struct {
	authService     auth.Service
	userService     user.Service
	passwordEncoder platform.PasswordEncoder
	dbClient        *sqlx.DB
	userRepository  user.UserRepository
	taskRepository  task.TaskRepository
	taskService     task.Service
	jwtHandler      platform.JWTHandler
	configs          platform.Configs
}

func (c *container) GetAuthService() auth.Service {
	if c.authService == nil {
		userService := c.GetUserService()
		passwordEncoder := c.GetPasswordEncoder()
		jwtHandler := c.GetJwtHandler()
		authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler)
		c.authService = authService
	}
	return c.authService
}

func (c *container) GetUserService() user.Service {
	if c.userService == nil {
		userRepository := c.GetUserRepository()
		userService := user.NewService(userRepository)
		c.userService = userService
	}
	return c.userService
}

func (c *container) GetPasswordEncoder() platform.PasswordEncoder {
	if c.passwordEncoder == nil {
		passwordEncoder := platform.NewPasswordEncoder()
		c.passwordEncoder = passwordEncoder
	}
	return c.passwordEncoder
}

func (c *container) GetConfig() platform.Configs {
	if c.configs == nil {
		viper := config.SetConfigs()
		configs := platform.NewConfigs(viper)
		c.configs = configs
	}
	return c.configs
}

func (c *container) GetDbClient() *sqlx.DB {
	if c.dbClient == nil {
		configs := c.GetConfig()
		dbClient := platform.NewDBClient(configs)
		c.dbClient = dbClient
	}
	return c.dbClient
}

func (c *container) GetUserRepository() user.UserRepository {
	if c.userRepository == nil {
		dbClient := c.GetDbClient()
		userRepository := user.NewUserRepository(dbClient)
		c.userRepository = userRepository
	}
	return c.userRepository
}

func (c *container) GetTaskRepository() task.TaskRepository {
	if c.taskRepository == nil {
		dbClient := c.GetDbClient()
		taskRepository := task.NewTaskRepository(dbClient)
		c.taskRepository = taskRepository
	}
	return c.taskRepository
}

func (c *container) GetTaskService() task.Service {
	if c.taskService == nil {
		taskRepository := c.GetTaskRepository()
		taskService := task.NewService(taskRepository)
		c.taskService = taskService
	}
	return c.taskService
}

func (c *container) GetJwtHandler() platform.JWTHandler {
	if c.jwtHandler == nil {
		jwtHandler := platform.NewJWTHandler()
		c.jwtHandler = jwtHandler
	}
	return c.jwtHandler
}

func NewContainer() *container {
	return &container{}
}
