package db

import (
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func openSqliteDB(dsn, dbName string) (*gorm.DB, error) {
	//自带db参数
	if !strings.Contains(dsn, "/?") {
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "tb_",
				SingularTable: true,
			},
			SkipDefaultTransaction: true, //关闭事务
		})
	}

	newDsn := strings.ReplaceAll(dsn, "/?", fmt.Sprintf("/%v?", dbName))
	db, err := gorm.Open(sqlite.Open(newDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_",
			SingularTable: true,
		},
	})
	if err != nil && !strings.Contains(err.Error(), "Unknown database") {
		return nil, err
	} else if err == nil {
		return db, nil
	}

	//创建数据库
	temp, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Open: "+dsn)
	}

	err = temp.Exec(fmt.Sprintf("create database %v", dbName)).Error
	if err != nil {
		return nil, err
	}

	return gorm.Open(sqlite.Open(newDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_",
			SingularTable: true,
		},
	})
}

// InitSqlite 初始化数据库
// results 添加测试数据，调用需要计算好MysqlDB的函数调用顺序
func InitSqlite(dsn, dbName string, tables []any, results ...*MockResult) (*MysqlDB, error) {
	var (
		err error
		md  = &MysqlDB{}
	)

	if len(results) > 0 {
		md.Mock = &Mock{
			Results: results,
		}
		return md, nil
	}

	md.DB, err = openSqliteDB(dsn, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Open")
	}

	//初始化表
	err = md.DB.AutoMigrate(tables...)
	if err != nil {
		return nil, err
	}
	return md, nil
}
