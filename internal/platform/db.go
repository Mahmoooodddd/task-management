package platform

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDBClient(c Configs) *sqlx.DB {
	username := c.GetString("database.username")
	password := c.GetString("database.password")
	host := c.GetString("database.host")
	port := c.GetInt("database.port")
	dbName := c.GetString("database.dbName")
	driver := c.GetString("database.driver")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", username, password, host, port, dbName)
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		panic("can not connect to database " + err.Error())
	}
	return db
}
