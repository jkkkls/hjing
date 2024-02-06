package user

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"webserver/common"
	dbadmin "webserver/db/admin"

	"gitee.com/jkkkls/huaijing/libs/utils"

	"github.com/buger/jsonparser"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	JwtSecret = "admin123123213"
)

func currentUser(c *gin.Context) {
	id := c.Request.Header.Get("admin-id")
	if id == "" {
		c.JSON(http.StatusUnauthorized, map[string]any{
			"errorCode":    401,
			"errorMessage": "请先登录1！",
			"success":      true,
			"data": map[string]any{
				"isLogin": false,
			},
		})
		return
	}

	user, _ := dbadmin.QuerySysUser(id, true)
	if user == nil {
		c.JSON(http.StatusUnauthorized, map[string]any{
			"errorCode":    401,
			"errorMessage": "请先登录2！",
			"success":      true,
			"data": map[string]any{
				"isLogin": false,
			},
		})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    user,
	})
}

func updateUser(c *gin.Context) {
	data := &dbadmin.SysUser{}
	err := c.BindJSON(data)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}

	if data.Username == "" {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": "用户名不能为空",
		})
		return
	}

	if slices.Contains(data.SysRoleIDs, 1) {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": "操作不允许",
		})
	}

	old := &dbadmin.SysUser{}
	err = dbadmin.Query(old, "username = ?", data.Username)
	if err == nil {
		if old.Id != data.Id {
			c.JSON(http.StatusOK, map[string]any{
				"success":      false,
				"errorMessage": "用户名已存在",
			})
			return
		}
	}

	now := time.Now()
	user := &dbadmin.SysUser{
		Id:       data.Id,
		Username: data.Username,
		Status:   data.Status,
		Type:     data.Type,
		CreateTs: &now,
	}
	if data.Id != 0 {
		old = &dbadmin.SysUser{}
		err := dbadmin.Query(old, "id = ?", data.Id)
		if err == nil {
			user.CreateTs = old.CreateTs
			user.Password = old.Password
			user.Name = old.Name
		}
	}

	if data.Name != "" {
		user.Name = data.Name
	}
	if data.Password != "" {
		user.Password = utils.Md5(user.Username + data.Password)
	}

	err = dbadmin.Save(user)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": "保存用户失败，原因：" + err.Error(),
		})
		return
	}

	//更新权限
	roles := []dbadmin.SysRoleUser{}
	err = dbadmin.QueryAll(&roles, "", "sys_user_id = ?", data.Id)
	if err != nil {
		log.Println(err)
	}

	var dels []int64
	var adds []*dbadmin.SysRoleUser
	for _, v := range roles {
		if !slices.Contains(data.SysRoleIDs, v.SysRoleID) {
			dels = append(dels, v.Id)
		}
	}
	for _, v := range data.SysRoleIDs {
		if !slices.ContainsFunc(roles, func(sru dbadmin.SysRoleUser) bool {
			return sru.SysRoleID == v
		}) {
			adds = append(adds, &dbadmin.SysRoleUser{
				SysRoleID: v,
				SysUserID: user.Id,
			})
		}
	}
	// log.Println(adds, dels)
	if len(adds) > 0 {
		err = dbadmin.Save(adds)
		if err != nil {
			log.Println(err)
		}
	}
	if len(dels) > 0 {
		err = dbadmin.Delete(&dbadmin.SysRoleUser{}, "id in (?)", dels)
		if err != nil {
			log.Println(err)
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"success":      true,
		"errorMessage": "操作成功",
	})
}

func deleteUser(c *gin.Context) {
	buff, _ := ioutil.ReadAll(c.Request.Body)
	id, _ := jsonparser.GetInt(buff, "id")
	if id == 1 {
		c.JSON(http.StatusOK, map[string]any{
			"success":      false,
			"errorMessage": "该用户不能删除",
		})
		return
	}

	dbadmin.AssociationDelete(&dbadmin.SysUser{Id: id})

	c.JSON(http.StatusOK, map[string]any{
		"success":      true,
		"errorMessage": "操作成功",
	})
}

