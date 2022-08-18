package user_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"task-management/internal/mocks"
	"task-management/internal/user"
	"testing"
)

func TestService_GetUserByEmail_RepoHasError(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := user.NewService(userRepository)
	userRepository.On("GetUserByEmail", "test@test.test").Once().Return(user.User{}, fmt.Errorf("user not found"))
	_, err := userService.GetUserByEmail("test@test.test")
	assert.Equal(t, "user not found", err.Error())
	userRepository.AssertExpectations(t)

}

func TestService_GetUserByEmail_Successful(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := user.NewService(userRepository)
	userRepository.On("GetUserByEmail","test@test.test").Once().Return(user.User{
		ID :1,
		Email:"test@test.test",
		Password:"123456789",
	},nil)
	u,err := userService.GetUserByEmail("test@test.test")
	assert.Nil(t,err)
	assert.Equal(t,u.ID,int64(1))
	assert.Equal(t,u.Email,"test@test.test")
	assert.Equal(t,u.Password,"123456789")
	userRepository.AssertExpectations(t)
}

func TestService_CreateUser_RepoHasError(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := user.NewService(userRepository)
	userRepository.On("CreateUser", user.User{ID:0,Email:"test@test.test",Password:"123456789"}).Once().Return(int64(0), fmt.Errorf("can not create user"))
	err := userService.CreateUser("test@test.test","123456789")
	assert.Equal(t, "can not create user", err.Error())
	userRepository.AssertExpectations(t)
}

func TestService_CreateUser_Successful(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := user.NewService(userRepository)
	userRepository.On("CreateUser", user.User{Email:"test@test.test",Password:"123456789"}).Once().Return(int64(1),nil)
	err := userService.CreateUser("test@test.test","123456789")
	assert.Nil(t,err)
	userRepository.AssertExpectations(t)
}

