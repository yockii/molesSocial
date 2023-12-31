package model

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

// Site 站点表
type Site struct {
	common.BaseModel
	Name             string `json:"name" gorm:"size:100;comment:站点名称"`
	Domain           string `json:"domain" gorm:"size:100;comment:站点域名"`
	Description      string `json:"description" gorm:"size:500;comment:站点描述"`
	ShortDesc        string `json:"shortDesc" gorm:"size:100;comment:站点简介"`
	Thumbnail        string `json:"thumbnail" gorm:"size:100;comment:站点缩略图"`
	Email            string `json:"email" gorm:"size:100;comment:站点邮箱"`
	Languages        string `json:"languages" gorm:"size:100;comment:站点语言"`
	Registrations    *bool  `json:"registrations" gorm:"comment:是否开放注册"`
	ApprovalRequired *bool  `json:"approvalRequired" gorm:"comment:是否需要审核"`
	InvitesEnabled   *bool  `json:"invitesEnabled" gorm:"comment:是否开放邀请"`
	Configuration    string `json:"configuration" gorm:"size:2000;default:{};comment:站点配置"`
	Rules            string `json:"rules" gorm:"size:2000;comment:站点规则"`
	CreateTime       int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (*Site) TableComment() string {
	return `站点表`
}

func (s *Site) AddRequired() string {
	if s.Name == "" || s.Domain == "" {
		return "name, code, host"
	}
	return ""
}

func (s *Site) CheckDuplicatedModel() database.Model {
	return &Site{
		Name: s.Name,
	}
}

func (s *Site) UpdateModel() database.Model {
	return &Site{
		Name:             s.Name,
		Domain:           s.Domain,
		Description:      s.Description,
		ShortDesc:        s.ShortDesc,
		Thumbnail:        s.Thumbnail,
		Email:            s.Email,
		Languages:        s.Languages,
		Registrations:    s.ApprovalRequired,
		ApprovalRequired: s.ApprovalRequired,
		InvitesEnabled:   s.InvitesEnabled,
		Configuration:    s.Configuration,
		Rules:            s.Rules,
	}
}

func (s *Site) FuzzyQueryMap() map[string]string {
	result := make(map[string]string)
	if s.Name != "" {
		result["name"] = s.Name
	}
	if s.Domain != "" {
		result["host"] = s.Domain
	}
	return result
}
func (s *Site) ListOmits() string {
	return "description"
}

func init() {
	database.Models = append(database.Models, &Site{})
}
