package main

import (
	"context"
	"fmt"
	"github.com/chen/download_pixiv_pic/cmd/utils"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"github.com/chen/download_pixiv_pic/log"
	"github.com/spf13/cobra"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

//var Rootpath = (conf.ConfigData["DownloadControl"]["Path"]).(string)

//const Spt = string(os.PathSeparator)

var IpI = &inputInfo{
	User:       "rank",
	Mode:       "daily",
	Date:       time.Now().AddDate(0, 0, -1).Format("20060102"),
	RankTop:    10,
	NotSeSe:    true,
	SkipUgoira: true,
	SkipManga:  true,
	Content:    "all",
}

type inputInfo struct {
	User            string
	NotSeSe         bool
	Mode            string // day,week,....
	Date            string // 初始下载时间
	RankTop         int
	Content         string //
	SkipIllust      bool
	SkipIllusts     bool
	SkipUgoira      bool
	SkipManga       bool
	SkipNovel       bool
	IncludeSkipTime bool
}

var l *log.Logs

// 日 周 月 。。。
var cmdRank = &cobra.Command{
	Use:   "rk [flags]",
	Short: "根据排行进行下载",
	Long:  "可以下载日排行，周排行，月排行的图片",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//  判断日期是否正确
		if err := utils.CheckTime(IpI.Date); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		var wg sync.WaitGroup
		cancelCtx, cancel := context.WithCancel(ctx)
		defer func() {
			cancel()
			wg.Wait()
		}()

		// TODO
		_ = cancelCtx
		// termstatus?  更新终端显示

		l = log.NewSlogGroup("Rank")
		DownloadRank(ctx, globalOptions, args)
		return

	},
}

func init() {
	cmdRoot.AddCommand(cmdRank)
	f := cmdRank.Flags()
	f.BoolVar(&IpI.NotSeSe, "R!*", IpI.NotSeSe, "是否是R18")
	f.StringVarP(&IpI.Date, "data", "d", IpI.Date, "设置下载的开始日期，默认是昨天")
	f.StringVarP(&IpI.Mode, "mode", "m", IpI.Mode, "设置需要下载的类型，有daily,weekly,monthly,rookie,original,male,female")
	f.IntVarP(&IpI.RankTop, "top", "t", IpI.RankTop, "下载前top个图片,max=500")
	f.StringVarP(&IpI.Content, "content", "c", IpI.Content, "设置下载的内容，有综合，插画，动画，漫画，小说")
	f.BoolVarP(&IpI.SkipIllust, "skipIllus", "i", IpI.SkipIllust, "是否跳过插图")
	f.BoolVarP(&IpI.SkipIllusts, "skipIlluss", "I", IpI.SkipIllusts, "是否跳过有多个插图的")
	f.BoolVarP(&IpI.SkipUgoira, "skipUgoira", "u", IpI.SkipUgoira, "是否跳过动画")
	f.BoolVarP(&IpI.SkipManga, "skipManga", "M", IpI.SkipManga, "是否跳过漫画")
	f.BoolVarP(&IpI.IncludeSkipTime, "include", "C", IpI.IncludeSkipTime, "下载的top里面是否包含的次数")
	//f.BoolVar(&IpI.SkipNovel, "skipNovel", false, "是否跳过小说")

}

func DownloadRank(ctx context.Context, gopts *GlobalOptions, args []string) {
	rootPath := gopts.Rootpath
	//l.Send(slog.LevelInfo, "Download Ranking by Post ID mode (15).", log.LogFiles|log.LogStdouts)
	//validModes := []string{"daily", "weekly", "monthly", "rookie", "original", "male", "female"}
	//validContents := []string{"all", "illust", "ugoira", "manga"}

	utils.QuoteOrCreateFile(rootPath)

	l.Send(slog.LevelInfo, fmt.Sprintf(
		"Downloading Setting View:R18 %v,Data %s, Mode %s, TopPage %d, Content %s, SkipIllust %v, SkipUgoira %v, SkipManga %v, SkipNovel %v",
		IpI.NotSeSe, IpI.Date, IpI.Mode, IpI.RankTop, IpI.Content, IpI.SkipIllust, IpI.SkipUgoira, IpI.SkipManga, IpI.SkipNovel), log.LogStdouts|log.LogFiles)

	var InSCount = 0
	var DwCount = 0
	var DwCounts = 0
	var AllCount = 0
	var AllCounts = 0
	var SkipNum = 0

	var page = 1
	var total = 50
	wd := sync.Once{}
bk:
	for page*50 <= total {

		l.Send(slog.LevelInfo, fmt.Sprintf("** Page %d **", page), 3)
		ranks, output := GetPixivRanking(IpI.Mode, IpI.Date, IpI.Content, "", page)
		l.Send(slog.LevelInfo, output, log.LogFiles|log.LogStdouts)

		if ranks == nil {
			gopts.Auth.Status = sql.Fail
			gopts.Auth.Output = output
			break
		}
		wd.Do(func() {
			l.Send(slog.LevelInfo, fmt.Sprintf("*Mode :%s", ranks.Mode), log.LogStdouts)
			l.Send(slog.LevelInfo, fmt.Sprintf("*Content :%s", ranks.Content), log.LogStdouts)
			l.Send(slog.LevelInfo, fmt.Sprintf("*Total :%d", ranks.RankTotal), log.LogStdouts)
			total = ranks.RankTotal
		})

		AllCounts += len(ranks.Contents)

	skip:
		for k := 0; k < len(ranks.Contents); k++ {

			if IpI.IncludeSkipTime {
				InSCount = AllCount
			} else {
				InSCount = DwCount
			}

			if InSCount == IpI.RankTop {
				break bk
			}

			post := ranks.Contents[k]
			AllCount += 1
			switch post.IllustType {
			case "0":
				if IpI.SkipIllust || (IpI.SkipIllusts && post.IllustPageCount != "1") {
					break skip
				}
			case "1":
				if IpI.SkipManga {
					break skip
				}
			case "2":
				if IpI.SkipUgoira {
					break skip
				}
			}

			SkipNum++
			imageId := fmt.Sprintf(strconv.FormatFloat(post.IllustId, 'f', 0, 64))
			l.Send(slog.LevelInfo, fmt.Sprintf("{%d}.Image id %s", InSCount, imageId), log.LogStdouts)
			pgCount, _ := strconv.Atoi(post.IllustPageCount)

			result := ProcessImage(pgCount, imageId, rootPath, IpI)

			if result == nil {
				DwCount += 1
				DwCounts += pgCount
			}
			l.Send(slog.LevelDebug, fmt.Sprintf("result:%v", result), log.LogStdouts)

		}

		// 判断是否是最后一页
		if ranks.NextPage == false {
			break
		}

	}

	feedback := fmt.Sprintf("需下载前top%d，成功下载前top%d，下载失败%d，过滤个%d，成功下载%d个图片", IpI.RankTop, DwCount, InSCount-DwCount, AllCount-InSCount, DwCounts)
	l.Send(slog.LevelInfo, feedback, log.LogFiles|log.LogStdouts)
	globalOptions.Auth.Output = feedback
	globalOptions.DB.SaveAuthDownloadRecord(globalOptions.Auth)

}
