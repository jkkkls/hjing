package utils

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlDB *sql.DB

var MysqlConnCount = 10

func InitMysql(db string) error {
	var err error
	mysqlDB, err = sql.Open("mysql", db)
	if err != nil {
		return err
	}
	err = mysqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}

// 获取数据库操作对象
func GetMysqlDB() *sql.DB {
	return mysqlDB
}
