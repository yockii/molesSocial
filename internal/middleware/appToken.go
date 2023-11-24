package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/qscore/pkg/cache"
)

func AppAccessTokenMiddleware(isAppToken bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 从header:Authorization中获取accessToken
		accessToken := ctx.Get("Authorization")
		if accessToken == "" {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Authorization Token")
		}

		// 从redis中获取accessToken对应的appID和accountID
		// 如果redis中不存在，则返回错误
		// 如果redis中存在，则将appID和accountID存入ctx.Locals中
		// 以便后续的controller中使用

		conn := cache.Get()
		defer conn.Close()
		appID, err := redis.Uint64(conn.Do("GET", constant.RedisPrefixAccessTokenApp+accessToken))
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Authorization Token")
		}
		if !isAppToken {
			accountID, err := redis.Uint64(conn.Do("GET", constant.RedisPrefixAccessTokenAccount+accessToken))
			if err != nil {
				return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Authorization Token")
			}
			ctx.Locals(constant.CtxAccountID, accountID)
		}
		ctx.Locals(constant.CtxAppID, appID)

		return ctx.Next()
	}
}
