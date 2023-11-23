package data

import (
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/model"
	model2 "github.com/yockii/molesSocial/internal/model/account"
	modelU "github.com/yockii/molesSocial/internal/model/user"
	"github.com/yockii/molesSocial/internal/service"
	"github.com/yockii/qscore/pkg/config"
	"github.com/yockii/qscore/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func InitData() {
	if config.GetBool("dev") {
		InitDevData()
	}
}

func InitDevData() {
	// 开启事务
	tx := database.DB.Begin()

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

	duplicated, err := service.SiteService.Add(site, tx)
	if err != nil {
		return
	}
	if duplicated {
		// 获取site
		site, err = service.SiteService.GetByDomain(site.Domain)
		if err != nil {
			return
		}
	}

	// 初始化admin用户
	account := &model2.Account{
		Username:    "admin",
		SiteID:      site.ID,
		Note:        "管理员",
		DisplayName: "测试管理员",
	}
	duplicated, err = service.AccountService.Add(account, tx)
	if err != nil {
		return
	}
	if duplicated {
		// 获取account
		account, err = service.AccountService.GetByUsernameAndSite(account.Username, site.ID)
		if err != nil {
			return
		}
	}

	pwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	user := &modelU.User{
		Email:             "xuyuqi@gmail.com",
		SiteID:            site.ID,
		Admin:             &constant.BoolTrue,
		AccountID:         account.ID,
		EncryptedPassword: string(pwd),
	}
	duplicated, err = service.UserService.Add(user, tx)
	if err != nil {
		return
	}
	if duplicated {
		// 获取user
		user, err = service.UserService.GetByAccountID(account.ID)
		if err != nil {
			return
		}
	}
	tx.Commit()
}
