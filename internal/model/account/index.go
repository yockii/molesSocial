package model

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

// Account 账号，包括本站和联邦账号
type Account struct {
	common.BaseModel
	Username                string  `json:"username" gorm:"uniqueIndex:idx_account_on_username_and_site;size:100;comment:用户名"`
	SiteID                  uint64  `json:"siteId" gorm:"uniqueIndex:idx_account_on_username_and_site;comment:站点ID"`
	PrivateKey              string  `json:"privateKey" gorm:"size:100;comment:私钥"`
	PublicKey               string  `json:"publicKey" gorm:"size:100;comment:公钥"`
	RemoteUrl               string  `json:"remoteUrl" gorm:"size:100;comment:远程地址"`
	SalmonUrl               string  `json:"salmonUrl" gorm:"size:100;comment:Salmon地址"`
	HubUrl                  string  `json:"hubUrl" gorm:"size:100;comment:Hub地址"`
	CreatedAt               int64   `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt               int64   `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	Note                    string  `json:"note" gorm:"size:100;comment:备注"`
	DisplayName             string  `json:"displayName" gorm:"size:100;comment:显示名称"`
	Uri                     string  `json:"uri" gorm:"index;size:100;comment:URI"`
	Url                     string  `json:"url" gorm:"index;size:100;comment:URL"`
	AvatarMediaAttachmentID uint64  `json:"avatarMediaAttachmentId,string" gorm:"comment:头像媒体附件ID"`
	AvatarRemoteUrl         string  `json:"avatarRemoteUrl" gorm:"size:100;comment:头像远程地址"`
	AvatarBlurHash          string  `json:"avatarBlurHash" gorm:"size:100;comment:头像BlurHash"`
	AvatarUpdatedAt         int64   `json:"avatarUpdatedAt" gorm:"comment:头像更新时间"`
	HeaderMediaAttachmentID uint64  `json:"headerMediaAttachmentId,string" gorm:"comment:头部媒体附件ID"`
	HeaderRemoteUrl         string  `json:"headerRemoteUrl" gorm:"size:100;comment:头部远程地址"`
	HeaderBlurHash          string  `json:"headerBlurHash" gorm:"size:100;comment:头部BlurHash"`
	HeaderUpdatedAt         int64   `json:"headerUpdatedAt" gorm:"comment:头部更新时间"`
	SubscriptionExpiresAt   int64   `json:"subscriptionExpiresAt" gorm:"comment:订阅过期时间"`
	Locked                  *bool   `json:"locked" gorm:"comment:是否锁定"`
	LastWebfingeredAt       int64   `json:"lastWebfingeredAt" gorm:"comment:上次Webfinger时间"`
	LastSignInAt            int64   `json:"lastSignInAt" gorm:"comment:上次登录时间"`
	InboxUrl                string  `json:"inboxUrl" gorm:"size:100;comment:收件箱地址"`
	OutboxUrl               string  `json:"outboxUrl" gorm:"size:100;comment:发件箱地址"`
	SharedInboxUrl          *string `json:"sharedInboxUrl" gorm:"size:100;comment:共享收件箱地址"`
	FollowingUrl            string  `json:"followingUrl" gorm:"size:100;comment:关注地址"`
	FollowersUrl            string  `json:"followersUrl" gorm:"size:100;comment:粉丝地址"`
	Memorial                *bool   `json:"memorial" gorm:"comment:是否纪念账户(即用户是否去世了)"`
	MovedToAccountID        *uint64 `json:"movedToAccountId,string,omitempty" gorm:"index;comment:转移到的账户ID"`
	FeaturedCollectionUrl   string  `json:"featuredCollectionUrl" gorm:"size:100;comment:特色集合地址"`
	Fields                  string  `json:"fields" gorm:"size:1000;comment:用户自定义的字段"`
	ActorType               string  `json:"actorType" gorm:"size:100;comment:Actor类型"`
	Discoverable            *bool   `json:"discoverable" gorm:"comment:是否可发现"`
	AlsoKnownAs             string  `json:"alsoKnownAs" gorm:"size:1000;comment:别名"`
	SilencedAt              int64   `json:"silencedAt" gorm:"comment:沉默时间（仅粉丝可见非公开）"`
	SuspendedAt             int64   `json:"suspendedAt" gorm:"comment:暂停时间（不可登录、发送或接受消息）"`
	TrustLevel              int     `json:"trustLevel" gorm:"comment:信任等级"`
}

func (*Account) TableComment() string {
	return `账号表`
}

func (a *Account) AddRequired() string {
	if a.Username == "" || a.SiteID == 0 {
		return "username, siteId"
	}
	return ""
}

func (a *Account) CheckDuplicatedModel() database.Model {
	return &Account{
		Username: a.Username,
		SiteID:   a.SiteID,
	}
}

func init() {
	database.Models = append(database.Models, &Account{})
}
