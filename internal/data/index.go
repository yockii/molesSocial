package data

import (
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/model"
	model2 "github.com/yockii/molesSocial/internal/model/account"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/config"
)

func InitData() {
	if config.GetBool("dev") {
		InitDevData()
	}
}

func InitDevData() {
	// 初始化localhost的站点信息
	site := &model.Site{
		Name:             "本地测试站点",
		Domain:           "localhost:" + config.GetString("server.port"),
		Description:      "这是一个本地测试站点",
		ShortDesc:        "测试用",
		Thumbnail:        "https://www.baidu.com/img/flexible/logo/pc/result.png",
		Email:            "xuyuqi@gmail.com",
		Languages:        "zh-CN",
		Registrations:    &constant.BoolFalse,
		ApprovalRequired: &constant.BoolFalse,
		InvitesEnabled:   &constant.BoolFalse,
		Configuration:    "{}",
		Rules:            "[]",
	}

	service.SiteService.Add(site)

	// 初始化admin用户
	account := &model2.Account{
		Username:    "admin",
		Domain:      site.Domain,
		Note:        "管理员",
		DisplayName: "测试管理员",
	}
	service.AccountService.Add(account)

	user := &model.User{
		Email:     "xuyuqi@gmail.com",
		SiteID:    site.ID,
		Admin:     &constant.BoolTrue,
		AccountID: account.ID,
	}
	service.UserService.Add(user)
}
