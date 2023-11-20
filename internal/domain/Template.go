package domain

import (
	model "github.com/yockii/molesSocial/internal/model/template"
	"github.com/yockii/qscore/pkg/common"
)

type Template struct {
	model.Template
	common.BaseDomain[*model.Template]
}

func (t *Template) GetModel() *model.Template {
	return &t.Template
}
