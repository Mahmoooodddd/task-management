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
	"task-management/internal/platform"
	"task-management/internal/task"
	"testing"
	"time"
)

type TaskTests struct {
	*suite.Suite
}

func TestTask(t *testing.T) {
	suite.Run(t, &TaskTests{
		Suite: new(suite.Suite),
	})
}

func (t *TaskTests) SetupSuite() {
}

func (t *TaskTests) SetupTest() {
	db := getDB()
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		fmt.Println(err)
	}
}

func (t *TaskTests) TearDownTest() {
}

func (t *TaskTests) TearDownSuite() {
}

func (t TaskTests) TestUpdateIsDone() {
	httpServer := api.NewHttpServer()
	db := getDB()
	userActor := getUserActor()
	taskModel, err := db.Exec(`INSERT INTO tasks(description, created_at,updated_at,user_id,is_done,is_deleted) VALUES(?,?,?,?,?,?)`,
		"doneTest", time.Now(), time.Now(), userActor.ID, false, false)
	taskId, err := taskModel.LastInsertId()
	data := fmt.Sprintf(`{ "id":%d,"isDone":true}`, taskId)
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/task/update-is-done", bytes.NewReader(body))
	token := "Bearer " + userActor.Token
	req.Header.Set("Authorization", token)
	httpServer.GetEngine().ServeHTTP(res, req)
	result := struct {
		Status  bool
		Message string
		Data    interface{}
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	//if err != nil {
	//	t.Fail(err.Error())
	//}
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), "", result.Message)
	assert.Nil(t.T(), err)
	updateTask := task.Task{}
	err = db.Get(&updateTask, `SELECT * from tasks where id=? limit 1`, taskId)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), true, updateTask.IsDone)

	data = fmt.Sprintf(`{ "id":%d,"isDone":false}`, taskId)
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/task/update-is-done", bytes.NewReader(body))
	token = "Bearer " + userActor.Token
	req.Header.Set("Authorization", token)
	httpServer.GetEngine().ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), "", result.Message)
	updateTask = task.Task{}
	err = db.Get(&updateTask, `SELECT * from tasks where id=? limit 1`, taskId)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), false, updateTask.IsDone)

}

func (t TaskTests) TestUpdateIsDone_TaskDoesNotBelongToUser() {
	httpServer := api.NewHttpServer()
	db := getDB()
	userActor := getUserActor()
	passwordEncoder := platform.NewPasswordEncoder()
	hashedPassword, _ := passwordEncoder.GenerateFromPassword("123456789")
	userModel, err := db.Exec(`INSERT INTO users(email,password) VALUES(?,?)`, "testtest@test.com", hashedPassword)
	assert.Nil(t.T(), err)
	userId, err := userModel.LastInsertId()
	assert.Nil(t.T(), err)
	taskModel, err := db.Exec(`INSERT INTO tasks(description, created_at,updated_at,user_id,is_done,is_deleted) VALUES(?,?,?,?,?,?)`,
		"doneTest", time.Now(), time.Now(), userId, false, false)
	taskId, err := taskModel.LastInsertId()
	data := fmt.Sprintf(`{ "id":%d,"isDone":true}`, taskId)
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/task/update-is-done", bytes.NewReader(body))
	token := "Bearer " + userActor.Token
	req.Header.Set("Authorization", token)
	httpServer.GetEngine().ServeHTTP(res, req)
	result := struct {
		Status  bool
		Message string
		Data    interface{}
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusNotFound, res.Code)
	assert.Equal(t.T(), "Not Found", result.Message)
	assert.Nil(t.T(), err)
}

