package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/domain"
	"github.com/yockii/molesSocial/internal/model/account"
	modelU "github.com/yockii/molesSocial/internal/model/user"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/cache"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

type authController struct{}

func (c *authController) SignIn(ctx *fiber.Ctx) (err error) {
	sir := new(domain.SignInRequest)
	if err = ctx.BodyParser(sir); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	host := ctx.Hostname()
	var user *modelU.User
	if sir.Username != "" {
		var account *model.Account
		account, err = service.AccountService.GetByUsernameAndDomain(sir.Username, host)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		if account != nil {
			user, err = service.UserService.GetByAccountID(account.ID)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}
	} else if sir.Email != "" {
		user, err = service.UserService.GetByEmailAndDomain(sir.Email, host)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
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
	return ctx.Redirect("/oauth/authorize")
}

func (c *authController) SignInPage(ctx *fiber.Ctx) error {
	host := ctx.Hostname()
	site, err := service.SiteService.GetByDomain(host)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	layout := ""
	if site.LayoutTemplateName != "" {
		layout = site.Domain + "/" + site.LayoutTemplateName
	}

	return ctx.Render(site.Domain+"/sign_in", fiber.Map{
		"site": site,
	}, layout)
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
