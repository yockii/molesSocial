package service

import (
	model "github.com/yockii/molesSocial/internal/model/manage"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

type roleService struct {
	common.BaseService[*model.Role]
}

var RoleService = newRoleService()

func newRoleService() *roleService {
	s := new(roleService)
	s.BaseService = common.BaseService[*model.Role]{
		Service: s,
	}
	return s
}

func (*roleService) Model() *model.Role {
	return new(model.Role)
}

func (s *roleService) ResourceCodes(id uint64) ([]string, error) {
	var codes []string
	err := database.DB.Model(new(model.Resource)).Where("id in (?)", database.DB.Model(new(model.RoleResource)).Where("role_id = ?", id).Select("resource_id")).Pluck("code", &codes).Error
	return codes, err
}
