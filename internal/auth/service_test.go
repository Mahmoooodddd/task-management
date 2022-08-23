package auth_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task-management/internal/auth"
	"task-management/internal/mocks"
	"task-management/internal/platform"
	"task-management/internal/user"

	"testing"
)

func TestService_Login_UserDoesNotExist(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	userService.On("GetUserByEmail", "test@test.test").Once().Return(user.User{}, fmt.Errorf("user not found"))
	logger.On("Error","fail to get user by email",mock.Anything,mock.Anything,mock.Anything,mock.Anything,mock.Anything).Once().Return(nil,logger.Error)

	params := auth.LoginParams{
		Email: "test@test.test",
	}
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	response, statusCode := authService.Login(params)
	assert.Equal(t, statusCode, 401)
		assert.Equal(t, response.Message, "Unauthorized")
	userService.AssertExpectations(t)
}

func TestService_Login_WrongPassword(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	userService.On("GetUserByEmail", "test@test.test").Once().Return(user.User{
		ID:       1,
		Email:    "test@test.test",
		Password: "123456789",
	}, nil)
	passwordEncoder.On("CompareHashAndPassword", "123456789", "987654321").Once().Return(fmt.Errorf("password is not correct"))
	params := auth.LoginParams{
		Email:    "test@test.test",
		Password: "987654321",
	}
	logger.On("Error","fail to compare hash and password",mock.Anything,mock.Anything,mock.Anything,mock.Anything,mock.Anything).Once().Return(nil,logger.Error)
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	response, statusCode := authService.Login(params)
	assert.Equal(t, statusCode, 401)
	assert.Equal(t, response.Message, "Unauthorized")
	userService.AssertExpectations(t)
	passwordEncoder.AssertExpectations(t)
}

func TestService_Login_Success(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	userService.On("GetUserByEmail", "test@test.test").Once().Return(user.User{
		ID:       1,
		Email:    "test@test.test",
		Password: "123456789",
	}, nil)
	passwordEncoder.On("CompareHashAndPassword", "123456789", "123456789").Once().Return(nil)
	params := auth.LoginParams{
		Email:    "test@test.test",
		Password: "123456789",
	}
	jwtHandler.On("GetToken", platform.GetTokenParams{Email: "test@test.test"}).Once().Return("token", nil)
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	response, statusCode := authService.Login(params)
	loginResponse := response.Data.(auth.LoginResponse)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, response.Message, "")
	assert.Equal(t, loginResponse.Token, "token")
	userService.AssertExpectations(t)
	passwordEncoder.AssertExpectations(t)
	jwtHandler.AssertExpectations(t)
}

func TestService_Register_CanNotGeneratePassword(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	passwordEncoder.On("GenerateFromPassword", "123456789").Once().Return([]byte{}, fmt.Errorf("can not generate password"))
	logger.On("Error","fail to encode password",mock.Anything,mock.Anything,mock.Anything,mock.Anything,mock.Anything).Once().Return(nil,logger.Error)
	params := auth.RegisterParams{
		Email:    "test@test.test",
		Password: "123456789",
	}
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	response, statusCode := authService.Register(params)
	assert.Equal(t, statusCode, 500)
	assert.Equal(t, response.Message, "something went wrong")
	userService.AssertExpectations(t)
	passwordEncoder.AssertExpectations(t)
}

func TestService_Register_UserServiceHasErr(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	passwordEncoder.On("GenerateFromPassword", "123456789").Once().Return([]byte("123456789"), nil)
	params := auth.RegisterParams{
		Email:    "test@test.test",
		Password: "123456789",
	}
	userService.On("CreateUser", "test@test.test", "123456789").Once().Return(fmt.Errorf("can not create user"))
	logger.On("Error","fail to create user",mock.Anything,mock.Anything,mock.Anything,mock.Anything,mock.Anything).Once().Return(nil,logger.Error)
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	response, statusCode := authService.Register(params)
	assert.Equal(t, statusCode, 500)
	assert.Equal(t, response.Message, "something went wrong")
	userService.AssertExpectations(t)
	passwordEncoder.AssertExpectations(t)
	jwtHandler.AssertExpectations(t)
}

func TestService_Register_Success(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	passwordEncoder.On("GenerateFromPassword", "123456789").Once().Return([]byte("123456789"), nil)
	params := auth.RegisterParams{
		Email:    "test@test.test",
		Password: "123456789",
	}
	userService.On("CreateUser", "test@test.test", "123456789").Once().Return(nil)
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	response, statusCode := authService.Register(params)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, response.Message, "")
	userService.AssertExpectations(t)
	passwordEncoder.AssertExpectations(t)
	jwtHandler.AssertExpectations(t)
}

func TestService_GetUser_JwtHasError(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	jwtHandler.On("GetUsernameFromToken", "token").Once().Return("", fmt.Errorf("can not get username from token"))
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	u, err := authService.GetUser("token")
	assert.Nil(t, u)
	assert.Error(t, err)
	userService.AssertExpectations(t)
	passwordEncoder.AssertExpectations(t)
	jwtHandler.AssertExpectations(t)
}

func TestService_GetUser_UserServiceHasError(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	jwtHandler.On("GetUsernameFromToken", "token").Once().Return("test@test.test", nil)
	userService.On("GetUserByEmail", "test@test.test").Once().Return(user.User{}, fmt.Errorf("can not get user by email"))
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	u, err := authService.GetUser("token")
	assert.Nil(t, u)
	assert.Error(t, err)
	userService.AssertExpectations(t)
	passwordEncoder.AssertExpectations(t)
	jwtHandler.AssertExpectations(t)
}

func TestService_GetUser_Success(t *testing.T) {
	userService := new(mocks.UserService)
	passwordEncoder := new(mocks.PasswordEncoder)
	jwtHandler := new(mocks.JwtHandler)
	logger := new(mocks.Logger)
	jwtHandler.On("GetUsernameFromToken", "token").Once().Return("test@test.test", nil)
	userService.On("GetUserByEmail", "test@test.test").Once().Return(user.User{
		ID:1,
		Email:"test@test.test",
		Password:"123456789",
	},nil)
	authService := auth.NewAuthService(userService, passwordEncoder, jwtHandler,logger)
	u, err := authService.GetUser("token")
	assert.Nil(t,err)
	assert.Equal(t,u.ID,int64(1))
	assert.Equal(t,u.Email,"test@test.test")
	assert.Equal(t,u.Password,"123456789")
	userService.AssertExpectations(t)
	passwordEncoder.AssertExpectations(t)
	jwtHandler.AssertExpectations(t)
}
