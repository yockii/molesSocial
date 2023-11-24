package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/domain"
	"github.com/yockii/molesSocial/internal/model"
	modelA "github.com/yockii/molesSocial/internal/model/account"
	modelO "github.com/yockii/molesSocial/internal/model/oauth"
	modelU "github.com/yockii/molesSocial/internal/model/user"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/cache"
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type authController struct{}

func (c *authController) SignIn(ctx *fiber.Ctx) (err error) {
	sir := new(domain.SignInRequest)
	if err = ctx.BodyParser(sir); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// 从cookies中获取constant.CookiesNameOauthAuthorizeCode，然后去缓存中获取
	var code string
	code = ctx.Cookies(constant.CookiesNameOauthAuthorizeCode)
	if code == "" {
		return ctx.Status(fiber.StatusUnauthorized).SendString("app is not authorized")
	}
	var app *modelO.OauthApplication
	var ar *domain.OauthAuthorizeRequest
	ar, err = c.popOauthAuthorizeRequest(code)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if ar == nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("app is not authorized")
	}
	app, err = service.OauthApplicationService.Instance(&modelO.OauthApplication{ClientID: ar.ClientID})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if app == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("app not found")
	}

	host := ctx.Hostname()

	var site *model.Site
	site, err = service.SiteService.GetByDomain(host)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if site == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("site not found")
	}

	var user *modelU.User
	if sir.Username != "" {
		if strings.Contains(sir.Username, "@") {
			user, err = service.UserService.GetByEmailAndDomain(sir.Username, host)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		} else {
			var account *modelA.Account
			account, err = service.AccountService.GetByUsernameAndSite(sir.Username, site.ID)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			if account != nil {
				user, err = service.UserService.GetByAccountID(account.ID)
				if err != nil {
					return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
				}
			}
		}
	}

	if user == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(sir.Password))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("password is wrong")
	}

	// 生成token放入cookies中，并加入缓存
	var token string
	token, err = c.generateAndSaveToken(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     constant.CookiesNameToken,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
	})

	// 获取account
	var account *modelA.Account
	account, err = service.AccountService.Instance(&modelA.Account{BaseModel: common.BaseModel{ID: user.AccountID}})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if account == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("account not found")
	}

	// 登录成功后，跳转到授权页面
	// 进入授权页面
	return ctx.Render("authorize", fiber.Map{
		"account": account,
		"site":    site,
		"app":     app,
		"state":   ar.State,
		"scopes":  strings.Split(ar.Scope, " "),
	}, "layouts/main")
}

func (c *authController) SignInPage(ctx *fiber.Ctx) error {
	host := ctx.Hostname()
	site, err := service.SiteService.GetByDomain(host)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.Render("sign_in", fiber.Map{
		"site": site,
	}, "layouts/main")
}

func init() {
	c := new(authController)
	Controllers = append(Controllers, c)
}

func (c *authController) InitRoute() {
	r := server.Group("/auth")
	r.Get("/sign_in", c.SignInPage)

	r.Post("/sign_in", c.SignIn)
}

func (c *authController) generateAndSaveToken(user *modelU.User) (string, error) {
	token := util.GenerateXid()
	// 存入缓存
	conn := cache.Get()
	defer conn.Close()
	// 存入user的json
	userJson, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	_, err = conn.Do("SET", constant.RedisPrefixUserToken+token, userJson)
	if err != nil {
		return "", err
	}
	_, err = conn.Do("EXPIRE", constant.RedisPrefixUserToken+token, constant.RedisUserTokenExpireTime) // 7天过期
	return token, nil
}

func (c *authController) popOauthAuthorizeRequest(code string) (*domain.OauthAuthorizeRequest, error) {
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
	return req, nil
}
