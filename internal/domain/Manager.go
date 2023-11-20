package domain

import (
	model "github.com/yockii/molesSocial/internal/model/manage"
	"github.com/yockii/qscore/pkg/common"
)

type Manager struct {
	model.Manager
	common.BaseDomain[*model.Manager]
	OldPassword string `json:"oldPassword,omitempty"`
}

func (m *Manager) GetModel() *model.Manager {
	return &m.Manager
}
