package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/molesSocial/internal/middleware"
	"github.com/yockii/qscore/pkg/server"
)

type emailController struct{}

func init() {
	c := new(emailController)
	Controllers = append(Controllers, c)
}

func (c *emailController) InitRoute() {
	r := server.Group("/api/v1//emails")
	r.Post("/confirmations", middleware.AppAccessTokenMiddleware(false), c.sendConfirmationEmail)
}

func (c *emailController) sendConfirmationEmail(ctx *fiber.Ctx) error {
	// TODO ?email=xxx 如果有则更新，否则使用当前登录用户的邮箱发送确认邮件
	return nil
}
