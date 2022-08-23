package auth

import (
	"go.uber.org/zap"
	"net/http"
	"task-management/internal/platform"
	"task-management/internal/response"
	"task-management/internal/user"
)

type Service interface {
	Login(params LoginParams) (apiResponse response.ApiResponse, statusCode int)
	Register(params RegisterParams) (apiResponse response.ApiResponse, statusCode int)
	GetUser(token string) (*user.User, error)
}

type service struct {
	userService     user.Service
	passwordEncoder platform.PasswordEncoder
	jwtHandler      platform.JWTHandler
	logger          platform.Logger
}

type RegisterParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (s *service) Login(params LoginParams) (apiResponse response.ApiResponse, statusCode int) {
	userModel, err := s.userService.GetUserByEmail(params.Email)
	if err != nil {
		s.logger.Error("fail to get user by email", err,
			zap.String("service", "authService"),
			zap.String("method", "Login"),
			zap.String("email", params.Email),
			zap.String("password", params.Password),
		)
		return response.Error("Unauthorized", http.StatusUnauthorized, nil)
	}
	err = s.passwordEncoder.CompareHashAndPassword(userModel.Password, params.Password)
	if err != nil {
		s.logger.Error("fail to compare hash and password", err,
			zap.String("service", "authService"),
			zap.String("method", "Login"),
			zap.String("email", params.Email),
			zap.String("password", params.Password),
		)
		return response.Error("Unauthorized", http.StatusUnauthorized, nil)
	}
	getTokenParams := platform.GetTokenParams{
		Email: userModel.Email,
	}
	token, err := s.jwtHandler.GetToken(getTokenParams)
	if err != nil {
		s.logger.Error("fail to get token", err,
			zap.String("service", "authService"),
			zap.String("method", "Login"),
			zap.String("email", params.Email),
			zap.String("password", params.Password),
		)
		return response.Error("Unauthorized", http.StatusUnauthorized, nil)
	}
	res := LoginResponse{
		Token: token,
	}
	return response.Success(res, "")
}

func (s *service) Register(params RegisterParams) (apiResponse response.ApiResponse, statusCode int) {
	encodePassword, err := s.passwordEncoder.GenerateFromPassword(params.Password)
	if err != nil {
		s.logger.Error("fail to encode password", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", params.Email),
			zap.String("password", params.Password),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	err = s.userService.CreateUser(params.Email, string(encodePassword))
	if err != nil {
		s.logger.Error("fail to create user", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", params.Email),
			zap.String("password", params.Password),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	return response.Success(nil, "")
}

func (s *service) GetUser(token string) (*user.User, error) {
	email, err := s.jwtHandler.GetUsernameFromToken(token)
	if err != nil {
		return nil, err
	}
	u, err := s.userService.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func NewAuthService(userService user.Service, passwordEncoder platform.PasswordEncoder, jwtHandler platform.JWTHandler, logger platform.Logger) Service {
	return &service{
		userService:     userService,
		passwordEncoder: passwordEncoder,
		jwtHandler:      jwtHandler,
		logger:          logger,
	}
}
