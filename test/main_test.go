package test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
	"testing"
)

const (
	DbDNS = "task_management_user:task_management_pass@tcp(localhost:3306)/task_management_test?multiStatements=true&parseTime=true"
)

var mainDB *sqlx.DB

func getDB() *sqlx.DB {
	if mainDB != nil {
		return mainDB
	}
	db, err := sqlx.Connect("mysql", DbDNS)
	if err != nil {
		fmt.Println("err",err)
		panic("can not establish connection to database")
	}
	mainDB = db
	return mainDB
}

func TestMain(m *testing.M) {
	setupDB()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupDB() {
	db := getDB()
	content, err := ioutil.ReadFile("./data/db.sql")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = db.Exec(string(content))
	if err != nil {
		fmt.Println(err.Error())
	}
}
