package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/sjson"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/domain"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/server"
	"strings"
)

type nodeController struct{}

func init() {
	c := new(nodeController)
	Controllers = append(Controllers, c)
}

func (c *nodeController) InitRoute() {
	server.Get("/nodeinfo/:version", c.VersionedNodeInfo)
	server.Get("/api/v1/instance", c.NodeInstance)
}

func (c *nodeController) VersionedNodeInfo(ctx *fiber.Ctx) error {
	version := ctx.Params("version")
	if version != constant.NodeInfoVersion {
		return ctx.Status(fiber.StatusBadRequest).SendString("version is invalid")
	}

	ctx.Set(fiber.HeaderContentType, "application/json; profile=\"http://nodeinfo.diaspora.software/ns/schema/"+version+"#\"")

	total, monthly, halfYears, err := service.AccountService.CountUsers(ctx.Hostname())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("user data process error")
	}
	// TODO 查询本地post数量

	var totalPosts int64
	totalPosts, err = service.StatusService.CountLocalStatuses(ctx.Hostname())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("status data process error")
	}

	return ctx.JSON(&domain.NodeInfo{
		Version: version,
		Software: &domain.NodeInfoSoftware{
			Name:    constant.NodeInfoSoftwareName,
			Version: constant.Version,
		},
		Protocols: constant.NodeInfoProtocols,
		Services: &domain.NodeInfoServices{
			Inbound:  constant.NodeInfoInbound,
			Outbound: constant.NodeInfoOutbound,
		},
		OpenRegistrations: true,
		Usage: &domain.NodeInfoUsage{
			Users: &domain.NodeInfoUsers{
				Total:          total,
				ActiveMonth:    monthly,
				ActiveHalfyear: halfYears,
			},
			LocalPosts: totalPosts,
		},
		Metadata: constant.NodeInfoMetadata,
	})
}

func (c *nodeController) NodeInstance(ctx *fiber.Ctx) error {
	host := ctx.Hostname()
	site, err := service.SiteService.GetByDomain(host)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("site data process error")
	}
	if site == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("site node not found")
	}
	instance := &domain.InstanceInfo{
		URI:         host,
		Title:       site.Name,
		Description: site.Description,
		Email:       site.Email,
		Version:     constant.Version,
		Urls: map[string]string{
			"streaming_api": "wss://" + host + "/api/v1/streaming",
		},
		Stats: map[string]int64{
			"domain_count": 0,
			"user_count":   0,
			"status_count": 0,
		},
		ShortDescription: site.ShortDesc,
		Thumbnail:        site.Thumbnail,
		Languages:        strings.Split(site.Languages, ","),
		ContactAccount:   &domain.AccountInfo{},
	}

	result, err := json.Marshal(instance)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("node data process error")
	}

	result, _ = sjson.SetRawBytes(result, "configuration", []byte(site.Configuration))

	// configuration
	// rules

	return ctx.Send(result)
}
