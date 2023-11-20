package model

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

const (
	ManagerStatusNormal = 1
	ManagerStatusBan    = -1
)

type Manager struct {
	common.BaseModel
	Username  string `json:"username" gorm:"size:100;comment:用户名"`
	Password  string `json:"password" gorm:"size:100;comment:密码"`
	Email     string `json:"email" gorm:"size:100;comment:邮箱"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime;comment:更新时间"`
	Status    int    `json:"status" gorm:"comment:状态 1-启用 -1-禁用"`
}

func (_ *Manager) TableComment() string {
	return "管理员表"
}

func (m *Manager) AddRequired() string {
	if m.Username == "" || m.Password == "" || m.Email == "" {
		return "username,password,email"
	}
	return ""
}

func (m *Manager) UpdateModel() database.Model {
	return &Manager{
		Username: m.Username,
		Password: m.Password,
		Email:    m.Email,
	}
}

func (m *Manager) FuzzyQueryMap() map[string]string {
	result := make(map[string]string)
	if m.Username != "" {
		result["username"] = "%" + m.Username + "%"
	}
	if m.Email != "" {
		result["email"] = "%" + m.Email + "%"
	}
	return result
}

func init() {
	database.Models = append(database.Models, &Manager{})
}
