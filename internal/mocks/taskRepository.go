package mocks

import (
	"github.com/stretchr/testify/mock"
	"task-management/internal/task"
)

type TaskRepository struct {
	mock.Mock
}

func (tr *TaskRepository) CreateTask(task task.Task) (id int64, err error) {
	args := tr.Called(task)
	return args.Get(0).(int64), args.Error(1)
}

func (tr *TaskRepository) UpdateIsDone(id int64,isDone bool) error {
	args:=tr.Called(id,isDone)
	return args.Error(0)
}

func (tr *TaskRepository) UpdateIsDeleted(id int64,isDeleted bool)  error {
	args := tr.Called(id,isDeleted)
	return args.Error(0)
}

func (tr *TaskRepository) GetUserTaskList(description string,userID int64) ([]task.Task,error) {
	args := tr.Called(description,userID)
	return args.Get(0).([]task.Task), args.Error(1)
}

func (tr *TaskRepository) GetTaskByID(ID int64) (task.Task,error)  {
	args := tr.Called(ID)
	return args.Get(0).(task.Task),args.Error(1)

}
