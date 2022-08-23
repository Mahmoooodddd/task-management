package task_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task-management/internal/mocks"
	"task-management/internal/task"
	"task-management/internal/user"
	"testing"
)

func TestService_GetUserTaskList_RepoHasError(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("GetUserTaskList", "test", int64(1)).Once().Return([]task.Task{}, fmt.Errorf("can not show list"))
	logger.On("Error","fail to get user task list",mock.Anything,mock.Anything,mock.Anything,mock.Anything).Once().Return([]task.Task{},logger.Error)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 1,
	}
	params := task.GetUserTasksListParams{
		Description: "test",
	}
	response, statusCode := taskService.GetUserTaskList(u, params)
	assert.Equal(t, statusCode, 500)
	assert.Equal(t, response.Message, "something went wrong")
	taskRepository.AssertExpectations(t)
}

func TestService_GetUserTaskList_Success(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	tasksRes := []task.Task{
		{
			ID:          1,
			Description: "test",
			IsDone:      false,
		}, {
			ID:          2,
			Description: "test2",
			IsDone:      true,
		},
	}
	taskRepository.On("GetUserTaskList", "test", int64(1)).Once().Return(tasksRes, nil)
	taskService := task.NewService(taskRepository,logger)

	u := &user.User{
		ID: 1,
	}
	params := task.GetUserTasksListParams{
		Description: "test",
	}
	response, statusCode := taskService.GetUserTaskList(u, params)
	res := response.Data.([]task.SingleGetUserTaskListRes)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, response.Message, "")
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0].ID, int64(1))
	assert.Equal(t, res[0].Description, "test")
	assert.Equal(t, res[0].IsDone, false)
	assert.Equal(t, res[1].ID, int64(2))
	assert.Equal(t, res[1].Description, "test2")
	assert.Equal(t, res[1].IsDone, true)
	taskRepository.AssertExpectations(t)
}

func TestService_CreateTask_RepoHasError(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("CreateTask", mock.Anything).Once().Return(int64(0), fmt.Errorf("can not create task"))
	logger.On("Error","fail to create task",mock.Anything,mock.Anything,mock.Anything,mock.Anything).Once().Return([]task.Task{},logger.Error)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 1,
	}
	params := task.CreateTaskParams{
		Description: "test",
	}
	response, statusCode := taskService.CreateTask(u, params)
	assert.Equal(t, statusCode, 500)
	assert.Equal(t, response.Message, "something went wrong")
	taskRepository.AssertExpectations(t)
}

func TestService_CreateTask_Success(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("CreateTask", mock.Anything).Once().Return(int64(1), nil)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 1,
	}
	params := task.CreateTaskParams{
		Description: "test",
	}
	response, statusCode := taskService.CreateTask(u, params)
	taskResponse := response.Data.(task.CreateTaskResponse)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, response.Message, "")
	assert.Equal(t, taskResponse.ID, int64(1))
	taskRepository.AssertExpectations(t)
}

func TestService_UpdateIsDone_GetTaskHasError(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("GetTaskByID", int64(1)).Once().Return(task.Task{},fmt.Errorf("can not get task"))
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 1,
	}
	params := task.ChangeTaskToDoneParams{
		ID:     1,
		IsDone: false,
	}
	response, statusCode := taskService.UpdateIsDone(u, params)
	assert.Equal(t, statusCode, 404)
	assert.Equal(t, response.Message, "Not Found")
	taskRepository.AssertExpectations(t)
}

func TestService_UpdateIsDone_TaskDoesNotBelongUser(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("GetTaskByID", int64(1)).Once().Return(task.Task{
		ID:     1,
		UserId: 1,
	}, nil)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 2,
	}
	params := task.ChangeTaskToDoneParams{
		ID:     1,
		IsDone: false,
	}
	response, statusCode := taskService.UpdateIsDone(u, params)
	assert.Equal(t, statusCode, 404)
	assert.Equal(t, response.Message, "Not Found")
	taskRepository.AssertExpectations(t)
}



func TestService_UpdateIsDone_RepoHasError(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("GetTaskByID", int64(1)).Once().Return(task.Task{
		ID:     1,
		UserId: 1,
	}, nil)
	taskRepository.On("UpdateIsDone", int64(1), false).Once().Return(fmt.Errorf("can not update is done"))
	logger.On("Error","fail to update is done",mock.Anything,mock.Anything,mock.Anything,mock.Anything,mock.Anything).Once().Return([]task.Task{},logger.Error)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 1,
	}
	params := task.ChangeTaskToDoneParams{
		ID:     1,
		IsDone: false,
	}
	response, statusCode := taskService.UpdateIsDone(u, params)
	assert.Equal(t, statusCode, 500)
	assert.Equal(t, response.Message, "something went wrong")
	taskRepository.AssertExpectations(t)
}

func TestService_UpdateIsDeleted_TaskDoesNotBelongUser(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("GetTaskByID", int64(1)).Once().Return(task.Task{
		ID:     1,
		UserId: 1,
	}, nil)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 2,
	}
	params := task.ChangeTaskToDoneParams{
		ID:     1,
		IsDone: false,
	}
	response, statusCode := taskService.UpdateIsDone(u, params)
	assert.Equal(t, statusCode, 404)
	assert.Equal(t, response.Message, "Not Found")
	taskRepository.AssertExpectations(t)
}


func TestService_UpdateIsDone_Success(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("GetTaskByID", int64(1)).Once().Return(task.Task{
		ID:     1,
		UserId: 1,
	}, nil)
	taskRepository.On("UpdateIsDone", int64(1), false).Once().Return(nil)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 1,
	}
	params := task.ChangeTaskToDoneParams{
		ID:     1,
		IsDone: false,
	}
	response, statusCode := taskService.UpdateIsDone(u, params)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, response.Message, "")
	taskRepository.AssertExpectations(t)
}

func TestService_UpdateIsDeleted_RepoHasError(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("GetTaskByID", int64(1)).Once().Return(task.Task{
		ID:     1,
		UserId: 1,
	}, nil)
	taskRepository.On("UpdateIsDeleted", int64(1), false).Once().Return(fmt.Errorf("can not update is deleted"))
	logger.On("Error","fail to update is deleted",mock.Anything,mock.Anything,mock.Anything,mock.Anything,mock.Anything,mock.Anything).Once().Return([]task.Task{},logger.Error)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 1,
	}
	params := task.ChangeTaskToDeletedParams{
		ID:        1,
		IsDeleted: false,
	}
	response, statusCode := taskService.UpdateIsDeleted(u, params)
	assert.Equal(t, statusCode, 500)
	assert.Equal(t, response.Message, "something went wrong")
	taskRepository.AssertExpectations(t)
}

func TestService_UpdateIsDeleted_Success(t *testing.T) {
	taskRepository := new(mocks.TaskRepository)
	logger := new(mocks.Logger)
	taskRepository.On("GetTaskByID", int64(1)).Once().Return(task.Task{
		ID:     1,
		UserId: 1,
	}, nil)
	taskRepository.On("UpdateIsDeleted", int64(1), false).Once().Return(nil)
	taskService := task.NewService(taskRepository,logger)
	u := &user.User{
		ID: 1,
	}
	params := task.ChangeTaskToDeletedParams{
		ID:        1,
		IsDeleted: false,
	}
	response, statusCode := taskService.UpdateIsDeleted(u, params)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, response.Message, "")
	taskRepository.AssertExpectations(t)
}
