package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/qscore/pkg/server"
)

type RouterController interface {
	// InitRoute 初始化路由
	InitRoute()
}

var Controllers []RouterController

func InitRouter() {
	InitRoute()

	server.All("/*", func(ctx *fiber.Ctx) error {
		// 记录请求
		logger.Traceln(ctx.Method(), ctx.Path())
		return ctx.SendStatus(fiber.StatusNotFound)
	})
}

func InitRoute() {
	for _, c := range Controllers {
		c.InitRoute()
	}
}
