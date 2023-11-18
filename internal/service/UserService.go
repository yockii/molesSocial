package service

import (
	"errors"
	modelU "github.com/yockii/molesSocial/internal/model/user"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
	"gorm.io/gorm"
)

var UserService = newUserService()

type userService struct {
	common.BaseService[*modelU.User]
}

func newUserService() *userService {
	s := new(userService)
	s.BaseService = common.BaseService[*modelU.User]{
		Service: s,
	}
	return s
}

func (*userService) Model() *modelU.User {
	return new(modelU.User)
}

func (s *userService) GetByAccountID(id uint64) (*modelU.User, error) {
	if id == 0 {
		return nil, nil
	}
	user := s.Model()
	err := database.DB.Where("account_id = ?", id).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return user, nil
}

func (s *userService) GetByEmailAndDomain(email string, domain string) (*modelU.User, error) {
	if email == "" || domain == "" {
		return nil, nil
	}
	user := s.Model()
	err := database.DB.Where("email = ? AND domain = ?", email, domain).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return user, nil
}
