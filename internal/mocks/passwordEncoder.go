package mocks

import "github.com/stretchr/testify/mock"

type PasswordEncoder struct {
	mock.Mock
}

func (p *PasswordEncoder) GenerateFromPassword(password string) ([]byte,error)  {
	args:=p.Called(password)
	return args.Get(0).([]byte),args.Error(1)
}

func (p *PasswordEncoder) CompareHashAndPassword(hashedPassword string,password string) error  {
	args := p.Called(hashedPassword,password)
	return args.Error(0)
}
