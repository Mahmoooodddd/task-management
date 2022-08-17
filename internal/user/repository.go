package user

import (
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user User) (id int64, err error)
	GetUserByEmail(email string) (user User, err error)
}

type userRepository struct {
	dbClient *sqlx.DB
}

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ur *userRepository) CreateUser(user User) (id int64, err error) {
	sqlStr := "INSERT INTO users(email, password) VALUES(?, ?)"
	result, err := ur.dbClient.Exec(sqlStr, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (ur *userRepository) GetUserByEmail(email string) (user User, err error) {
	sqlStr := "SELECT * from users where email = ? LIMIT 1"
	user = User{}
	row := ur.dbClient.QueryRow(sqlStr, email)
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}

func NewUserRepository(dbClient *sqlx.DB) UserRepository {
	return &userRepository{
		dbClient: dbClient,
	}
}
