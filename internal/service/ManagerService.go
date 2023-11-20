package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/molesSocial/internal/constant"
	model "github.com/yockii/molesSocial/internal/model/manage"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type managerService struct {
	common.BaseService[*model.Manager]
}

var ManagerService = newManagerService()

func newManagerService() *managerService {
	s := new(managerService)
	s.BaseService = common.BaseService[*model.Manager]{
		Service: s,
	}
	return s
}

func (*managerService) Model() *model.Manager {
	return new(model.Manager)
}

// Add 添加用户
func (s *managerService) Add(instance *model.Manager, tx ...*gorm.DB) (duplicated bool, err error) {
	if instance.Username == "" {
		err = errors.New("username is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.Manager{}).Where(&model.Manager{Username: instance.Username}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	instance.ID = util.SnowflakeId()
	if instance.Password != "" {
		pwd, _ := bcrypt.GenerateFromPassword([]byte(instance.Password), bcrypt.DefaultCost)
		instance.Password = string(pwd)
	}
	instance.Status = model.ManagerStatusNormal

	// 获取默认角色
	defaultRole := &model.Role{DefaultRole: &constant.BoolTrue}
	if err = database.DB.Where(defaultRole).First(defaultRole).Error; err != nil {
		logger.Errorln(err)
		return
	}
	if defaultRole != nil && defaultRole.ID > 0 {
		if len(tx) == 0 {
			// 添加用户的同时要添加默认角色
			err = database.DB.Transaction(func(tx *gorm.DB) error {
				if err = tx.Create(instance).Error; err != nil {
					logger.Errorln(err)
					return err
				}
				managerRole := &model.ManagerRole{
					BaseModel: common.BaseModel{ID: util.SnowflakeId()},
					ManagerID: instance.ID,
					RoleID:    defaultRole.ID,
				}
				if err = tx.Create(managerRole).Error; err != nil {
					logger.Errorln(err)
					return err
				}
				return nil
			})
		} else {
			if err = tx[0].Create(instance).Error; err != nil {
				logger.Errorln(err)
				return
			}
			managerRole := &model.ManagerRole{
				BaseModel: common.BaseModel{ID: util.SnowflakeId()},
				ManagerID: instance.ID,
				RoleID:    defaultRole.ID,
			}
			if err = tx[0].Create(managerRole).Error; err != nil {
				logger.Errorln(err)
				return
			}
		}
	} else {
		if len(tx) == 0 {
			err = database.DB.Transaction(func(tx *gorm.DB) error {
				if err = tx.Create(instance).Error; err != nil {
					logger.Errorln(err)
					return err
				}
				return nil
			})
		} else {
			err = tx[0].Create(instance).Error
		}
	}
	if err != nil {
		logger.Errorln(err)
		return
	}
	// 完成后密码置空
	instance.Password = ""
	return
}

// ManagerRoles 获取用户角色
func (s *managerService) ManagerRoles(userId uint64) ([]*model.Role, error) {
	var roles []*model.Role
	// 关联查询
	err := database.DB.Model(new(model.Role)).Where("id in (?)", database.DB.Model(new(model.ManagerRole)).Where("manager_id = ?", userId).Select("role_id")).Find(&roles).Error
	return roles, err
}

// LoginWithUsernameAndPassword 用户登录
func (s *managerService) LoginWithUsernameAndPassword(username, password string) (instance *model.User, passwordNotMatch bool, err error) {
	instance = new(model.User)
	err = database.DB.Where(&model.User{Username: username}).First(instance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			passwordNotMatch = true
			return
		}
		logger.Errorln(err)
		return
	}
	if instance.Status != model.UserStatusNormal {
		err = errors.New("用户已被禁用")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(instance.Password), []byte(password))
	if err != nil {
		passwordNotMatch = true
		err = nil
		return
	}
	// 完成后密码置空
	instance.Password = ""
	return
}
