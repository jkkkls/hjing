package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"webserver/controllers/role"
	"webserver/controllers/user"
	dbadmin "webserver/db/admin"
	"webserver/models"

	"gitee.com/jkkkls/huaijing/libs/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options = []Option{}

// 注册app的路由配置
func addRouter(opts ...Option) {
	options = append(options, opts...)
}

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		pass := []string{"/api/login/account", "/api/shared/img"}
		if utils.In(pass, path) {
			c.Next()
			return
		}

		var username string
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(user.JwtSecret), nil
			})
		if err == nil {
			if token.Valid {
				info := token.Claims.(jwt.MapClaims)
				if v, ok := info["id"]; ok {
					username = v.(string)
					c.Request.Header.Set("admin-id", username)
				}
			} else {
				c.JSON(http.StatusUnauthorized, map[string]any{
					"errorCode":    401,
					"errorMessage": "请先登录1！",
					"success":      true,
					"data": map[string]any{
						"isLogin": false,
					},
				})
				c.Abort()
				return
			}

			//检查权限
			//username
			//path
			user, _ := dbadmin.QuerySysUser(username, true)
			if user == nil {
				c.JSON(http.StatusUnauthorized, map[string]any{
					"errorCode":    401,
					"errorMessage": "请先登录2！",
					"success":      true,
					"data": map[string]any{
						"isLogin": false,
					},
				})
				c.Abort()
				return
			}

			m := map[string]string{
				"GET":   "R",
				"PUT":   "U",
				"POST":  "U",
				"DELTE": "D",
			}
			if !models.CheckPermission(user, path, m[c.Request.Method]) {
				c.JSON(http.StatusOK, map[string]any{
					"errorMessage": "没有相应权限，请联系管理员!",
					"success":      false,
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, map[string]any{
				"errorCode":    401,
				"errorMessage": "请先登录4！",
				"success":      true,
				"data": map[string]any{
					"isLogin": false,
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 初始化
func Init() *gin.Engine {
	r := gin.New()
	r.Use(MiddleWare())
	for _, opt := range options {
		opt(r)
	}
	return r
}

func initDb() {
	dbadmin.UpdateSysPermissions(models.TreeData)

	roles, _ := dbadmin.QueryAllSysRole()
	if len(roles) > 0 {
		return
	}

	now := time.Now()
	role := &dbadmin.SysRole{
		Name:     "超级管理员",
		Desc:     "拥有最高权限",
		CreateTs: &now,
		UpdateTs: &now,
	}

	err := dbadmin.Save(role)
	if err != nil {
		fmt.Println(err)
		return
	}

	user := &dbadmin.SysUser{
		Username: "admin123",
		Name:     "羊过",
		Password: utils.Md5("admin123" + "admin123"), //用户名和密码都是admin123
		// Type:        1,
		Status:   "normal",
		CreateTs: &now,
	}
	err = dbadmin.Save(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dbadmin.Save(&dbadmin.SysRoleUser{
		SysUserID: user.Id,
		SysRoleID: role.Id,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	err := models.LoadRunConf("./server.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dbadmin.InitAdminDB(models.Conf.Db.Type, models.Conf.Db.Dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.WriteFile("./web.pid", []byte(strconv.Itoa(os.Getpid())), 0666)

	initDb()

	addRouter(
		user.Routers,
		role.Routers,
	)
	r := Init()

	log.Println("server start")

	r.Run(models.Conf.Net.Address)
}
