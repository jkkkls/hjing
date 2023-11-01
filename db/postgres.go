package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// select setval(pg_get_serial_sequence('tb_player', 'id'), 100000, false);
// dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
func openPgDB(dsn, dbName string) (*gorm.DB, error) {
	conf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_",
			SingularTable: true,
		},
		//Logger: logger.Default.LogMode(logger.Info),
	}
	//注册json标签处理逻辑
	// schema.RegisterSerializer("json", JSONSerializer{})

	//自带db参数
	if strings.Contains(dsn, "dbname=") {
		return gorm.Open(postgres.Open(dsn), conf)
	}

	newDsn := dsn + " dbname=" + dbName
	db, err := gorm.Open(postgres.Open(newDsn), conf)
	log.Println(newDsn, err)
	if err != nil && !strings.Contains(err.Error(), "SQLSTATE 3D000") {
		return nil, err
	} else if err == nil {
		return db, nil
	}

	//创建数据库
	temp, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Open: "+dsn)
	}

	err = temp.Exec(fmt.Sprintf("create database %v", dbName)).Error
	if err != nil {
		return nil, err
	}

	return gorm.Open(postgres.Open(newDsn), conf)
}

// InitMysql 初始化数据库
// results 添加测试数据，调用需要计算好MysqlDB的函数调用顺序
func InitPg(dsn, dbName string, tables []interface{}, results ...*MockResult) (*MysqlDB, error) {
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

	md.DB, err = openPgDB(dsn, dbName)
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
