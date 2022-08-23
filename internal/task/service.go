package task

import (
	"go.uber.org/zap"
	"net/http"
	"task-management/internal/platform"
	"task-management/internal/response"
	"task-management/internal/user"
	"time"
)

type Service interface {
	CreateTask(u *user.User, params CreateTaskParams) (apiResponse response.ApiResponse, statusCode int)
	UpdateIsDone(u *user.User, params ChangeTaskToDoneParams) (apiResponse response.ApiResponse, statusCode int)
	UpdateIsDeleted(u *user.User, params ChangeTaskToDeletedParams) (apiResponse response.ApiResponse, statusCode int)
	GetUserTaskList(u *user.User, params GetUserTasksListParams) (apiResponse response.ApiResponse, statusCode int)
}

type service struct {
	taskRepository TaskRepository
	logger         platform.Logger
}

type CreateTaskParams struct {
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	ID int64 `json:"id"`
}

type ChangeTaskToDoneParams struct {
	ID     int64 `json:"id"`
	IsDone bool  `json:"isDone"`
}

type ChangeTaskToDeletedParams struct {
	ID        int64 `json:"id"`
	IsDeleted bool  `json:"isDeleted"`
}

type GetUserTasksListParams struct {
	Description string `form:"description"`
}

type SingleGetUserTaskListRes struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	IsDone      bool   `json:"isDone"`
}

func (s *service) CreateTask(u *user.User, params CreateTaskParams) (apiResponse response.ApiResponse, statusCode int) {
	task := Task{
		Description: params.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserId:      u.ID,
		IsDone:      false,
		IsDeleted:   false,
	}
	result, err := s.taskRepository.CreateTask(task)
	if err != nil {
		s.logger.Error("fail to create task", err,
			zap.String("service", "taskService"),
			zap.String("method", "CreateTask"),
			zap.String("description", params.Description),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	res := CreateTaskResponse{
		ID: result,
	}
	return response.Success(res, "")
}

func (s *service) UpdateIsDone(u *user.User, params ChangeTaskToDoneParams) (apiResponse response.ApiResponse, statusCode int) {
	task, err := s.taskRepository.GetTaskByID(params.ID)
	if err != nil {
		return response.Error("Not Found", http.StatusNotFound, nil)
	}
	if task.UserId != u.ID {
		return response.Error("Not Found", http.StatusNotFound, nil)
	}
	err = s.taskRepository.UpdateIsDone(params.ID, params.IsDone)
	if err != nil {
		s.logger.Error("fail to update is done", err,
			zap.String("service", "taskService"),
			zap.String("method", "UpdateIsDone"),
			zap.Int64("id", params.ID),
			zap.Bool("isDone", params.IsDone),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	return response.Success(nil, "")
}

func (s *service) UpdateIsDeleted(u *user.User, params ChangeTaskToDeletedParams) (apiResponse response.ApiResponse, statusCode int) {
	task, err := s.taskRepository.GetTaskByID(params.ID)
	if err != nil {
		return response.Error("Not Found", http.StatusNotFound, nil)
	}
	if task.UserId != u.ID {
		return response.Error("Not Found", http.StatusNotFound, nil)
	}

	err = s.taskRepository.UpdateIsDeleted(params.ID, params.IsDeleted)
	if err != nil {
		s.logger.Error("fail to update is deleted", err,
			zap.String("service", "taskService"),
			zap.String("method", "UpdateIsDeleted"),
			zap.Int64("id", params.ID),
			zap.Bool("isDone", params.IsDeleted),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	return response.Success(nil, "")
}

func (s *service) GetUserTaskList(u *user.User, params GetUserTasksListParams) (apiResponse response.ApiResponse, statusCode int) {
	tasks, err := s.taskRepository.GetUserTaskList(params.Description, u.ID)
	if err != nil {
		s.logger.Error("fail to get user task list", err,
			zap.String("service", "taskService"),
			zap.String("method", "GetUserTaskList"),
			zap.String("description", params.Description),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	var res []SingleGetUserTaskListRes
	for _, task := range tasks {
		singleRes := SingleGetUserTaskListRes{
			ID:          task.ID,
			Description: task.Description,
			IsDone:      task.IsDone,
		}
		res = append(res, singleRes)
	}
	return response.Success(res, "")
}

func NewService(taskRepository TaskRepository,logger platform.Logger) Service {
	return &service{
		taskRepository: taskRepository,
		logger:         logger,
	}
}
