package dbadmin

import (
	"fmt"
	"webserver/db"
)

var (
	dbName  = "axjmin"
	adminDB *db.MysqlDB
	tables  []interface{}
)

// InitAdminDB 初始化数据库
func InitAdminDB(t, dsn string) error {
	var err error
	switch t {
	case "mysql":
		adminDB, err = db.InitMysql(dsn, dbName, tables)
	case "sqlite":
		adminDB, err = db.InitSqlite(dsn, dbName, tables)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", t)
	}
	if err != nil {
		return err
	}

	return nil
}

func Save(a any, fields ...string) error {
	return adminDB.Save(a, fields...)
}

func QueryList(data any, preload, order string, limit, index int, query interface{}, args ...interface{}) (n int64, err error) {
	n, err = adminDB.QueryList(data, preload, order, limit, index, query, args...)
	return
}

// Query 查找
func Query(data any, query any, args ...any) error {
	err := adminDB.Query(data, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除
func Delete(data any, query any, args ...any) error {
	return adminDB.Delete(data, query, args...)
}

func AssociationDelete(data any) error {
	return adminDB.AssociationDelete(data)
}

func QueryAll(data any, order string, query any, args ...any) (err error) {
	err = adminDB.QueryAll(data, order, 0, query, args...)
	return
}