func (t TaskTests) TestUpdateIsDeleted() {
	httpServer := api.NewHttpServer()
	db := getDB()
	userActor := getUserActor()
	taskModel, err := db.Exec(`INSERT INTO tasks(description, created_at,updated_at,user_id,is_done,is_deleted) VALUES(?,?,?,?,?,?)`,
		"deletedTest", time.Now(), time.Now(), userActor.ID, false, false)
	assert.Nil(t.T(), err)
	taskId, err := taskModel.LastInsertId()
	data := fmt.Sprintf(`{ "id":%d,"isDeleted":true}`, taskId)
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/task/update-is-deleted", bytes.NewReader(body))
	token := "Bearer " + userActor.Token
	req.Header.Set("Authorization", token)
	httpServer.GetEngine().ServeHTTP(res, req)
	result := struct {
		Status  bool
		Message string
		Data    interface{}
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	//if err != nil {
	//	t.Fail(err.Error())
	//}
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), "", result.Message)
	updateTask := task.Task{}
	err = db.Get(&updateTask, `SELECT * from tasks where id=? limit 1`, taskId)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), true, updateTask.IsDeleted)

	data = fmt.Sprintf(`{ "id":%d,"isDeleted":false}`, taskId)
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/task/update-is-deleted", bytes.NewReader(body))
	token = "Bearer " + userActor.Token
	req.Header.Set("Authorization", token)
	httpServer.GetEngine().ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), "", result.Message)
	assert.Nil(t.T(), err)
	updateTask = task.Task{}
	err = db.Get(&updateTask, `SELECT * from tasks where id=? limit 1`, taskId)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), false, updateTask.IsDeleted)

}

func (t TaskTests) TestCreateTask() {
	httpServer := api.NewHttpServer()
	db := getDB()
	userActor := getUserActor()
	data := `{"description":"new-test"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/task/create", bytes.NewReader(body))
	token := "Bearer " + userActor.Token
	req.Header.Set("Authorization", token)
	httpServer.GetEngine().ServeHTTP(res, req)
	result := struct {
		Status  bool
		Message string
		Data    task.CreateTaskResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	//if err != nil {
	//	t.Fail(err.Error())
	//}
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), "", result.Message)
	assert.NotEqual(t.T(), 0, result.Data.ID)
	createTask := task.Task{}
	err = db.Get(&createTask, `SELECT * from tasks where description = "new-test" limit 1`)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), false, createTask.IsDone)
	assert.Equal(t.T(), false, createTask.IsDeleted)
	assert.NotNil(t.T(), createTask.CreatedAt)
	assert.NotNil(t.T(), createTask.UpdatedAt)
}

func (t TaskTests) TestGetTaskUserList() {
	httpServer := api.NewHttpServer()
	db := getDB()
	userActor := getUserActor()
	_, err := db.Exec(`INSERT INTO tasks(id,description, created_at,updated_at,user_id,is_done,is_deleted) VALUES(?,?,?,?,?,?,?)`,
		3,"first-test", time.Now(), time.Now(), userActor.ID, false, false)
	_, err = db.Exec(`INSERT INTO tasks(id,description, created_at,updated_at,user_id,is_done,is_deleted) VALUES(?,?,?,?,?,?,?)`,
		4,"second-test", time.Now(), time.Now(), userActor.ID, true, true)
	assert.Nil(t.T(), err)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/task/list", nil)
	token := "Bearer " + userActor.Token
	req.Header.Set("Authorization", token)
	httpServer.GetEngine().ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	result := struct {
		Status  bool
		Message string
		Data    []task.SingleGetUserTaskListRes
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.Nil(t.T(), err)
	//
	//if err != nil {
	//	t.Fail(err.Error())
	//}
	assert.Equal(t.T(), "", result.Message)
	assert.Equal(t.T(), 2, len(result.Data))
	//assert.Equal(t.T(), int64(3), result.Data[0].ID)
	assert.Greater(t.T(),result.Data[0].ID,int64(0))
	assert.Equal(t.T(), "first-test", result.Data[0].Description)
	assert.Equal(t.T(), false, result.Data[0].IsDone)
	assert.Greater(t.T(),result.Data[1].ID,int64(0))
	assert.Equal(t.T(), "second-test", result.Data[1].Description)
	assert.Equal(t.T(), true, result.Data[1].IsDone)
}
