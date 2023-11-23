package controller

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/domain"
	modelA "github.com/yockii/molesSocial/internal/model/account"
	modelO "github.com/yockii/molesSocial/internal/model/oauth"
	modelU "github.com/yockii/molesSocial/internal/model/user"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/cache"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/util"
	"strconv"
	"strings"
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
func (c *oauthController) Authorize(ctx *fiber.Ctx) (err error) {
	// 从cookies中取出molesocial_oauth_authorize_code信息
	var authorizeRequest *domain.OauthAuthorizeRequest
	var arCached bool
	code := ctx.Cookies(constant.CookiesNameOauthAuthorizeCode)
	token := ctx.Cookies(constant.CookiesNameToken)
	if code != "" {
		// 有code，从redis中取出并删除OauthAuthorizeRequest信息
		authorizeRequest, err = c.popOauthAuthorizeRequest(code)
		// cookies中删除
		ctx.Cookie(&fiber.Cookie{
			Name:     constant.CookiesNameOauthAuthorizeCode,
			Value:    "",
			Path:     "/",
			Expires:  time.Now(),
			HTTPOnly: true,
			Secure:   true,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	if authorizeRequest == nil {
		authorizeRequest = new(domain.OauthAuthorizeRequest)
		if err = ctx.QueryParser(authorizeRequest); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if authorizeRequest.ResponseType == "" || authorizeRequest.ClientID == "" || authorizeRequest.RedirectURI == "" {
			return ctx.Status(fiber.StatusBadRequest).SendString("invalid request")
		}
		// 判断response_type是否为code
		if authorizeRequest.ResponseType != "code" {
			return ctx.Status(fiber.StatusBadRequest).SendString("invalid response_type")
		}

		if token == "" {
			// 则为新进入，需要存储信息并跳转
			// 信息放redis中，并设置过期时间，生成唯一的code作为cookies传出，然后重定向到登录页面
			code, err = c.cacheAuthorizeRequest(ctx, authorizeRequest)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			arCached = true
			return ctx.Redirect("/auth/sign_in")
		}
	}

	// 验证token是否有效
	cachedUser, err := c.verifyToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if cachedUser == nil {
		if !arCached {
			code, err = c.cacheAuthorizeRequest(ctx, authorizeRequest)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}
		return ctx.Redirect("/auth/sign_in")
	}

	// 获取对应账号
	var account *modelA.Account
	account, err = service.AccountService.Instance(&modelA.Account{BaseModel: common.BaseModel{ID: cachedUser.AccountID}})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if account == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("account not found")
	}

	// 获取对应站点信息
	host := ctx.Hostname()
	site, err := service.SiteService.GetByDomain(host)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// 获取app
	var app *modelO.OauthApplication
	app, err = service.OauthApplicationService.Instance(&modelO.OauthApplication{ClientID: authorizeRequest.ClientID})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// 展示授权页面
	return ctx.Render("authorize", fiber.Map{
		"account": account,
		"site":    site,
		"app":     app,
		"state":   authorizeRequest.State,
		"scopes":  strings.Split(authorizeRequest.Scope, " "),
	}, "layouts/main")
}

func (c *oauthController) cacheAuthorizeRequest(ctx *fiber.Ctx, authorizeRequest *domain.OauthAuthorizeRequest) (code string, err error) {
	code, err = c.cacheOauthAuthorizeRequest(authorizeRequest)
	if err != nil {
		return
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     constant.CookiesNameOauthAuthorizeCode,
		Value:    code,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
	})
	return
}

func (c *oauthController) popOauthAuthorizeRequest(code string) (*domain.OauthAuthorizeRequest, error) {
	conn := cache.Get()
	defer conn.Close()
	oarJson, err := redis.Bytes(conn.Do("GET", constant.RedisPrefixOauthAuthorizeInfo+code))
	if err != nil {
		return nil, err
	}
	req := new(domain.OauthAuthorizeRequest)
	if err = json.Unmarshal(oarJson, req); err != nil {
		return nil, err
	}
	_, err = conn.Do("DEL", constant.RedisPrefixOauthAuthorizeInfo+code)
	return req, nil
}

func (c *oauthController) cacheOauthAuthorizeRequest(req *domain.OauthAuthorizeRequest) (string, error) {
	conn := cache.Get()
	defer conn.Close()
	code := util.GenerateXid()
	reqJson, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	_, err = conn.Do("SETEX", constant.RedisPrefixOauthAuthorizeInfo+code, 60*5, reqJson)
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
	r.Post("/authorize", c.PostAuthorize)

	r.Post("/token", c.Token)
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

// Token  http://localhost:9980/oauth/token
// request： client_id=nVPq20OCKrW8vw297NWMl22iz5qKiD9J4yvWMw_EHSk&client_secret=uALhB7JT7Ocfj3YA2lsr1ICZGYC67B9jqJbmVhg-arU&redirect_uri=https%3A%2F%2Fpinafore.social%2Fsettings%2Finstances%2Fadd&grant_type=authorization_code&code=ZQ7ByYqb8J1rC6dCOrCH
//
//	response: {
//	   "access_token": "JxzJThdnrriruOjSU2HjCrBWmJzSaHEAH3HTL7rVJgA",
//	   "token_type": "Bearer",
//	   "scope": "read write follow push",
//	   "created_at": 1700532528
//	}
func (c *oauthController) Token(ctx *fiber.Ctx) error {
	req := new(domain.OauthTokenRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if req.GrantType != "authorization_code" {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid grant_type")
	}
	if req.ClientID == "" || req.ClientSecret == "" || req.Code == "" || req.RedirectURI == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request")
	}

	// 获取对应app
	app, err := service.OauthApplicationService.Instance(&modelO.OauthApplication{ClientID: req.ClientID})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if app == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("app not found")
	}

	// 验证code是否有效
	conn := cache.Get()
	defer conn.Close()
	userId, err := redis.Uint64(conn.Do("GET", constant.RedisPrefixAccessCode+req.Code))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("invalid code")
	}
	if userId == 0 {
		return ctx.Status(fiber.StatusInternalServerError).SendString("invalid code")
	}
	// 删除code
	_, err = conn.Do("DEL", constant.RedisPrefixAccessCode+req.Code)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// 获取对应账号
	var account *modelA.Account
	account, err = service.AccountService.Instance(&modelA.Account{BaseModel: common.BaseModel{ID: userId}})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if account == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("account not found")
	}

	// 生成accessToken，缓存账号ID
	accessToken := util.GenerateXid()
	_, err = conn.Do("SETEX", constant.RedisPrefixUserToken+accessToken, 60*60*24*7, account.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.JSON(&domain.OauthTokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		Scope:       app.Scopes,
		CreatedAt:   time.Now().Unix(),
	})
}

