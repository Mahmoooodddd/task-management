package platform

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordEncoder interface {
	GenerateFromPassword(password string) ([]byte, error)
	CompareHashAndPassword(hashedPassword string, password string) error
}

type passwordEncoder struct {
}

func (p *passwordEncoder) GenerateFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}
func (p *passwordEncoder) CompareHashAndPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
}

func NewPasswordEncoder() PasswordEncoder {
	return &passwordEncoder{}
}