func getAllRoles(c *gin.Context) {
	roles := []dbadmin.SysRole{}
	err := dbadmin.QueryAll(&roles, "", "id <> 1")
	res := &common.ResData{
		Succ:  true,
		Code:  0,
		Count: int64(len(roles)),
	}
	if err != nil {
		log.Println("000", err)
	}

	for i := 0; i < len(roles); i++ {
		res.Data = append(res.Data, &dbadmin.SysRole{
			Id:   roles[i].Id,
			Name: roles[i].Name,
		})
	}
	c.JSON(http.StatusOK, res)
}

func getAllUsers(c *gin.Context) {
	users := []dbadmin.SysUser{}
	err := dbadmin.QueryAll(&users, "", nil)
	res := &common.ResData{
		Succ:  true,
		Code:  0,
		Count: int64(len(users)),
	}
	if err != nil {
		log.Println("000", err)
	}

	for i := 0; i < len(users); i++ {
		res.Data = append(res.Data, &dbadmin.SysUser{
			Id:   users[i].Id,
			Name: users[i].Name,
		})
	}
	c.JSON(http.StatusOK, res)
}

func getUsers(c *gin.Context) {
	current, _ := strconv.Atoi(c.Query("current"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	id := c.Query("id")
	name := c.Query("name")
	role := c.Query("role")

	var where []string
	var args []any
	var users []*dbadmin.SysUser
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

	if role != "" {
		where = append(where, "role = ?")
		args = append(args, role)
	}

	//隐藏
	// where = append(where, "id <> ?")
	// args = append(args, 1)

	if len(args) == 0 {
		count, err = dbadmin.QueryList(&users, "SysRoleUsers", "", pageSize, current, nil)
	} else {
		count, err = dbadmin.QueryList(&users, "SysRoleUsers", "", pageSize, current, strings.Join(where, " and "), args...)
	}
	res := &common.ResData{
		Succ:  true,
		Code:  0,
		Count: int64(count),
	}
	if err != nil {
		log.Println("000", err)
	}
	for i := 0; i < len(users); i++ {
		users[i].Password = ""
		for _, v := range users[i].SysRoleUsers {
			users[i].SysRoleIDs = append(users[i].SysRoleIDs, v.SysRoleID)
		}
		res.Data = append(res.Data, users[i])
	}
	c.JSON(http.StatusOK, res)
}

type Login struct {
	Username string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	Pssword  string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
	T        string `form:"type" json:"type" uri:"type" xml:"type" binding:"required"`
}

func login(c *gin.Context) {
	var data Login
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"status":           err.Error(),
			"type":             data.T,
			"currentAuthority": "guest",
		})
		return
	}

	id := c.Request.Header.Get("admin-id")
	if id != "" {
		user := &dbadmin.SysUser{}
		err := dbadmin.Query(user, "username = ?", id)
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"status":           "ok",
				"type":             data.T,
				"currentAuthority": "admin",
			})
			return
		}
	}

	user := &dbadmin.SysUser{}
	err := dbadmin.Query(user, "username = ?", data.Username)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"status":           "用户不存在",
			"type":             data.T,
			"currentAuthority": "guest",
		})
		return
	}
	if user.Password != utils.Md5(data.Username+data.Pssword) {
		c.JSON(http.StatusOK, map[string]any{
			"status":           "密码错误",
			"type":             data.T,
			"currentAuthority": "guest",
		})
		return
	}

	now := time.Now()
	user.LastLoginTs = &now
	dbadmin.Save(user, "last_login_ts")

	dbadmin.SaveLog(data.Username, "登录", map[string]any{
		"time": now,
		"ip":   c.ClientIP(),
	})

	c.JSON(http.StatusOK, map[string]any{
		"status":           "ok",
		"type":             data.T,
		"currentAuthority": "admin",
		"token":            newToken(data.Username),
	})
}

func logout(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]any{
		"success": "true",
	})
}

func newToken(id string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = id
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return ""
	}

	return tokenString
}
