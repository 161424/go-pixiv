package handle

import (
	"fmt"
	sql2 "github.com/chen/download_pixiv_pic/dao/dao"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"time"
)

var DB = sql2.GetClient()

func MainLoop() {
	BuildMenu()
	rootpath := Rootpath

	for {

		selection := menu()
		if selection == "16" {
			CostTime(DownloadRank, rootpath)
		}
		if selection == "-1" {
			CostTime(FRTUDB, "")
		}

		if selection == "2" {
			CostTime(DownLoadByAuth, "")
		}

		if Menu[selection] != "" {
			DB.DB.Create(Auth)
			Auth = sql.InitAuth()
		}

	}
}

func CostTime(f func(string2 string), path string) {
	time2 := time.Now()
	f(path)
	time3 := time.Now()
	Auth.CostTime = time3.Sub(time2).Minutes()
}

func menu() string {
	var sel string
	fmt.Printf("\033[1;32;40m%s\033[0m\n", "── Pixiv ────────────────────────────────────────────────────────────")
	for key, value := range Menu {
		fmt.Printf(" %s. %s\n", key, value)
	}

	//fmt.Printf("\033[1;32;40m%s\033[0m\n", "── FANBOX ────────────────────────────────────────────────────────────")

	fmt.Scanln(&sel)
	return sel

}
