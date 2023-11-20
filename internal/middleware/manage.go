package middleware

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/yockii/molesSocial/internal/constant"
	model "github.com/yockii/molesSocial/internal/model/manage"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/cache"
	"github.com/yockii/qscore/pkg/config"
	"strconv"
	"strings"
)

// NeedAuthorization 需要授权的中间件
// code: 空 或 anon 不需要授权
// code: user 需要用户授权
// code: 其他 需要用户授权并且需要对应的权限
func NeedAuthorization(codes ...string) fiber.Handler {
	if config.GetBool("isDev") {
		return func(ctx *fiber.Ctx) error {
			return ctx.Next()
		}
	}

	for _, code := range codes {
		code = strings.ToLower(code)
		if code == "" || code == "anon" {
			return func(ctx *fiber.Ctx) error {
				return ctx.Next()
			}
		}
	}

	codeMap := make(map[string]bool)
	for _, code := range codes {
		codeMap[code] = true
	}

	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(constant.JwtSecret),
		ContextKey:    "jwt-subject",
		SigningMethod: "HS256",
		TokenLookup:   "header:Authorization,cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).SendString("Invalid Authorization Token")
			} else {
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Authorization Token")
			}
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			jwtToken := c.Locals("jwt-subject").(*jwt.Token)
			claims := jwtToken.Claims.(jwt.MapClaims)
			uid := claims[constant.JwtClaimManagerId].(string)
			sid := claims[constant.JwtClaimSessionId].(string)
			tenantId, hasTenantId := claims[constant.JwtClaimTenantId].(string)

			conn := cache.Get()
			defer func(conn redis.Conn) {
				_ = conn.Close()
			}(conn)
			sessionKey := constant.RedisPrefixSessionId + sid
			cachedUid, err := redis.String(conn.Do("GET", sessionKey))
			if err != nil {
				if !errors.Is(err, redis.ErrNil) {
					logrus.Errorln(err)
				}
				return c.Status(fiber.StatusUnauthorized).SendString("token information has expired")
			}
			if cachedUid != uid {
				return c.Status(fiber.StatusUnauthorized).SendString("token信息不正确")
			}

			// 获取用户角色
			userRolesKey := constant.RedisPrefixUserRoles + uid
			roleIds, _ := redis.Uint64s(conn.Do("SMEMBERS", userRolesKey))
			if len(roleIds) == 0 {
				// 获取该用户的角色id存入缓存
				userId, _ := strconv.ParseUint(uid, 10, 64)
				if userId == 0 {
					return c.Status(fiber.StatusUnauthorized).SendString("token信息已失效")
				}
				var roles []*model.Role
				roles, err = service.ManagerService.ManagerRoles(userId)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString("系统错误")
				}
				for _, role := range roles {
					// 缓存用户的角色
					_, _ = conn.Do("SADD", userRolesKey, role.ID)
				}
			}
			_, _ = conn.Do("EXPIRE", userRolesKey, 3*24*60*60)

			hasAuth := false
			if _, ok := codeMap["user"]; ok {
				hasAuth = true
			}
			for _, roleId := range roleIds {
				roleResourceKey := fmt.Sprintf("%s:%d", constant.RedisPrefixRoleResourceCode, roleId)
				cachedCodes, _ := redis.Strings(conn.Do("SMEMBERS", roleResourceKey))
				if len(cachedCodes) == 0 {
					// 缓存没有，那么就去数据库取出来放进去
					cachedCodes, err = service.RoleService.ResourceCodes(roleId)
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).SendString("系统错误")
					}
					for _, resourceCode := range cachedCodes {
						rc := resourceCode
						_, _ = conn.Do("SADD", roleResourceKey, rc)
						if _, ok := codeMap[rc]; ok {
							hasAuth = true
						} else {
							for _, code := range codes {
								if strings.HasPrefix(code, rc+":") {
									hasAuth = true
									break
								}
							}
						}
					}
				}
				_, _ = conn.Do("EXPIRE", roleResourceKey, 3*24*60*60)
				if hasAuth {
					break
				}
				for _, resourceCode := range cachedCodes {
					rc := resourceCode
					if _, ok := codeMap[rc]; ok {
						hasAuth = true
						break
					} else {
						for _, code := range codes {
							if strings.HasPrefix(code, rc+":") {
								hasAuth = true
								break
							}
						}
					}
				}
				if hasAuth {
					break
				}
			}

			if !hasAuth {
				return c.Status(fiber.StatusUnauthorized).SendString("无权限")
			}

			// 有权限，那么就把用户信息放到上下文中
			c.Locals(constant.JwtClaimManagerId, uid)
			// 如果有租户，则租户信息也放入
			if hasTenantId {
				c.Locals(constant.JwtClaimTenantId, tenantId)
			}
			// token续期
			_, _ = conn.Do("EXPIRE", sessionKey, config.GetInt("userTokenExpire"))
			return c.Next()
		},
	})
}
