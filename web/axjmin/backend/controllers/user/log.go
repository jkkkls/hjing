package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"webserver/common"
	dbadmin "webserver/db/admin"

	"github.com/gin-gonic/gin"
)

func getLog(c *gin.Context) {
	current, _ := strconv.Atoi(c.Query("current"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	operate := c.Query("operate")
	name := c.Query("name")

	var (
		where []string
		args  []any
		logs  []*dbadmin.Log
		count int64
		err   error
	)
	if name != "" {
		where = append(where, "name = ?")
		args = append(args, name)
	}

	if operate != "" {
		where = append(where, "operate like ?")
		args = append(args, fmt.Sprintf("%%%v%%", operate))
	}

	if len(args) == 0 {
		logs, count, err = dbadmin.QueryLogList(pageSize, current, nil)
	} else {
		logs, count, err = dbadmin.QueryLogList(pageSize, current, strings.Join(where, " and "), args...)
	}
	res := &common.ResData{
		Succ:  true,
		Code:  0,
		Count: int64(count),
	}
	if err != nil {
		log.Println("000", err)
	}
	for i := 0; i < len(logs); i++ {
		res.Data = append(res.Data, logs[i])
	}
	c.JSON(http.StatusOK, res)
}
