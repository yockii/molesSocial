package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/molesSocial/internal/domain"
	model "github.com/yockii/molesSocial/internal/model/oauth"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/server"
)

type ApplicationController struct{}

func (c *ApplicationController) ApplicationCreate(ctx *fiber.Ctx) error {
	d := new(domain.ApplicationCreateRequest)
	if err := ctx.BodyParser(d); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	host := ctx.Hostname()
	//TODO 判断必填项

	app, err := service.OauthApplicationService.Save(&model.OauthApplication{
		Name:        d.ClientName,
		Domain:      host,
		RedirectUri: d.RedirectURIs,
		Scopes:      d.Scopes,
		Website:     d.Website,
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return ctx.JSON(app)
}

func init() {
	c := new(ApplicationController)
	Controllers = append(Controllers, c)
}

func (c *ApplicationController) InitRoute() {
	r := server.Group("/api/v1/apps")
	r.Post("/", c.ApplicationCreate)
}
