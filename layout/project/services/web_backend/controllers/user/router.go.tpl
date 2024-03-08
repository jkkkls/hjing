package user

import (
	"net/http"
	"{{projectName}}/services/monitor"
	"{{projectName}}/services/web_backend/common"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	e.GET("/api/users", getUsers)
	e.GET("/api/currentUser", currentUser)
	e.POST("/api/user", updateUser)
	e.DELETE("/api/user", deleteUser)
	e.POST("/api/login/account", login)
	e.POST("/api/login/outLogin", logout)
	// e.GET("/comment", commentHandler)
	e.GET("/api/log", getLog)

	e.GET("/api/all_users", getAllUsers)
	e.GET("/api/all_roles", getAllRoles)

	e.GET("/api/monitor", func(c *gin.Context) {
		ps := monitor.GetProcess()

		res := &common.ResData{
			Succ:  true,
			Code:  0,
			Count: int64(len(ps)),
		}
		for _, v := range ps {
			res.Data = append(res.Data, v)
		}

		c.JSON(http.StatusOK, res)
	})
}
