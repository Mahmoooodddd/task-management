package task

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type TaskRepository interface {
	CreateTask(task Task) (id int64, err error)
	UpdateIsDone(ID int64, isDone bool) error
	UpdateIsDeleted(ID int64, isDeleted bool) error
	GetUserTaskList(description string, userID int64) ([]Task, error)
}

type taskRepository struct {
	dbClient *sqlx.DB
}

type Task struct {
	ID          int64     `db:"id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	UserId      int64     `db:"user_id"`
	IsDone      bool      `db:"is_done"`
	IsDeleted   bool      `db:"is_deleted"`
}

func (tr *taskRepository) CreateTask(task Task) (id int64, err error) {
	sqlStr := "INSERT INTO tasks(description, created_at,updated_at,user_id,is_done,is_deleted) VALUES(?,?,?,?,?,?)"
	result, err := tr.dbClient.Exec(sqlStr, task.Description, task.CreatedAt, task.UpdatedAt, task.UserId, task.IsDone, task.IsDeleted)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (tr *taskRepository) UpdateIsDone(ID int64, isDone bool) error {
	sqlStr := "UPDATE tasks set is_done=? where id=?"
	_, err := tr.dbClient.Exec(sqlStr, isDone, ID)
	if err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) UpdateIsDeleted(ID int64, isDeleted bool) error {
	sqlStr := "UPDATE tasks set is_deleted=? where id=?"
	_, err := tr.dbClient.Exec(sqlStr, isDeleted, ID)
	if err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) GetUserTaskList(description string, userID int64) ([]Task, error) {
	tasks := []Task{}
	sqlStr := "SELECT * from tasks where user_id=?"
	rows, err := tr.dbClient.Query(sqlStr, userID)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		t := Task{}
		if err := rows.Scan(
			&t.ID,
			&t.Description,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.IsDone,
			&t.IsDeleted,
			&t.UserId,
		); err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}

	return tasks, err
}

func NewTaskRepository(dbClient *sqlx.DB) TaskRepository {
	return &taskRepository{
		dbClient: dbClient,
	}
}
