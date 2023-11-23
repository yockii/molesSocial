package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/domain"
	"github.com/yockii/molesSocial/internal/model"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/server"
	"net/url"
	"strings"
)

const (
	WebFingerHostMeta = `<?xml version="1.0" encoding="UTF-8"?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
  <Link rel="lrdd" template="%s://%s/.well-known/webfinger?resource={uri}" type="application/xrd+xml" />
</XRD>`
)

type WebFingerController struct{}

func (c *WebFingerController) WebFinger(ctx *fiber.Ctx) error {
	resource := ctx.Query("resource")
	if resource == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("resource is required")
	}
	u, err := url.ParseRequestURI(resource)
	if err != nil {
		u, err = url.ParseRequestURI("acct:" + resource)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString("resource is invalid")
		}
	}

	var host, username string
	if u.Scheme == "http" || u.Scheme == "https" {
		host = u.Host
		p := strings.TrimLeft(u.Path, "/")
		segs := strings.Split(p, "/")
		if segs[0] == "users" {
			username = segs[1]
		} else {
			username = segs[0]
		}
	} else if u.Scheme == "acct" {
		p := strings.TrimPrefix(u.Opaque, "@")
		segs := strings.Split(p, "@")
		if len(segs) != 2 {
			return ctx.Status(fiber.StatusBadRequest).SendString("resource is invalid")
		}
		username = segs[0]
		host = segs[1]
	}
	username = strings.TrimPrefix(username, "@")
	if username == "" || host == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("resource is invalid")
	}
	if host != ctx.Hostname() && !service.SiteService.Contains(host) {
		return ctx.Status(fiber.StatusBadRequest).SendString("resource is invalid")
	}

	var site *model.Site
	site, err = service.SiteService.GetByDomain(host)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("resource is invalid")
	}
	if site == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("site is invalid")
	}

	// 根据username和host获取用户信息
	account, err := service.AccountService.GetByUsernameAndSite(username, site.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("resource is invalid")
	}

	return ctx.JSON(&domain.WellKnownResponse{
		Subject: constant.WebFingerAccountPrefix + ":" + account.Username + "@" + site.Domain,
		Aliases: []string{
			account.Uri,
			account.Url,
		},
		Links: []domain.WellKnownLink{
			{
				Rel:  constant.WebFingerProfilePage,
				Type: constant.WebFingerProfilePageContentType,
				Href: account.Url,
			},
			{
				Rel:  constant.WebFingerSelf,
				Type: constant.WebFingerSelfContentType,
				Href: account.Uri,
			},
		},
	})
}

func (c *WebFingerController) HostMeta(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, "application/xrd+xml")
	return ctx.SendString(fmt.Sprintf(WebFingerHostMeta, ctx.Protocol(), ctx.Hostname()))
}

func (c *WebFingerController) NodeInfo(ctx *fiber.Ctx) error {
	host := ctx.Hostname()
	protocol := ctx.Protocol()
	return ctx.JSON(&domain.WellKnownResponse{
		Links: []domain.WellKnownLink{
			{
				Rel:  constant.NodeInfoRel,
				Href: fmt.Sprintf("%s://%s/nodeinfo/%s", protocol, host, constant.NodeInfoVersion),
			},
		},
	})
}

func init() {
	c := new(WebFingerController)

	Controllers = append(Controllers, c)
}

func (c *WebFingerController) InitRoute() {
	wk := server.Group("/.well-known")

	wk.Get("/webfinger", c.WebFinger)
	wk.Get("/nodeinfo", c.NodeInfo)
	wk.Get("/host-meta", c.HostMeta)
}
