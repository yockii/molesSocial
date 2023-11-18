package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/model/account"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
	"time"
)

var AccountService = newAccountService()

type accountService struct {
	common.BaseService[*model.Account]
}

func newAccountService() *accountService {
	s := new(accountService)
	s.BaseService = common.BaseService[*model.Account]{
		Service: s,
	}
	return s
}

func (*accountService) Model() *model.Account {
	return new(model.Account)
}

func (s *accountService) GetByUsernameAndDomain(username string, domain string) (*model.Account, error) {
	instance := new(model.Account)
	err := database.DB.Where(&model.Account{Username: username, Domain: domain}).First(instance).Error
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return instance, nil
}

func (s *accountService) CountUsers(domain string) (total, monthly, halfYears int64, err error) {
	err = database.DB.Model(s.Model()).Where("domain = ?", domain).Count(&total).Error
	if err != nil {
		logger.Errorln("CountUsers", err)
		return
	}
	err = database.DB.Model(s.Model()).Where("domain = ? AND last_sign_in_at > ?", domain, time.Now().UnixMilli()-constant.MonthMilliseconds).Count(&monthly).Error
	if err != nil {
		logger.Errorln("CountUserMonthly", err)
		return
	}
	err = database.DB.Model(s.Model()).Where("domain = ? AND last_sign_in_at > ?", domain, time.Now().UnixMilli()-constant.HalfYearMilliseconds).Count(&halfYears).Error
	if err != nil {
		logger.Errorln("CountUserHalfYears", err)
		return
	}
	return
}
