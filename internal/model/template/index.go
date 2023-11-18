package model

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

type Template struct {
	common.BaseModel
	Name       string `json:"name" gorm:"size:100:comment:模板名称"`
	SiteID     uint64 `json:"siteId" gorm:"index;comment:站点ID"`
	Content    string `json:"content" gorm:"comment:模板内容"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (*Template) TableComment() string {
	return `模板表`
}

func (t *Template) AddRequired() string {
	if t.Name == "" || t.SiteID == 0 {
		return "name, siteId"
	}
	return ""
}

func (t *Template) CheckDuplicatedModel() database.Model {
	return &Template{
		Name:   t.Name,
		SiteID: t.SiteID,
	}
}

func (t *Template) UpdateModel() database.Model {
	return &Template{
		Name:    t.Name,
		SiteID:  t.SiteID,
		Content: t.Content,
	}
}

func (t *Template) FuzzyQueryMap() map[string]string {
	result := make(map[string]string)
	if t.Name != "" {
		result["name"] = "%" + t.Name + "%"
	}
	return result
}

func (t *Template) ExactMatchModel() database.Model {
	return &Template{
		SiteID: t.SiteID,
	}
}

func (t *Template) ListOmits() string {
	return "content"
}

func init() {
	database.Models = append(database.Models, &Template{})
}
