package db

import (
	"fmt"

	"github.com/jkkkls/hjing/db"
)

var (
	wdDB   *db.MysqlDB
	models = []any{
		// end models
	}
)

// InitDB 初始化数据库
func InitDB(dbType, dbName, dsn string) error {
	var err error
	switch dbType {
	case "mysql":
		wdDB, err = db.InitMysql(dsn, dbName, models)
	case "sqlite":
		wdDB, err = db.InitMysql(dsn, dbName, models)
	case "postgres":
		wdDB, err = db.InitPg(dsn, dbName, models)
	default:
		return fmt.Errorf("错误的数据库配置，请检查配置db项，需要type和dsn字段")
	}
	_ = wdDB
	return err
}
