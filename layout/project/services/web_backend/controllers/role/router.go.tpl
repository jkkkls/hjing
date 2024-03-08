package role

import "github.com/gin-gonic/gin"

func Routers(e *gin.Engine) {
	e.GET("/api/role", getRoles)
	e.POST("/api/role", updateRole)
	e.DELETE("/api/role", deleteRole)

	e.GET("/api/permission", getPermissions)
	e.POST("/api/permission", updatePermission)

	e.GET("/api/role_permission", getRolePermissions)
	e.POST("/api/role_permission", updateRolePermission)
}
