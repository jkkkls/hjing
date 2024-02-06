package user

import "github.com/gin-gonic/gin"

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
}
