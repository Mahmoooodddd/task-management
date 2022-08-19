package mocks

import (
	"github.com/stretchr/testify/mock"
	"task-management/internal/platform"
)

type JwtHandler struct {
	mock.Mock
}

func (j *JwtHandler) GetToken(params platform.GetTokenParams) (string, error) {
	args := j.Called(params)
	return args.Get(0).(string), args.Error(1)
}

func (j *JwtHandler) GetUsernameFromToken(signedToken string) (string, error) {
	args := j.Called(signedToken)
	return args.Get(0).(string),args.Error(1)
}
