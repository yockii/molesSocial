package constant

const (
	RedisPrefix                   = "moleSocial"
	RedisPrefixOauthAuthorizeInfo = RedisPrefix + ":oauth:authorize:"
	RedisPrefixUserToken          = RedisPrefix + ":user:token:"
)

const (
	RedisUserTokenExpireTime = 60 * 60 * 24 * 7 // 7 days
)
