package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gomodule/redigo/redis"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/domain"
	"github.com/yockii/molesSocial/internal/helper"
	"github.com/yockii/molesSocial/internal/middleware"
	model "github.com/yockii/molesSocial/internal/model/manage"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/cache"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/config"
	"github.com/yockii/qscore/pkg/crypto"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/util"
	"strconv"
)

type managerController struct {
	common.BaseController[*model.Manager, *domain.Manager]
}

func (c *managerController) GetService() common.Service[*model.Manager] {
	return service.ManagerService
}

func (c *managerController) NewModel() *model.Manager {
	return new(model.Manager)
}

func (c *managerController) NewDomain() *domain.Manager {
	return new(domain.Manager)
}

func (c *managerController) InitRoute() {
	r := server.Group("/api/v1/manager", middleware.NeedAuthorization("manager"))

	r.Post("/add", c.Add)
	r.Post("/update", c.Update)
	r.Put("/update", c.Update)
	r.Post("/delete", c.Delete)
	r.Delete("/delete", c.Delete)
	r.Get("/list", c.List)
	r.Get("/detail", c.Detail)
}

func (c *managerController) Add(ctx *fiber.Ctx) error {
	instance := new(model.Manager)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Username == "" || instance.Password == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " managername, password",
		})
	}

	// 解析密码
	if pwd, err := crypto.Sm2Decrypt(instance.Password); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + " password decrypt failed",
		})
	} else {
		instance.Password = pwd
	}

	if instance.Password != "" {
		isStrongEnough := util.PasswordStrengthCheck(8, 50, 3, instance.Password)
		if !isStrongEnough {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodePasswordStrengthInvalid,
				Msg:  server.ResponseMsgPasswordStrengthInvalid + ", password is not strong enough",
			})
		}
	}

	duplicated, err := service.ManagerService.Add(instance)
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if duplicated {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		})
	}
	instance.Password = ""
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *managerController) UpdatePassword(ctx *fiber.Ctx) (err error) {
	instance := new(domain.Manager)
	if err = ctx.BodyParser(instance); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if instance.OldPassword == "" || instance.Password == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": 密码(原/新)",
		})
	}

	instance.ID, err = helper.GetCurrentManagerID(ctx)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  "获取当前用户失败",
		})
	}

	// 解析密码
	var pwd string
	if pwd, err = crypto.Sm2Decrypt(instance.OldPassword); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + "密码加密数据不准确",
		})
	}
	var newPwd string
	if newPwd, err = crypto.Sm2Decrypt(instance.Password); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + "密码加密数据不准确",
		})
	}
	if newPwd != "" {
		isStrong := util.PasswordStrengthCheck(8, 50, 3, instance.Password)
		if !isStrong {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodePasswordStrengthInvalid,
				Msg:  server.ResponseMsgPasswordStrengthInvalid,
			})
		}
	}
	instance.Password = newPwd

	// 检查旧密码是否匹配
	var pass bool
	pass, err = service.ManagerService.CheckManagerPassword(instance.ID, pwd)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	} else if !pass {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  "原密码" + server.ResponseMsgDataNotMatch,
		})
	}

	var success bool
	success, err = service.ManagerService.UpdatePassword(&instance.Manager)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *managerController) ResetPassword(ctx *fiber.Ctx) error {
	instance := new(model.Manager)
	if err := ctx.BodyParser(instance); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if instance.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}

	// 生成密码
	result := util.RandomPassword(8, 20, 3)
	instance.Password = result
	success, err := service.ManagerService.UpdatePassword(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if !success {
		result = ""
	}
	return ctx.JSON(&server.CommonResponse{
		Data: result,
	})
}

func (c *managerController) Login(ctx *fiber.Ctx) error {
	instance := new(model.Manager)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Username == "" || instance.Password == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": 用户名及密码",
		})
	}

	// 解析密码
	if pwd, err := crypto.Sm2Decrypt(instance.Password); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + "密码不准确",
		})
	} else {
		instance.Password = pwd
	}

	isStrong := util.PasswordStrengthCheck(8, 50, 4, instance.Password)
	if !isStrong {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodePasswordStrengthInvalid,
			Msg:  server.ResponseMsgPasswordStrengthInvalid,
		})
	}

	manager, notMatch, err := service.ManagerService.LoginWithManagernameAndPassword(instance.Username, instance.Password)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if notMatch {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  "用户名与密码" + server.ResponseMsgDataNotMatch,
		})
	}
	return c.generateLoginResponse(manager, ctx)
}

func (c *managerController) generateLoginResponse(manager *model.Manager, ctx *fiber.Ctx) error {
	jwtToken, err := generateJwtToken(strconv.FormatUint(manager.ID, 10), "")
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeGeneration,
			Msg:  server.ResponseMsgGeneration + err.Error(),
		})
	}
	manager.Password = ""
	managerDomain := &domain.Manager{
		Manager: *manager,
	}

	return ctx.JSON(&server.CommonResponse{
		Data: map[string]interface{}{
			"token":   jwtToken,
			"manager": managerDomain,
		},
	})
}

func (c *managerController) UpdateMyInfo(ctx *fiber.Ctx) error {
	instance := new(domain.Manager)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	managerId, err := helper.GetCurrentManagerID(ctx)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  "未知错误",
		})
	}
	instance.ID = managerId
	success, err := service.ManagerService.UpdateDomain(&domain.Manager{
		Manager: model.Manager{
			BaseModel: common.BaseModel{ID: instance.ID},
			Email:     instance.Email,
		},
	})
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func generateJwtToken(managerId, tenantId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	sid := util.GenerateXid()

	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	sessionKey := constant.RedisPrefixSessionId + sid

	_, err := conn.Do("SETEX", sessionKey, config.GetInt("managerTokenExpire"), managerId)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	claims[constant.JwtClaimManagerId] = managerId
	claims[constant.JwtClaimTenantId] = tenantId
	claims[constant.JwtClaimSessionId] = sid

	t, err := token.SignedString([]byte(constant.JwtSecret))
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	return t, nil
}

func init() {
	c := new(managerController)
	c.BaseController = common.BaseController[*model.Manager, *domain.Manager]{
		Controller: c,
	}
	Controllers = append(Controllers, c)
}
