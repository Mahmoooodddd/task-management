package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"task-management/internal/api"
	"task-management/internal/auth"
	"task-management/internal/platform"
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
		fmt.Println("err", err)
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

var testUser *userActor

type userActor struct {
	ID    int64
	Token string
	Email string
}

func getUserActor() *userActor {
	httpServer := api.NewHttpServer()

	if testUser != nil {
		return testUser
	}
	db := getDB()
	passwordEncoder := platform.NewPasswordEncoder()
	encodedPassword, _ := passwordEncoder.GenerateFromPassword("123456789")
	result, err := db.Exec(`INSERT INTO users(email,password) VALUES(?,?)`, "test@test.com", encodedPassword)
	if err != nil {
		fmt.Println("err", err)
	}
	id,err := result.LastInsertId()
	res := httptest.NewRecorder()
	data := `{"email":"test@test.com","password":"123456789"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	httpServer.GetEngine().ServeHTTP(res, req)
	response := struct {
		Status  bool
		Message string
		Data    auth.LoginResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		fmt.Println("err ", err)
	}

	ua := userActor{
		ID:    id,
		Token: response.Data.Token,
		Email: "test@test.com",
	}

	testUser = &ua
	return testUser
}
