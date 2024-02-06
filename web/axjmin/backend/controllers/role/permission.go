package role

import (
	"log"
	"net/http"
	"strconv"
	"webserver/common"
	dbadmin "webserver/db/admin"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getPermissions(c *gin.Context) {
	var ps, menusData []*dbadmin.SysPermission
	err := dbadmin.QueryAll(&ps, "", nil)
	res := &common.ResData{
		Succ:  true,
		Code:  0,
		Count: int64(len(ps)),
	}
	if err != nil {
		res.Succ = false
		res.ErrorMessage = err.Error()
		return
	}

	m := map[int64]*dbadmin.SysPermission{}
	for _, v := range ps {
		m[v.Id] = v
	}

	for _, v := range m {
		if v.ParentId == 0 {
			menusData = append(menusData, v)
		} else {
			if p, ok := m[v.ParentId]; ok {
				p.Children = append(p.Children, v)
			}
		}
	}

	dbadmin.SortUserPermission(menusData)

	if err != nil {
		log.Println("000", err)
	}
	for i := 0; i < len(menusData); i++ {

		res.Data = append(res.Data, menusData[i])
	}
	c.JSON(http.StatusOK, res)
}

func updatePermission(c *gin.Context) {
	r := &dbadmin.SysPermission{}
	err := c.BindJSON(r)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}

	old := &dbadmin.SysPermission{}
	err = dbadmin.Query(old, "id = ?", r.Id)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}

	dbadmin.Save(r)

	c.JSON(http.StatusOK, map[string]any{
		"success":      true,
		"errorMessage": "操作成功",
	})
}

// getRolePermissions 函数用于获取指定角色的权限
func getRolePermissions(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)

	var ps, menusData []*dbadmin.SysPermission
	err := dbadmin.QueryAll(&ps, "", nil)
	res := &common.ResData{
		Succ:  true,
		Code:  0,
		Count: int64(len(ps)),
	}
	if err != nil {
		res.Succ = false
		res.ErrorMessage = err.Error()
		return
	}

	rps := []*dbadmin.SysRolePermission{}

	err = dbadmin.QueryAll(&rps, "", "sys_role_id = ?", id)
	if err != nil {
		log.Println(err)
	}

	m := map[int64]*dbadmin.SysPermission{}
	for _, v := range ps {
		if v.Status == "hide" {
			continue
		}
		m[v.Id] = v
	}
	for _, v := range rps {
		if p, ok := m[v.SysPermissionID]; ok {
			p.Code = v.Code
		}
	}

	for _, v := range m {
		if v.ParentId == 0 {
			menusData = append(menusData, v)
		} else {
			if p, ok := m[v.ParentId]; ok {
				p.Children = append(p.Children, v)
			}
		}
	}

	dbadmin.SortUserPermission(menusData)

	if err != nil {
		log.Println("000", err)
	}
	for i := 0; i < len(menusData); i++ {
		res.Data = append(res.Data, menusData[i])
	}
	c.JSON(http.StatusOK, res)
}

// updateRolePermission 用于更新角色权限信息
func updateRolePermission(c *gin.Context) {
	r := &dbadmin.SysPermission{}
	err := c.BindJSON(r)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}

	old := &dbadmin.SysRolePermission{}
	err = dbadmin.Query(old, "sys_role_id = ? and sys_permission_id = ?", r.SysRoleID, r.Id)
	if err == gorm.ErrRecordNotFound {
		old.SysRoleID = r.SysRoleID
		old.SysPermissionID = r.Id
	} else if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	old.Code = r.Code

	if old.Id == 0 {
		err = dbadmin.Save(old)
	} else {
		err = dbadmin.Save(old, "code")
	}
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"success":      true,
		"errorMessage": "操作成功",
	})
}
