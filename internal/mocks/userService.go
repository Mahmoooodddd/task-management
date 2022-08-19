package mocks

import (
	"github.com/stretchr/testify/mock"
	"task-management/internal/user"
)

type UserService struct {
	mock.Mock
}

func (u *UserService) CreateUser(email string,encodedPassword string) error  {
	args :=u.Called(email,encodedPassword)
	return args.Error(0)
}

func (u *UserService) GetUserByEmail(email string)(user.User, error)  {
	args := u.Called(email)
	return args.Get(0).(user.User), args.Error(1)
}