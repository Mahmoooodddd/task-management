package platform

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDBClient() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "task_management_user:task_management_pass@tcp(localhost:3306)/task_management?parseTime=true")
	if err != nil {
		panic("can not connect to database "+err.Error())
	}
	return db
}
