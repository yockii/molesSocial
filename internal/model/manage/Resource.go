package model

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

type Resource struct {
	common.BaseModel
	Name      string `json:"name" gorm:"size:100;comment:资源名称"`
	Code      string `json:"code" gorm:"size:100;comment:资源编码"`
	Remark    string `json:"remark" gorm:"size:100;comment:备注"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime;comment:更新时间"`
}

func (_ *Resource) TableComment() string {
	return "资源表"
}

func (m *Resource) FuzzyQueryMap() map[string]string {
	result := make(map[string]string)
	if m.Name != "" {
		result["name"] = "%" + m.Name + "%"
	}
	if m.Code != "" {
		result["code"] = "%" + m.Code + "%"
	}
	return result
}

func init() {
	database.Models = append(database.Models, &Resource{})
}
