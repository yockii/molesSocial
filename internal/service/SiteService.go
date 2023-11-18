package service

import (
	"errors"
	"github.com/yockii/molesSocial/internal/model"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
	"gorm.io/gorm"
	"sync"
)

var SiteService = newSiteService()

type siteService struct {
	common.BaseService[*model.Site]
	siteDomains map[string]*model.Site
	SiteIdMap   map[uint64]*model.Site
	lock        sync.Mutex
}

func newSiteService() *siteService {
	s := new(siteService)
	s.BaseService = common.BaseService[*model.Site]{
		Service: s,
	}
	s.siteDomains = make(map[string]*model.Site)
	s.SiteIdMap = make(map[uint64]*model.Site)
	return s
}

func (*siteService) Model() *model.Site {
	return new(model.Site)
}

func (s *siteService) Contains(host string) bool {
	_, ok := s.siteDomains[host]
	return ok
}

func (s *siteService) GetByDomain(host string) (*model.Site, error) {
	if site, ok := s.siteDomains[host]; ok {
		return site, nil
	}
	return s.getFromDB(&model.Site{
		Domain: host,
	})
}

func (s *siteService) getFromDB(condition *model.Site) (*model.Site, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if condition.ID != 0 {
		if site, ok := s.SiteIdMap[condition.ID]; ok {
			return site, nil
		}
	}
	if condition.Domain != "" {
		if site, ok := s.siteDomains[condition.Domain]; ok {
			return site, nil
		}
	}
	site := new(model.Site)
	if err := database.DB.First(site, condition).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	s.siteDomains[site.Domain] = site
	s.SiteIdMap[site.ID] = site
	return site, nil
}

func (s *siteService) GetByID(id uint64) (*model.Site, error) {
	if site, ok := s.SiteIdMap[id]; ok {
		return site, nil
	}
	return s.getFromDB(&model.Site{
		BaseModel: common.BaseModel{
			ID: id,
		},
	})
}
