package models

import (
	"slices"
	dbadmin "webserver/db/admin"
)

// passPaths 无需鉴权路径
var passPaths = []string{
	"/api/currentUser",    //
	"/api/login/outLogin", //
	"/api/all_users",
	"/api/all_roles",
}

// TreeData 菜单数据
var TreeData = []*dbadmin.SysPermission{
	{
		Name: "系统管理",
		Path: "/system",
		Icon: "setting",
		Children: []*dbadmin.SysPermission{
			{
				Name: "用户管理",
				Path: "/system/user",
				Apis: []string{"/api/users", "/api/user"},
			},
			{
				Name: "角色管理",
				Path: "/system/role",
				Apis: []string{"/api/role", "/api/role_permission"},
			},
			{
				Name: "权限管理",
				Path: "/system/permission",
				Apis: []string{"/api/permission"},
			},
			{
				Name: "后台日志",
				Path: "/system/log",
				Apis: []string{"/api/log"},
			},
		},
	},
}

func queryMenus(parent []*dbadmin.SysPermission, selects []string) []*dbadmin.SysPermission {
	var ret []*dbadmin.SysPermission

	for _, v := range parent {
		if slices.Contains(selects, v.Path) {
			d := &dbadmin.SysPermission{
				Name: v.Name,
				Path: v.Path,
				Icon: v.Icon,
			}
			if len(v.Children) > 0 {
				d.Children = queryMenus(v.Children, selects)
			}
			ret = append(ret, d)
		}
	}

	return ret
}

func QueryMenus(selects []string) []*dbadmin.SysPermission {
	return queryMenus(TreeData, selects)
}

func getChildren(d *dbadmin.SysPermission, arr *[]string) {
	(*arr) = append((*arr), d.Path)
	for _, v := range d.Children {
		getChildren(v, arr)
	}
}

func GetAllRoutes() []string {
	var arr []string

	for _, v := range TreeData {
		getChildren(v, &arr)
	}

	return arr
}

func CheckPermission(user *dbadmin.SysUser, path, method string) bool {
	//超级管理员
	for _, v := range user.SysRoleUsers {
		if v.SysRoleID == 1 {
			return true
		}
	}

	if slices.Contains(passPaths, path) {
		return true
	}

	// buff, _ := json.MarshalIndent(user.SysRoleUsers, "", "  ")
	// log.Println(string(buff))
	for _, v := range user.SysRoleUsers {
		for _, rp := range v.SysRole.SysRolePermissions {
			// dump.Println(rp, user.Username, path, method)
			if slices.Contains(rp.SysPermission.Apis, path) && slices.Contains(rp.Code, method) {
				return true
			}
		}
	}

	return false
}
