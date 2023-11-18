package ap

type Actor struct {
	Base
	Type              string          `json:"type"`              // Person
	Id                string          `json:"id"`                // ID url
	Url               string          `json:"url,omitempty"`     // 主页url，如果与ID不同则需要
	Followers         string          `json:"followers"`         // 粉丝url
	Following         string          `json:"following"`         // 关注url
	Inbox             string          `json:"inbox"`             // 收件箱
	Outbox            string          `json:"outbox"`            // 发件箱
	PreferredUsername string          `json:"preferredUsername"` // 首选用户名，不保证唯一性
	Name              string          `json:"name"`              // 昵称或显示名
	Summary           string          `json:"summary"`           // 简介
	Liked             string          `json:"liked"`
	Streams           []string        `json:"streams"`
	Endpoints         []ActorEndpoint `json:"endpoints,omitempty"` // 其他服务端点
	Icon              []string        `json:"icon"`                // 个人资料、头像等图片链接
}

type ActorEndpoint struct {
	ProxyUrl                   string `json:"proxyUrl,omitempty"`                   // 端点url，客户端会POST一个x-www-form-urlencoded请求，并将id作为参数
	OauthAutohrizationEndpoint string `json:"oauthAuthorizationEndpoint,omitempty"` // Bearer token验证，则使用该端点指定url来获取授权
	OauthTokenEndpoint         string `json:"oauthTokenEndpoint,omitempty"`         // Bearer token验证，则使用该端点指定url来获取token
	ProvideClientKey           string `json:"provideClientKey,omitempty"`           // 公钥授权
	SignClientKey              string `json:"signClientKey,omitempty"`              // ？？授权
	SharedInBox                string `json:"sharedInbox,omitempty"`                // 共享收件箱
}
