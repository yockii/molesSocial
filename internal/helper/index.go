package helper

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/molesSocial/internal/constant"
	"strconv"
)

func GetCurrentManagerID(ctx *fiber.Ctx) (uint64, error) {
	// 获取当前登录的用户
	uidStr, _ := ctx.Locals(constant.JwtClaimManagerId).(string)
	if uidStr == "" {
		return 0, nil
	}
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		logger.Errorln(err)
		return 0, err
	}
	return uid, nil
}
