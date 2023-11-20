package model

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

type Role struct {
	common.BaseModel
	Name        string `json:"name" gorm:"size:100;comment:角色名称"`
	Remark      string `json:"remark" gorm:"size:100;comment:备注"`
	CreatedAt   int64  `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   int64  `json:"updatedAt" gorm:"autoUpdateTime;comment:更新时间"`
	DefaultRole *bool  `json:"defaultRole" gorm:"comment:是否为默认角色"`
}

func (_ *Role) TableComment() string {
	return "角色表"
}

func (m *Role) AddRequired() string {
	if m.Name == "" {
		return "name"
	}
	return ""
}

func (m *Role) UpdateModel() database.Model {
	return &Role{
		Name:   m.Name,
		Remark: m.Remark,
	}
}

func (m *Role) FuzzyQueryMap() map[string]string {
	result := make(map[string]string)
	if m.Name != "" {
		result["name"] = "%" + m.Name + "%"
	}
	return result
}

func init() {
	database.Models = append(database.Models, &Role{})
}
