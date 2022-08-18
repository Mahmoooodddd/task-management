package mocks

import (
	"github.com/stretchr/testify/mock"
	"task-management/internal/user"
)

type UserRepository struct {
	mock.Mock
}

func (ur *UserRepository) CreateUser(user user.User) (id int64, err error) {
	args := ur.Called(user)
	return args.Get(0).(int64), args.Error(1)
}
func (ur *UserRepository) GetUserByEmail(email string) (user.User, error) {
	args := ur.Called(email)
	return args.Get(0).(user.User), args.Error(1)
}
