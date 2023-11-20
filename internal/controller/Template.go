package controller

import (
	"github.com/yockii/molesSocial/internal/domain"
	"github.com/yockii/molesSocial/internal/middleware"
	model "github.com/yockii/molesSocial/internal/model/template"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/server"
)

type templateController struct {
	common.BaseController[*model.Template, *domain.Template]
}

func (c *templateController) GetService() common.Service[*model.Template] {
	return service.TemplateService

}
func (c *templateController) NewModel() *model.Template {
	return new(model.Template)
}
func (c *templateController) NewDomain() *domain.Template {
	return new(domain.Template)
}

func (c *templateController) InitRoute() {
	r := server.Group("/api/v1/template", middleware.NeedAuthorization("template"))

	r.Post("/add", c.Add)
	r.Post("/update", c.Update)
	r.Put("/update", c.Update)
	r.Post("/delete", c.Delete)
	r.Delete("/delete", c.Delete)
	r.Get("/list", c.List)
	r.Get("/detail", c.Detail)
}

func init() {
	c := new(templateController)
	c.BaseController = common.BaseController[*model.Template, *domain.Template]{
		Controller: c,
	}
	Controllers = append(Controllers, c)
}
