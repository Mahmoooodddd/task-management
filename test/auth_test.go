package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"task-management/internal/api"
	"task-management/internal/auth"
	"task-management/internal/response"
	"task-management/internal/user"
	"testing"
)

type AuthTests struct {
	*suite.Suite
}

func TestAuth(t *testing.T) {
	suite.Run(t, &AuthTests{
		Suite: new(suite.Suite),
	})
}

func (t *AuthTests) SetupSuite() {
}

func (t *AuthTests) SetupTest() {
	db := getDB()
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		fmt.Println(err)
	}
}

func (t *AuthTests) TearDownTest() {
}

func (t *AuthTests) TearDownSuite() {
}

func (t *AuthTests) TestRegister() {
	httpServer := api.NewHttpServer()
	data := `{"email":"test@test.test","password":"123456789"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	httpServer.GetEngine().ServeHTTP(res, req)
	result := response.ApiResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), "", result.Message)
	db := getDB()
	u := user.User{}
	err = db.Get(&u, `SELECT * FROM users where email="test@test.test" LIMIT 1`)
	assert.Nil(t.T(), err)
}

func (t *AuthTests) TestLogin() {
	httpServer := api.NewHttpServer()
	db := getDB()
	passwordEncoder := getContainer().GetPasswordEncoder()
	hashedPassword, _ := passwordEncoder.GenerateFromPassword("123456789")
	_, err := db.Exec(`INSERT INTO users(email,password) VALUES(?,?)`, "test2@test.com", hashedPassword)
	data := `{"email":"test2@test.com","password":"123456789"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	httpServer.GetEngine().ServeHTTP(res, req)
	result := struct {
		Status  bool
		Message string
		Data    auth.LoginResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), "", result.Message)
	assert.Nil(t.T(), err)
	assert.NotEqual(t.T(), result.Data.Token, "")
}
