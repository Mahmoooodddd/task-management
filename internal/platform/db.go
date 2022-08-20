package platform

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDBClient() *sqlx.DB {
	if flag.Lookup("test.v") != nil{
		//here means test env
		db, err := sqlx.Connect("mysql", "task_management_user:task_management_pass@tcp(localhost:3306)/task_management_test?parseTime=true")
		if err != nil {
			panic("can not connect to database "+err.Error())
		}
		return db

	}

	db, err := sqlx.Connect("mysql", "task_management_user:task_management_pass@tcp(localhost:3306)/task_management?parseTime=true")
	if err != nil {
		panic("can not connect to database "+err.Error())
	}
	return db
}