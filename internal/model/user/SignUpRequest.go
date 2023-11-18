package user

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

type SignUpRequest struct {
	common.BaseModel
	UserID    uint64 `json:"userId,string" gorm:"index;comment:用户ID"`
	Reason    string `json:"reason" gorm:"size:100;comment:原因"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`
}

func (*SignUpRequest) TableComment() string {
	return `注册请求表`
}

func (u *SignUpRequest) AddRequired() string {
	if u.UserID == 0 {
		return "userId"
	}
	return ""
}

func (u *SignUpRequest) CheckDuplicatedModel() database.Model {
	return &SignUpRequest{
		UserID: u.UserID,
	}
}

func init() {
	database.Models = append(database.Models, &SignUpRequest{})
}
