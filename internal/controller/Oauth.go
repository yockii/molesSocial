package controller

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/domain"
	modelU "github.com/yockii/molesSocial/internal/model/user"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/cache"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/util"
	"time"
)

type oauthController struct{}

// Authorize  http://localhost:9980/oauth/authorize
// 有可能是其他客户端定位过来
// ?client_id=claujf0udbsi782n19d0
// &redirect_uri=https%3A%2F%2Fpinafore.social%2Fsettings%2Finstances%2Fadd
// &response_type=code
// &scope=read%20write%20follow%20push
// 有可能是登录后定位过来
func (c *oauthController) Authorize(ctx *fiber.Ctx) error {
	// 从cookies中获取token
	token := ctx.Cookies(constant.CookiesNameToken)
	if token == "" {
		req := new(domain.OauthAuthorizeRequest)
		if err := ctx.QueryParser(req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if req.ResponseType == "" || req.ClientID == "" || req.RedirectURI == "" {
			return ctx.Status(fiber.StatusBadRequest).SendString("invalid request")
		}

		// 判断response_type是否为code
		if req.ResponseType != "code" {
			return ctx.Status(fiber.StatusBadRequest).SendString("invalid response_type")
		}

		// 信息放redis中，并设置过期时间，生成唯一的code作为cookies传出，然后重定向到登录页面
		code, err := c.cacheOauthAuthorizeRequest(req)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		ctx.Cookie(&fiber.Cookie{
			Name:     constant.CookiesNameOauthAuthorizeCode,
			Value:    code,
			Path:     "/",
			Expires:  time.Now().Add(5 * time.Minute),
			HTTPOnly: true,
			Secure:   true,
		})
		return ctx.Redirect("/auth/sign_in")
	}
	// 验证token是否有效
	cachedUser, err := c.verifyToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if cachedUser == nil {
		return ctx.Redirect("/auth/sign_in")
	}

	// 验证完成，清楚缓存的OauthAuthorizeRequest
	// 从cookies中取出molesocial_oauth_authorize_code信息
	code := ctx.Cookies(constant.CookiesNameOauthAuthorizeCode)
	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request")
	}
	// 从redis中取出并删除OauthAuthorizeRequest信息
	authorizeRequest, err := c.popOauthAuthorizeRequest(code)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// 获取对应站点信息
	host := ctx.Hostname()
	site, err := service.SiteService.GetByDomain(host)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// 展示授权页面
	return ctx.Render(site.Domain+"/authorize", fiber.Map{
		"site":             site,
		"authorizeRequest": authorizeRequest,
	})
}

func (c *oauthController) popOauthAuthorizeRequest(code string) (*domain.OauthAuthorizeRequest, error) {
	conn := cache.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HGETALL", constant.RedisPrefixOauthAuthorizeInfo+code))
	if err != nil {
		return nil, err
	}
	req := new(domain.OauthAuthorizeRequest)
	if err = redis.ScanStruct(values, req); err != nil {
		return nil, err
	}
	_, err = conn.Do("DEL", constant.RedisPrefixOauthAuthorizeInfo+code)
	return req, nil
}

func (c *oauthController) cacheOauthAuthorizeRequest(req *domain.OauthAuthorizeRequest) (string, error) {
	conn := cache.Get()
	defer conn.Close()
	code := util.GenerateXid()
	_, err := conn.Do("HMSET", constant.RedisPrefixOauthAuthorizeInfo+code, "client_id", req.ClientID, "redirect_uri", req.RedirectURI, "scope", req.Scope, "state", req.State)
	if err != nil {
		return "", err
	}
	_, err = conn.Do("EXPIRE", code, 60*5) // 5分钟过期
	if err != nil {
		return "", err
	}
	return code, nil
}

func init() {
	c := new(oauthController)
	Controllers = append(Controllers, c)
}

func (c *oauthController) InitRoute() {
	r := server.Group("/oauth")
	r.Get("/authorize", c.Authorize)
}

func (c *oauthController) verifyToken(token string) (*modelU.User, error) {
	conn := cache.Get()
	defer conn.Close()
	userJson, err := redis.Bytes(conn.Do("GET", constant.RedisPrefixUserToken+token))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return nil, nil
		}
		return nil, err
	}
	user := new(modelU.User)
	if err = json.Unmarshal(userJson, user); err != nil {
		return nil, err
	}
	return user, nil
}
