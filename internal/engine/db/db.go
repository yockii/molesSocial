package db

import (
	"fmt"
	core "github.com/gofiber/template"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/molesSocial/internal/model"
	modelT "github.com/yockii/molesSocial/internal/model/template"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/database"
	"html/template"
	"io"
)

type Engine struct {
	core.Engine
	Templates *template.Template
}

func New() *Engine {
	engine := &Engine{
		Engine: core.Engine{
			Left:       "{{",
			Right:      "}}",
			LayoutName: "embed",
			Funcmap:    make(map[string]interface{}),
		},
	}
	engine.AddFunc(engine.LayoutName, func() error {
		return fmt.Errorf("layoutName called unexpectedly")
	})
	return engine
}

// Load 加载所有模板数据
func (e *Engine) Load() error {
	if e.Loaded {
		return nil
	}
	e.Mutex.Lock()
	defer e.Mutex.Unlock()
	e.Templates = template.New("")

	e.Templates.Delims(e.Left, e.Right)
	e.Templates.Funcs(e.Funcmap)
	e.Loaded = true

	// 查询所有模板并进行加载
	var tl []*modelT.Template
	err := database.DB.Find(&tl).Error
	if err != nil {
		logger.Errorln(err)
		return err
	}

	for _, t := range tl {
		var site *model.Site
		site, err = service.SiteService.GetByID(t.SiteID)
		if err != nil {
			logger.Errorln(err)
			return err
		}
		if site == nil {
			continue
		}
		_, err = e.Templates.New(site.Domain + "/" + t.Name).Parse(t.Content)
		if err != nil {
			logger.Errorln(err)
			return err
		}
	}
	return nil
}

func (e *Engine) Render(out io.Writer, name string, binding interface{}, layout ...string) error {
	if !e.Loaded || e.ShouldReload {
		if e.ShouldReload {
			e.Loaded = false
		}
		if err := e.Load(); err != nil {
			return err
		}
	}
	tmpl := e.Templates.Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("template %s not found", name)
	}

	if len(layout) > 0 && layout[0] != "" {
		l := e.Templates.Lookup(layout[0])
		if l == nil {
			return fmt.Errorf("layout %s not found", layout[0])
		}
		l.Funcs(map[string]interface{}{
			e.LayoutName: func() error {
				return tmpl.Execute(out, binding)
			},
		})
		return l.Execute(out, binding)
	}
	return tmpl.Execute(out, binding)
}
