package model

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

// OauthApplication 应用表
type OauthApplication struct {
	common.BaseModel
	Name         string `json:"name" gorm:"size:100;comment:应用名称"`
	Domain       string `json:"domain" gorm:"size:100;comment:域名"`
	RedirectUri  string `json:"redirect_uri" gorm:"size:100;comment:回调地址"`
	Scopes       string `json:"scopes" gorm:"size:100;comment:授权范围"`
	CreatedAt    int64  `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt    int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	Website      string `json:"website" gorm:"size:100;comment:网站"`
	ClientID     string `json:"client_id" gorm:"size:100;comment:客户端ID"`
	ClientSecret string `json:"client_secret" gorm:"size:100;comment:客户端密钥"`
}

func (*OauthApplication) TableComment() string {
	return `应用表`
}

func (o *OauthApplication) AddRequired() string {
	if o.Name == "" || o.Domain == "" || o.RedirectUri == "" {
		return "name, domain, redirect_uri are required"
	}
	return ""
}

func (o *OauthApplication) CheckDuplicatedModel() database.Model {
	return &OauthApplication{
		Name:   o.Name,
		Domain: o.Domain,
	}
}

func (o *OauthApplication) UpdateModel() database.Model {
	return &OauthApplication{
		Name:         o.Name,
		Domain:       o.Domain,
		RedirectUri:  o.RedirectUri,
		Scopes:       o.Scopes,
		Website:      o.Website,
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
	}
}

func (o *OauthApplication) FuzzyQueryMap() map[string]string {
	result := make(map[string]string)
	if o.Name != "" {
		result["name"] = o.Name
	}
	return result
}

func (o *OauthApplication) ListOmits() string {
	return "clientSecret"
}

func init() {
	database.Models = append(database.Models, &OauthApplication{})
}
