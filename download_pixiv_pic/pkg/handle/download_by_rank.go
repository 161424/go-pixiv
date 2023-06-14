package handle

import (
	"fmt"
	"github.com/chen/download_pixiv_pic/common/conf"
	"github.com/chen/download_pixiv_pic/database/sql"
	"github.com/chen/download_pixiv_pic/pkg/Browser"
	"os"
	"strconv"
	"strings"
	"time"
)

var Rootpath = (conf.ConfigData["DownloadControl"]["Path"]).(string)

const Spt = string(os.PathSeparator)

var InP inputInfo

type inputInfo struct {
	Mode    string
	Content string
	Date    string
}

var Auth = sql.DefaultAuth
var Menu = make(map[string]string, 16)

func BuildMenu() {
	//Menu = append(Menu, []string{"1": "12"})
	Menu["1"] = "Download by ImageId"
	Menu["2"] = "Download by Author"
	Menu["16"] = "Download by Rank"
	Menu["-1"] = "From RootPath to Update DB"
}

func DownloadRank(string2 string) {
	rootpath := Rootpath
	fmt.Println("Download Ranking by Post ID mode (15).")

	validModes := []string{"daily", "weekly", "monthly", "rookie", "original", "male", "female"}
	validContents := []string{"all", "illust", "ugoira", "manga"}

	var mode string
	var contents string
	var date string

md:
	for {
		fmt.Printf("Valid Modes are: %s\n", Map2String(validModes, ", "))
		fmt.Print("Mode: ")
		fmt.Scanln(&mode)
		for _, i := range validModes {
			if strings.ToLower(mode) == i {
				rootpath += Spt
				rootpath += i
				InP.Mode = i
				break md
			}
		}
		fmt.Printf("\033[1;31;40m%s\033[0m\n", "Invalid mode.")
	}
ty:
	for {
		fmt.Printf("Valid Content Types are: %s\n", Map2String(validContents, ", "))
		fmt.Print("Type: ")
		fmt.Scanln(&contents)
		for _, i := range validContents {
			if strings.ToLower(contents) == i {
				contents = strings.ToLower(contents)
				InP.Content = i
				break ty
			}
		}
		fmt.Println("Invalid Content Type.")
	}
de:
	for {
		fmt.Println("Specify the ranking date, valid type is YYYYMMDD (default: today)")
		fmt.Print("Date: ")
		fmt.Scanln(&date)
		if date != "" {
			_date, err := time.Parse("20060102", date)
			if err != nil {
				continue
			}
			rootpath += Spt
			date = _date.Format("20060102")
			rootpath += _date.Format("2006-01-02")
			InP.Date = date
			break de
		}
		fmt.Println("Invalid Date.")
	}

	Browser.CheckPathIsExit(rootpath)

	sp, ep, num := GetStartAndEndNumber()
	//fmt.Println(InP)
	ProcessRanking(mode, contents, rootpath, date, sp, ep, num)

}

func GetStartAndEndNumber() (int, int, int) {

	var mode string
	var sp, ep, num int
	fmt.Print("Start Page (default=1): ")
	fmt.Scanln(&mode)
	if mode == "" {
		fmt.Print("Start Page = 1 ")
		sp = 1
	} else {
		sp, _ = strconv.Atoi(mode)
	}

	fmt.Printf("请输入下载的数量或者页数  1.个数(default=10) or 2.结束页数(default=sp+1): ")
	fmt.Scanln(&mode)
	tep := strings.Split(mode, ".")
	if tep[0] == "1" {
		if len(tep) == 1 || tep[1] == "" {
			num = 10
		} else {
			num, _ = strconv.Atoi(tep[1])
		}
	} else if tep[0] == "2" || tep[1] == "" {
		if len(tep) == 1 {
			//fmt.Println("End Page = 2 ")
			ep = 2
		} else {
			ep, _ = strconv.Atoi(tep[1])
		}
		if ep <= sp {
			fmt.Println("ep err. Make ep = sp + 1")
			ep = sp + 1
		}
	} else {

	}

	fmt.Printf("sp = %d, ep = %d, num = %d\n", sp, ep, num)
	return sp, ep, num
}

func ProcessRanking(mode, concent, rootpath, date string, sp, ep, num int) {
	fmt.Printf("Processing Pixiv Ranking: %s.\n", mode)
	var br *Browser.Br
	var i = 1
	var YesCount = 0
	var AllCount = 0
	var AllImage = 0
Nm:
	for {
		fmt.Printf("Pixiv Ranking: %s page %d and date %s\n", mode, sp, date)
		ranks := br.GetPixivRanking(mode, date, concent, "", sp)
		if ranks == nil {
			Auth.Status = 1
			Auth.Output = "Cookie 失效或其他莫名错误"
			break
		}
		fmt.Printf("Mode :%s\n", ranks.Mode)
		fmt.Printf("Total :%.0f\n", ranks.RankTotal)
		fmt.Printf("Next Page :%.0f\n", ranks.NextPage)

		if _, isok := (ranks.NextDate).(bool); isok {
			fmt.Printf("Next Page :%s\n", strconv.FormatBool((ranks.NextDate).(bool)))
		} else {
			fmt.Printf("Next Date :%s\n", ranks.NextDate)
		}
		fmt.Println(ranks.Contents)
		AllCount += len(ranks.Contents)

		for _, post := range ranks.Contents {
			imageId := fmt.Sprintf(strconv.FormatFloat((post["illust_id"]).(float64), 'f', 0, 64))
			fmt.Printf("#%d. Image id %s\n", i, imageId)
			result := br.ProcessImage(imageId, rootpath, "rk")
			if result == "YES" {
				YesCount += 1
			}
			AllImage += Browser.ImageCount
			fmt.Printf(result)
			//fmt.Println(post)
			i += 1
			if num > 0 && num < YesCount {
				fmt.Printf("Reach max download num %d", num)
				break Nm
			}
		}

		sp += 1
		if ep > 0 && sp > ep {
			fmt.Printf("Reach max page %d", ep)
			break
		}
	}
	Auth.Output = fmt.Sprintf("共完成从 sp-%d 到 ep-%d 共 %d 页 %d 项 %d 个图片", sp, ep, ep-sp, AllCount, AllImage)
}
