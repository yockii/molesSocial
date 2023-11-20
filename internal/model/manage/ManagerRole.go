package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

type ManagerRole struct {
	common.BaseModel
	ManagerID uint64 `json:"managerId" gorm:"index;comment:管理员ID"`
	RoleID    uint64 `json:"roleId" gorm:"index;comment:角色ID"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
}

func (_ *ManagerRole) TableComment() string {
	return "管理员角色表"
}

func (ur *ManagerRole) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ur.ID = j.Get("id").Uint()
	ur.ManagerID = j.Get("managerId").Uint()
	ur.RoleID = j.Get("roleId").Uint()
	return nil
}

func init() {
	database.Models = append(database.Models, &ManagerRole{})
}
