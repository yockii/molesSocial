package service

import (
	"github.com/yockii/molesSocial/internal/model/template"
	"github.com/yockii/qscore/pkg/common"
)

var TemplateService = newTemplateService()

type templateService struct {
	common.BaseService[*model.Template]
}

func newTemplateService() *templateService {
	s := new(templateService)
	s.BaseService = common.BaseService[*model.Template]{
		Service: s,
	}
	return s
}

func (*templateService) Model() *model.Template {
	return new(model.Template)
}