func (c *oauthController) PostAuthorize(ctx *fiber.Ctx) error {
	accepted := ctx.FormValue("authorize") == "1"
	appIdStr := ctx.FormValue("appId")
	state := ctx.FormValue("state")
	token := ctx.Cookies(constant.CookiesNameToken)
	if token == "" || appIdStr == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request")
	}

	appId, err := strconv.ParseUint(appIdStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request")
	}
	app, err := service.OauthApplicationService.Instance(&modelO.OauthApplication{BaseModel: common.BaseModel{ID: appId}})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if app == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("app not found")
	}

	if !accepted {
		return ctx.Redirect(app.RedirectUri + "?error=access_denied&state=" + state)
	}

	// 验证token是否有效
	cachedUser, err := c.verifyToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if cachedUser == nil {
		return ctx.Redirect("/auth/sign_in")
	}

	// 生成code并跳转
	code, err := c.generateCode(cachedUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return ctx.Redirect(app.RedirectUri + "?code=" + code + "&state=" + state)
}

func (c *oauthController) generateCode(user *modelU.User) (string, error) {
	conn := cache.Get()
	defer conn.Close()
	code := util.GenerateXid()
	_, err := conn.Do("SETEX", constant.RedisPrefixAccessCode+code, 60*5, user.AccountID)
	if err != nil {
		return "", err
	}
	return code, nil
}
