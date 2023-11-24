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

func init() {
	c := new(oauthController)
	Controllers = append(Controllers, c)
}

func (c *oauthController) InitRoute() {
	r := server.Group("/oauth")
	r.Get("/authorize", c.GetAuthorize)   // 授权页面
	r.Post("/authorize", c.PostAuthorize) // 授权请求

	r.Post("/token", c.Token) // 获得访问令牌
	r.Post("/revoke", c.RevokeToken)
}

func (c *oauthController) GetAuthorize(ctx *fiber.Ctx) error {
	ar := new(domain.OauthAuthorizeRequest)
	if err := ctx.QueryParser(ar); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if ar.ResponseType == "" || ar.ClientID == "" || ar.RedirectURI == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&domain.Response{
			Error:            "invalid_request",
			ErrorDescription: "response_type, client_id, redirect_uri are required",
		})
	}
	if ar.ResponseType != "code" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&domain.Response{
			Error:            "unsupported_response_type",
			ErrorDescription: "response_type must be code",
		})
	}

	// 获取对应app
	app, err := service.OauthApplicationService.Instance(&modelO.OauthApplication{ClientID: ar.ClientID})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if app == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("app not found")
	}

	// 判断redirect_uri是否一致
	if ar.RedirectURI != app.RedirectUri {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&domain.Response{
			Error:            "invalid_request",
			ErrorDescription: "redirect_uri is not match",
		})
	}

	// 获取对应站点信息
	host := ctx.Hostname()
	site, err := service.SiteService.GetByDomain(host)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if site == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("site not found")
	}

	// 判断是否已经登录
	token := ctx.Cookies(constant.CookiesNameToken)
	if !ar.ForceLogin && token != "" {
		// 验证token是否有效
		cachedUser, err := c.verifyToken(token)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		if cachedUser != nil {
			// 获取对应账号
			var account *modelA.Account
			account, err = service.AccountService.Instance(&modelA.Account{BaseModel: common.BaseModel{ID: cachedUser.AccountID}})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			if account == nil {
				return ctx.Status(fiber.StatusNotFound).SendString("account not found")
			}

			// 进入授权页面
			return ctx.Render("authorize", fiber.Map{
				"account": account,
				"site":    site,
				"app":     app,
				"state":   ar.State,
				"scopes":  strings.Split(ar.Scope, " "),
			}, "layouts/main")
		}
	}

	// 未登录，跳转到登录页面
	// 信息放redis中，并设置过期时间，生成唯一的code作为cookies传出，然后重定向到登录页面
	code, err := c.cacheOauthAuthorizeRequest(ar)
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
// request： client_id=nVPq20OCKrW8vw297NWMl22iz5qKiD9J4yvWMw_EHSk
//
//	         &client_secret=uALhB7JT7Ocfj3YA2lsr1ICZGYC67B9jqJbmVhg-arU
//	         &redirect_uri=https%3A%2F%2Fpinafore.social%2Fsettings%2Finstances%2Fadd - 可以是url作为重定向，如果是urn:ietf:wg:oauth:2.0:oob，则显示token；这里必须与注册app时保持一致
//	         &grant_type=authorization_code - 必须是authorization_code来根据code获取token，或者也可以是client_credentials获取程序级别的token
//	         &code=ZQ7ByYqb8J1rC6dCOrCH - 仅当grant_type为authorization_code时有效
//
//		response: {
//		   "access_token": "JxzJThdnrriruOjSU2HjCrBWmJzSaHEAH3HTL7rVJgA",
//		   "token_type": "Bearer",
//		   "scope": "read write follow push",
//		   "created_at": 1700532528
//		}
//
// 400 client error
// 401 unauthorized
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
	_, err = conn.Do("SETEX", constant.RedisPrefixAccessTokenAccount+accessToken, 60*60*24*7, account.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	// 缓存accessToken对应的appID
	_, err = conn.Do("SETEX", constant.RedisPrefixAccessTokenApp+accessToken, 60*60*24*7, app.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// 删除cookies中的CookiesNameOauthAuthorizeCode及对应缓存
	code := ctx.Cookies(constant.CookiesNameOauthAuthorizeCode)
	if code != "" {
		ctx.Cookie(&fiber.Cookie{
			Name:     constant.CookiesNameOauthAuthorizeCode,
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Minute),
			HTTPOnly: true,
			Secure:   true,
		})
		_, _ = conn.Do("DEL", constant.RedisPrefixOauthAuthorizeInfo+code)
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
		return ctx.Redirect(app.RedirectUri + "?error=access_denied&error_description=The+user+has+cancelled+entering+self-asserted+information&state=" + state)
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

// RevokeToken - 撤销访问令牌 client_id, client_secret, token | 200 ok, 403 forbidden
func (c *oauthController) RevokeToken(ctx *fiber.Ctx) error {
	ortr := new(domain.OauthRevokeTokenRequest)
	if err := ctx.BodyParser(ortr); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if ortr.ClientID == "" || ortr.ClientSecret == "" || ortr.Token == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request")
	}

	// 检查clientID和clientSecret是否匹配
	app, err := service.OauthApplicationService.Instance(&modelO.OauthApplication{ClientID: ortr.ClientID})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if app == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("app not found")
	}
	if app.ClientSecret != ortr.ClientSecret {
		return ctx.Status(fiber.StatusForbidden).SendString("invalid client_secret")
	}

	// 检查token是否有效
	conn := cache.Get()
	defer conn.Close()
	_, err = redis.Uint64(conn.Do("GET", constant.RedisPrefixAccessTokenAccount+ortr.Token))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	_, err = redis.Uint64(conn.Do("GET", constant.RedisPrefixAccessTokenApp+ortr.Token))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// 删除token
	_, _ = conn.Do("DEL", constant.RedisPrefixAccessTokenAccount+ortr.Token)
	_, _ = conn.Do("DEL", constant.RedisPrefixAccessTokenApp+ortr.Token)

	return ctx.JSON(&fiber.Map{})
}
