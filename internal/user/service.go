package user

type Service interface {
	CreateUser(email string, encodedPassword string) error
	GetUserByEmail(email string) (user User,err error)
}

type service struct {
	userRepository UserRepository
}

func (s *service) CreateUser(email string, encodedPassword string) error {
	u := User{
		Email:    email,
		Password: encodedPassword,
	}
	_, err := s.userRepository.CreateUser(u)
	return err
}

func (s *service) GetUserByEmail(email string) (user User,err error) {
	user,err = s.userRepository.GetUserByEmail(email)
	return user,err
}

func NewService(userRepository UserRepository) Service {
	return &service{
		userRepository: userRepository,
	}
}
