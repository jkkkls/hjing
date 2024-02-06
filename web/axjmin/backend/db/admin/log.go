package dbadmin

import (
	"encoding/json"
	"strings"
	"time"

	"gitee.com/jkkkls/huaijing/libs/utils"
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Name    string     `json:"name,omitempty"`
	Time    *time.Time `json:"time,omitempty"`
	Operate string     `json:"operate,omitempty"`
	Data    string     `json:"data,omitempty"`
}

func init() {
	tables = append(tables, &Log{})
}

// SaveLog ...
func SaveLog(user, operate string, data any) error {
	buff, _ := json.Marshal(data)
	return adminDB.Save(&Log{
		Name:    user,
		Operate: operate,
		Time:    utils.GetNow(),
		Data:    string(buff),
	})
}

func QueryLogList(limit, index int, query interface{}, args ...interface{}) (ls []*Log, n int64, err error) {
	n, err = adminDB.QueryList(&ls, "", "time DESC", limit, index, query, args)
	return
}

func QueryLog(name, operate string, start, end uint64, limit, index int) (ls []*Log, n int64, err error) {
	var (
		arr  []string
		args []interface{}
	)
	if name != "" {
		arr = append(arr, "name = ?")
		args = append(args, name)
	}
	if operate != "" {
		arr = append(arr, "operate = ?")
		args = append(args, operate)
	}
	if start != 0 {
		arr = append(arr, "time >= ?")
		args = append(args, start)
	}
	if end != 0 {
		arr = append(arr, "time <= ?")
		args = append(args, end)
	}

	if len(args) > 0 {
		n, err = adminDB.QueryList(&ls, "", "time DESC", limit, index, strings.Join(arr, " AND "), args...)
	} else {
		n, err = adminDB.QueryList(&ls, "", "time DESC", limit, index, nil)
	}
	return
}
