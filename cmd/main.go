package main

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/molesSocial/internal/constant"
	"github.com/yockii/molesSocial/internal/controller"
	"github.com/yockii/molesSocial/internal/data"
	"github.com/yockii/qscore/pkg/cache"
	"github.com/yockii/qscore/pkg/config"
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/task"
	"github.com/yockii/qscore/pkg/util"
)

var VERSION = "0.0.1"

func init() {
	constant.Version = VERSION
}

func main() {
	// 初始化日志
	config.InitialLogger()

	cache.InitWithDefault()
	defer cache.Close()

	// 初始化数据库
	database.Initial()
	defer database.Close()

	// 雪花算法初始节点信息
	_ = util.InitNode(0)

	// 初始化数据库，需要将model加入到common.Models中
	_ = database.AutoMigrate()
	data.InitData()

	// 启动定时任务
	task.Start()
	defer task.Stop()

	server.InitServer()

	// 初始化路由
	controller.InitRouter() // 自行编码，仅为了引用出controller包，也可尝试直接导入 _ "github.com/xxx/xxx/controller"

	// 启动服务
	for {
		err := server.Start()
		if err != nil {
			logger.Errorln(err)
		}
	}
}
