package dbadmin

import (
	"slices"
	"time"

	"gitee.com/jkkkls/huaijing/libs/utils"
)

// SysUser 用户
type SysUser struct {
	Id          int64      `gorm:"primarykey" json:"id"`
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	Name        string     `json:"name"`
	Type        int        `json:"type"`
	CreateTs    *time.Time `json:"createTs"`
	LastLoginTs *time.Time `json:"lastLoginTs"`
	Status      string     `json:"status"`

	MenuData   []*SysPermission `gorm:"-" json:"menuData"`
	SysRoleIDs []int64          `gorm:"-" json:"sysRoleIDs"`

	SysRoleUsers []*SysRoleUser `json:"-"`
}

// SysRole 角色
type SysRole struct {
	Id                 int64                `gorm:"primarykey" json:"id"`
	Name               string               `json:"name"`
	Desc               string               `json:"desc"`
	CreateTs           *time.Time           `json:"createTs"`
	UpdateTs           *time.Time           `json:"updateTs"`
	CreateBy           string               `json:"createBy"`
	SysRolePermissions []*SysRolePermission `json:"-"`
}

// SysRoleUser 用户的角色数据
type SysRoleUser struct {
	Id        int64    `gorm:"primarykey" json:"id"`
	SysUserID int64    `json:"sysUserID"`
	SysRoleID int64    `json:"sysRoleID"`
	SysRole   *SysRole `json:"-"`
}

// SysRolePermission 角色的权限数据
type SysRolePermission struct {
	Id        int64 `gorm:"primarykey" json:"id"`
	SysRoleID int64 `json:"sysRoleID"`

	SysPermissionID int64          `json:"sysPermissionID"`
	SysPermission   *SysPermission `json:"-"`
	Code            []string       `json:"code" gorm:"serializer:json"` //CRUD
}

// SysPermission 权限
type SysPermission struct {
	Id       int64    `gorm:"primarykey" json:"id"`
	ParentId int64    `json:"parentId"`
	Name     string   `json:"name"`
	Path     string   `gorm:"index" json:"path"`
	Icon     string   `json:"icon"`
	Level    int      `json:"level"`
	Redirect string   `json:"redirect"`
	Status   string   `json:"status"`
	Apis     []string `json:"apis" gorm:"serializer:json"`

	Children []*SysPermission `gorm:"-" json:"children"`

	//角色授权用途
	Code      []string `gorm:"-" json:"code"` //CRUD
	SysRoleID int64    `gorm:"-" json:"sysRoleID"`
}

func init() {
	tables = append(tables, &SysRole{}, &SysUser{}, &SysPermission{}, &SysRolePermission{}, &SysRoleUser{})
}

func QuerySysUser(username string, includeMenu bool) (*SysUser, error) {
	user := &SysUser{}
	err := adminDB.DB.Where("username = ?", username).Preload("SysRoleUsers.SysRole.SysRolePermissions.SysPermission").First(user).Error
	if err != nil {
		return nil, err
	}

	if !includeMenu {
		return user, err
	}

	// permissions := []*SysPermission{}

	isAdmin := false
	for _, roleUser := range user.SysRoleUsers {
		if roleUser.SysRoleID == 1 {
			isAdmin = true
			break
		}
	}

	m := map[int64]*SysPermission{}
	if isAdmin {
		//管理员拥有所有权限
		ps := []*SysPermission{}
		err := adminDB.QueryAll(&ps, "", 0, nil)
		if err != nil {
			utils.Warn("QueryAll SysPermission fail", "err", err)
		}
		for _, v := range ps {
			// 超管可以查看隐藏界面
			// if v.Status == "hide" {
			// 	continue
			// }
			m[v.Id] = v
		}
	} else {
		//查询玩家权限
		for _, roleUser := range user.SysRoleUsers {
			for _, rolePermission := range roleUser.SysRole.SysRolePermissions {
				if rolePermission.SysPermission.Status == "hide" || !slices.Contains(rolePermission.Code, "R") {
					continue
				}

				m[rolePermission.SysPermissionID] = rolePermission.SysPermission
			}
		}
	}
	for _, v := range m {
		if v.ParentId == 0 {
			user.MenuData = append(user.MenuData, v)
		} else {
			if p, ok := m[v.ParentId]; ok {
				p.Children = append(p.Children, v)
			}
		}
	}

	SortUserPermission(user.MenuData)
	return user, nil
}

func SortUserPermission(ps []*SysPermission) {
	slices.SortFunc(ps, func(i, j *SysPermission) int {
		return j.Level - i.Level
	})

	for _, v := range ps {
		if len(v.Children) > 1 {
			SortUserPermission(v.Children)
		}
	}
}

func UpdateSysPermissions(ps []*SysPermission) error {
	//先查询
	old := []*SysPermission{}
	err := adminDB.QueryAll(&old, "", 0, nil)
	if err != nil {
		return err
	}
	m := map[string]*SysPermission{}
	for _, v := range old {
		m[v.Path] = v
	}

	var f func(p *SysPermission)
	f = func(p *SysPermission) {
		v, ok := m[p.Path]
		if !ok {
			p.Status = "normal"
			adminDB.Save(p)
		} else {
			p.Id = v.Id
			p.ParentId = v.ParentId
		}

		for _, c := range p.Children {
			c.ParentId = p.Id
			f(c)
		}
	}

	for _, v := range ps {
		f(v)
	}

	return nil
}
func QueryAllSysRole() (data []*SysRole, err error) {
	err = adminDB.QueryAll(&data, "", 0, nil)
	return
}
