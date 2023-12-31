package model

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
	"gorm.io/gorm"
)

// Status 用户近况表(推文/嘟文)
type Status struct {
	common.BaseModel
	AccountID          uint64         `json:"accountId,string" gorm:"index;comment:推文/嘟文作者ID"`
	SiteID             uint64         `json:"site_id" gorm:"index;comment:推文/嘟文作者域名"`
	Uri                string         `json:"uri" gorm:"size:100;comment:推文/嘟文地址"`
	Text               string         `json:"text" gorm:"comment:推文/嘟文内容"`
	CreatedAt          int64          `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt          int64          `json:"updated_at" gorm:"autoUpdateTime:milli"`
	InReplyToID        uint64         `json:"in_reply_to_id,string,omitempty" gorm:"index;comment:回复的推文/嘟文ID"`
	InReplyToUri       string         `json:"in_reply_to_uri" gorm:"size:100;comment:回复的推文/嘟文地址"`
	InReplyToAccountID uint64         `json:"in_reply_to_account_id,string,omitempty" gorm:"index;comment:回复的推文/嘟文作者ID"`
	ReblogOfID         uint64         `json:"reblog_of_id,string,omitempty" gorm:"index;comment:转发的推文/嘟文ID"`
	ReblogOfUri        string         `json:"reblog_of_uri" gorm:"size:100;comment:转发的推文/嘟文地址"`
	ReblogOfAccountID  uint64         `json:"reblogOfAccountId,string,omitempty" gorm:"index;comment:转发的推文/嘟文作者ID"`
	ApplicationID      uint64         `json:"applicationId,string" gorm:"index;comment:应用ID"`
	Visibility         int            `json:"visibility" gorm:"comment:可见性"`
	Sensitive          *bool          `json:"sensitive" gorm:"comment:是否敏感"`
	SpoilerText        string         `json:"spoiler_text" gorm:"size:100;comment:内容警告"`
	Local              *bool          `json:"local" gorm:"comment:是否本地"`
	PollID             uint64         `json:"pollId,string,omitempty" gorm:"comment:投票ID"`
	ConversationID     uint64         `json:"conversationId,string,omitempty" gorm:"comment:对话ID"`
	Language           string         `json:"language" gorm:"size:100;comment:语言"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (*Status) TableComment() string {
	return `用户近况表(推文/嘟文)`
}

func (s *Status) AddRequired() string {
	if s.Text == "" || s.AccountID == 0 || s.SiteID == 0 {
		return "text, accountId, domain"
	}
	return ""
}

func init() {
	database.Models = append(database.Models, &Status{})
}
