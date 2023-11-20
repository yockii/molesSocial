package constant

const (
	RedisPrefix                   = "moleSocial"
	RedisPrefixOauthAuthorizeInfo = RedisPrefix + ":oauth:authorize:"
	RedisPrefixUserToken          = RedisPrefix + ":user:token:"

	RedisPrefixSessionId        = RedisPrefix + ":session:id:"
	RedisPrefixUserRoles        = RedisPrefix + ":user:roles:"
	RedisPrefixRoleResourceCode = RedisPrefix + ":role:resource:code:"
)

const (
	RedisUserTokenExpireTime = 60 * 60 * 24 * 7 // 7 days
)
