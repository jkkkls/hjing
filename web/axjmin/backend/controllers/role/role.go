package role

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"webserver/common"
	dbadmin "webserver/db/admin"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func getRoles(c *gin.Context) {
	current, _ := strconv.Atoi(c.Query("current"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	id := c.Query("id")
	name := c.Query("name")

	if current < 1 {
		current = 1
	}
	if pageSize < 1 {
		pageSize = 100
	}

	var where []string
	var args []any
	var roles []*dbadmin.SysRole
	var count int64
	var err error

	if id != "" {
		where = append(where, "id = ?")
		args = append(args, id)
	}

	if name != "" {
		where = append(where, "name like ?")
		args = append(args, fmt.Sprintf("%%%v%%", name))
	}

	where = append(where, "id <> ?")
	args = append(args, 1)

	if len(args) == 0 {
		count, err = dbadmin.QueryList(&roles, "", "", pageSize, current, nil)
	} else {
		count, err = dbadmin.QueryList(&roles, "", "", pageSize, current, strings.Join(where, " and "), args...)
	}
	res := &common.ResData{
		Succ:  true,
		Code:  0,
		Count: int64(count),
	}
	if err != nil {
		log.Println("000", err)
	}
	for i := 0; i < len(roles); i++ {
		// if roles[i].ID == "admin" {
		// 	roles[i].Selected = models.GetAllRoutes()
		// }
		res.Data = append(res.Data, roles[i])
	}
	c.JSON(http.StatusOK, res)
}

func updateRole(c *gin.Context) {
	r := &dbadmin.SysRole{}
	err := c.BindJSON(r)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}

	if r.Id == 1 {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": "超级管理员不能修改",
		})
		return
	}

	now := time.Now()
	username := c.Request.Header.Get("admin-id")
	if r.Id == 0 {
		r.CreateTs = &now
		r.UpdateTs = &now
		r.CreateBy = username
	} else {
		old := &dbadmin.SysRole{}
		err = dbadmin.Query(old, "id = ?", r.Id)
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"success":      false,
				"errorMessage": err.Error(),
			})
			return
		}
		r.CreateTs = old.CreateTs
		r.CreateBy = old.CreateBy
		r.UpdateTs = &now
	}
	dbadmin.Save(r)

	c.JSON(http.StatusOK, map[string]any{
		"success":      true,
		"errorMessage": "操作成功",
	})
}

func deleteRole(c *gin.Context) {
	buff, _ := ioutil.ReadAll(c.Request.Body)
	id, _ := jsonparser.GetInt(buff, "id")

	if id == 1 {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": "该角色不能删除",
		})
		return
	}

	dbadmin.AssociationDelete(&dbadmin.SysRole{Id: id})

	c.JSON(http.StatusOK, map[string]any{
		"success":      true,
		"errorMessage": "操作成功",
	})
}
