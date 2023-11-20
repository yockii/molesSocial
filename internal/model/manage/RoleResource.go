package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

type RoleResource struct {
	common.BaseModel
	RoleID     uint64 `json:"roleId" gorm:"index;comment:角色ID"`
	ResourceID uint64 `json:"resourceId" gorm:"index;comment:资源ID"`
	CreatedAt  int64  `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
}

func (*RoleResource) TableComment() string {
	return "角色资源表"
}

func (rr *RoleResource) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	rr.ID = j.Get("id").Uint()
	rr.RoleID = j.Get("roleId").Uint()
	return nil
}

func init() {
	database.Models = append(database.Models, &RoleResource{})
}
