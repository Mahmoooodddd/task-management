package task

import (
	"fmt"
	"net/http"
	"task-management/internal/response"
	"task-management/internal/user"
	"time"
)

type Service interface {
	CreateTask(u *user.User,params CreateTaskParams) (apiResponse response.ApiResponse, statusCode int)
	UpdateIsDone(u *user.User,params ChangeTaskToDoneParams) (apiResponse response.ApiResponse, statusCode int)
	UpdateIsDeleted(u *user.User,params ChangeTaskToDeletedParams) (apiResponse response.ApiResponse, statusCode int)
	GetUserTaskList(u *user.User,params GetUserTasksListParams) (apiResponse response.ApiResponse, statusCode int)
}

type service struct {
	taskRepository TaskRepository
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

func (s *service) CreateTask(u *user.User,params CreateTaskParams) (apiResponse response.ApiResponse, statusCode int) {
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
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	res := CreateTaskResponse{
		ID: result,
	}
	return response.Success(res, "")
}

func (s *service) UpdateIsDone(u *user.User,params ChangeTaskToDoneParams) (apiResponse response.ApiResponse, statusCode int) {
	err := s.taskRepository.UpdateIsDone(params.ID, params.IsDone)
	if err != nil {
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	return response.Success(nil, "")
}

func (s *service) UpdateIsDeleted(u *user.User,params ChangeTaskToDeletedParams) (apiResponse response.ApiResponse, statusCode int) {
	err := s.taskRepository.UpdateIsDeleted(params.ID, params.IsDeleted)
	if err != nil {
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	return response.Success(nil, "")
}

func (s *service) GetUserTaskList( u *user.User,params GetUserTasksListParams) (apiResponse response.ApiResponse, statusCode int) {
	tasks, err := s.taskRepository.GetUserTaskList(params.Description, u.ID)
	fmt.Println("==============================",err)
	if err != nil {
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

func NewService(taskRepository TaskRepository) Service {
	return &service{
		taskRepository: taskRepository,
	}
}
