package main

import (
	"github.com/chen/download_pixiv_pic/conf"
	"github.com/chen/download_pixiv_pic/dao/dao"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"github.com/chen/download_pixiv_pic/log"
	"github.com/chen/download_pixiv_pic/pkg/browser"
	"log/slog"
)

type GlobalOptions struct {
	Rootpath string
	Auth     *sql.Auth
	DB       *dao.PsgrDB
	Br       *browser.Br
}

var globalOptions = &GlobalOptions{}

func init() {
	c := conf.GetConf()
	globalOptions.Rootpath = c.Get("DownloadControl.Path").(string)
	globalOptions.Auth = sql.DefaultAuth()
	if db, ok := dao.GetClient().Open(c.GetStringMap("postgres")); ok != "" {
		l.Send(slog.LevelWarn, "参数错误，数据库未连接", log.LogFiles|log.LogStdouts, "err", ok)
	} else {
		globalOptions.DB = db
		if err := globalOptions.DB.CreateDb(); err != nil {
			l.Send(slog.LevelWarn, "表创建失败", log.LogFiles|log.LogStdouts, "err", err)
		}
	}

}
