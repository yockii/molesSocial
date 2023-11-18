package service

import (
	"errors"
	model "github.com/yockii/molesSocial/internal/model/oauth"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/util"
)

var OauthApplicationService = newOauthApplicationService()

func newOauthApplicationService() *oauthApplicationService {
	s := new(oauthApplicationService)
	s.BaseService = common.BaseService[*model.OauthApplication]{
		Service: s,
	}
	return s
}

type oauthApplicationService struct {
	common.BaseService[*model.OauthApplication]
}

func (*oauthApplicationService) Model() *model.OauthApplication {
	return new(model.OauthApplication)
}

func (s *oauthApplicationService) Save(d *model.OauthApplication) (*model.OauthApplication, error) {
	if d.Name == "" || d.Website == "" || d.RedirectUri == "" {
		return nil, errors.New("client_name, website, redirect_uris are required")
	}

	// 查重？  不查重，直接生成，虽然可能会有很多记录
	//instance := new(model.OauthApplication)
	//if err := database.DB.Model(s.Model()).Where("name = ?", d.Name).First(instance).Error; err != nil {
	//	if !errors.Is(err, gorm.ErrRecordNotFound) {
	//		logger.Errorln(err)
	//		return nil, err
	//	} else {
	//		return instance, nil
	//	}
	//}

	if d.Scopes == "" {
		d.Scopes = "read"
	}
	d.ClientID = util.GenerateXid()
	d.ClientSecret = util.GenerateDatabaseID()
	d.ID = util.SnowflakeId()

	if err := database.DB.Create(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}
