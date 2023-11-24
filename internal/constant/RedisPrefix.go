package constant

const (
	RedisPrefix                   = "moleSocial"
	RedisPrefixOauthAuthorizeInfo = RedisPrefix + ":oauth:authorize:"
	RedisPrefixUserToken          = RedisPrefix + ":user:token:"
	RedisPrefixAccessCode         = RedisPrefix + ":access:code:"
	RedisPrefixAccessTokenAccount = RedisPrefix + "accessToken:account:"
	RedisPrefixAccessTokenApp     = RedisPrefix + "accessToken:app:"
)

const (
	RedisUserTokenExpireTime = 60 * 60 * 24 * 7 // 7 days
)
