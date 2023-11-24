package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/domain"
	"github.com/yockii/molesSocial/internal/middleware"
	model "github.com/yockii/molesSocial/internal/model/oauth"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/server"
)

type ApplicationController struct{}

func (c *ApplicationController) ApplicationCreate(ctx *fiber.Ctx) error {
	d := new(domain.ApplicationCreateRequest)
	if err := ctx.BodyParser(d); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	host := ctx.Hostname()
	// 判断必填项
	if d.ClientName == "" || d.RedirectURIs == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&domain.Response{
			Error: "client_name, redirect_uris are required",
		})
	}

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
	r.Get("/verify_credentials", middleware.AppAccessTokenMiddleware(false), c.VerifyApplication)
}

func (c *ApplicationController) VerifyApplication(ctx *fiber.Ctx) error {
	appID, _ := ctx.Locals(constant.CtxAppID).(uint64)
	if appID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Authorization Token")
	}
	app, err := service.OauthApplicationService.Instance(&model.OauthApplication{BaseModel: common.BaseModel{ID: appID}})
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Authorization Token")
	}
	app.ClientID = ""
	app.ClientSecret = ""
	return ctx.JSON(app)
}
