package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/molesSocial/internal/middleware"
	"github.com/yockii/qscore/pkg/server"
)

// accountController @see https://docs.joinmastodon.org/methods/accounts/
type accountController struct{}

func init() {
	c := new(accountController)
	Controllers = append(Controllers, c)
}

func (c *accountController) InitRoute() {
	r := server.Group("/api/v1/accounts")
	r.Post("/", middleware.AppAccessTokenMiddleware(true), c.RegeditAccount)

}

func (c *accountController) RegeditAccount(ctx *fiber.Ctx) error {
	// TODO 注册一个账号
	// form: username\email\password\agreement\locale\reason
	// return {access_token, token_type, scope, created_at}

	return nil
}
