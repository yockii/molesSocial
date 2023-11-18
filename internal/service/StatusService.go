package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/molesSocial/internal/model/status"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

var StatusService = newStatusService()

type statusService struct {
	common.BaseService[*model.Status]
}

func newStatusService() *statusService {
	s := new(statusService)
	s.BaseService = common.BaseService[*model.Status]{
		Service: s,
	}
	return s
}

func (*statusService) Model() *model.Status {
	return new(model.Status)
}

func (s *statusService) CountLocalStatuses(domain string) (total int64, err error) {
	err = database.DB.Model(s.Model()).Where("domain = ?", domain).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
