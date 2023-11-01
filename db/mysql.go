package db

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type MockResult struct {
	Data  interface{}
	Err   error
	Count int64
}

type Mock struct {
	Results []*MockResult
}

type MysqlDB struct {
	DB   *gorm.DB
	Mock *Mock //测试数据
}

func openDB(dsn, dbName string) (*gorm.DB, error) {
	conf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_",
			SingularTable: true,
		},
		//Logger: logger.Default.LogMode(logger.Info),
	}
	//注册json标签处理逻辑
	schema.RegisterSerializer("json", JSONSerializer{})

	//自带db参数
	if !strings.Contains(dsn, "/?") {
		return gorm.Open(mysql.Open(dsn), conf)
	}

	newDsn := strings.ReplaceAll(dsn, "/?", fmt.Sprintf("/%v?", dbName))
	db, err := gorm.Open(mysql.Open(newDsn), conf)
	if err != nil && !strings.Contains(err.Error(), "Unknown database") {
		return nil, err
	} else if err == nil {
		return db, nil
	}

	//创建数据库
	temp, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Open: "+dsn)
	}

	err = temp.Exec(fmt.Sprintf("create database %v", dbName)).Error
	if err != nil {
		return nil, err
	}

	return gorm.Open(mysql.Open(newDsn), conf)
}

// InitMysql 初始化数据库
// results 添加测试数据，调用需要计算好MysqlDB的函数调用顺序
func InitMysql(dsn, dbName string, tables []interface{}, results ...*MockResult) (*MysqlDB, error) {
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

	md.DB, err = openDB(dsn, dbName)
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

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
	return err
}

func (md *MysqlDB) Exec(sql string, values ...interface{}) error {
	//mock
	if md.Mock != nil {
		return nil
	}

	return md.DB.Exec(sql, values...).Error
}

// QueryOne 查找用户
func (md *MysqlDB) QueryOne(k string, v interface{}, data interface{}, fields ...string) error {
	//mock
	if md.Mock != nil && len(md.Mock.Results) > 0 {
		m := md.Mock.Results[0]
		md.Mock.Results = md.Mock.Results[1:]
		deepCopy(data, m.Data)
		return m.Err
	}

	db := md.DB
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	tx := db.Where(map[string]interface{}{k: v}).First(data)
	return tx.Error
}

// Query 查找用户
func (md *MysqlDB) Query(data any, query any, args ...any) error {
	//mock
	if md.Mock != nil && len(md.Mock.Results) > 0 {
		m := md.Mock.Results[0]
		md.Mock.Results = md.Mock.Results[1:]
		deepCopy(data, m.Data)
		return m.Err
	}

	db := md.DB
	tx := db.Where(query, args...).First(data)
	return tx.Error
}
func (md *MysqlDB) QueryAll(data interface{}, order string, limit int, query interface{}, args ...interface{}) error {
	//mock
	if md.Mock != nil && len(md.Mock.Results) > 0 {
		m := md.Mock.Results[0]
		md.Mock.Results = md.Mock.Results[1:]
		deepCopy(data, m.Data)
		return m.Err
	}

	Db := md.DB
	if query != nil {
		Db = Db.Where(query, args...)
	}
	if len(order) > 0 {
		Db = Db.Order(order)
	}
	if limit > 0 {
		Db = Db.Limit(limit)
	}
	return Db.Find(data).Error
}

// Save 更新用户，fields空时，数据不存在会写入。非空时，不存在会更新失败
func (md *MysqlDB) Save(data interface{}, fields ...string) error {
	//mock
	if md.Mock != nil && len(md.Mock.Results) > 0 {
		m := md.Mock.Results[0]
		md.Mock.Results = md.Mock.Results[1:]
		return m.Err
	}

	if len(fields) == 0 {
		return md.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(data).Error
	}
	return md.DB.Model(data).Select(fields).Updates(data).Error
}

// Save 更新用户，fields空时，数据不存在会写入。非空时，不存在会更新失败
func (md *MysqlDB) SaveByWhere(data interface{}, k string, v interface{}, fields ...string) error {
	//mock
	if md.Mock != nil && len(md.Mock.Results) > 0 {
		m := md.Mock.Results[0]
		md.Mock.Results = md.Mock.Results[1:]
		return m.Err
	}

	return md.DB.Model(data).Where(k, v).Updates(data).Error
}

// DeleteOne 删除用户
func (md *MysqlDB) Delete(data any, query any, args ...any) error {
	if md.Mock != nil && len(md.Mock.Results) > 0 {
		m := md.Mock.Results[0]
		md.Mock.Results = md.Mock.Results[1:]
		return m.Err
	}

	db := md.DB.Unscoped()
	if query != nil {
		db = db.Where(query, args...)
	}

	tx := db.Delete(data)
	return tx.Error
}

// QueryList 分页查询
func (md *MysqlDB) QueryList(data interface{}, order string, limit, index int, query interface{}, args ...interface{}) (int64, error) {
	//mock
	if md.Mock != nil && len(md.Mock.Results) > 0 {
		m := md.Mock.Results[0]
		md.Mock.Results = md.Mock.Results[1:]
		deepCopy(data, m.Data)
		return m.Count, m.Err
	}

	var (
		r  int64
		Db = md.DB.Model(data)
	)

	if query != nil {
		Db = Db.Where(query, args...)
	}
	Db = Db.Count(&r)

	if len(order) > 0 {
		Db = Db.Order(order)
	}

	tx := Db.Limit(limit).Offset((index - 1) * limit).Find(data)
	if err := tx.Error; err != nil {
		return 0, err
	}
	return r, nil
}

// Count 查询数量
func (md *MysqlDB) Count(data interface{}, query interface{}, args ...interface{}) (int64, error) {
	//mock
	if md.Mock != nil && len(md.Mock.Results) > 0 {
		m := md.Mock.Results[0]
		md.Mock.Results = md.Mock.Results[1:]
		return m.Count, m.Err
	}

	Db := md.DB.Model(data)
	if query != nil {
		Db = Db.Where(query, args...)
	}

	var r int64
	err := Db.Count(&r).Error
	return r, err
}

func Unmarshal(p *string, data interface{}) error {
	if len(*p) == 0 {
		return nil
	}

	err := json.Unmarshal([]byte(*p), data)
	if err != nil {
		return err
	}
	*p = ""
	return nil
}

func Marshal(data interface{}, field string, fields ...string) string {
	if len(fields) > 0 {
		ok := false
		for i := 0; i < len(fields); i++ {
			if field == fields[i] {
				ok = true
				break
			}
		}
		if !ok {
			return ""
		}
	}

	buff, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(buff)
}
