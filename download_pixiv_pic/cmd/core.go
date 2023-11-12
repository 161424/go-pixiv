package main

import (
	"fmt"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"github.com/chen/download_pixiv_pic/log"
	"github.com/chen/download_pixiv_pic/pkg/browser"
	"log/slog"
	"os"
)

func GetPixivRanking(mode, date, content string, filter string, page int) (*browser.RpJs, string) {
	return globalOptions.Br.GetPixivRanking(mode, date, content, filter, page)
}

func ProcessImage(pgcount int, imageId, rootpath string, ipi *inputInfo) error {
	var inDb = false
	var dbImgInfo sql.ImageInfo
	var fileName string

	// GetImageInfo
	imgInfo, err := globalOptions.Br.GetImageInfo(imageId)
	if err != nil {
		l.Send(slog.LevelError, "获取图片信息失败", log.LogFiles|log.LogStdouts, "err", err)
		return err
	}
	imgInfo.ImageCount = pgcount
	// CreateDB or pass
	if ipi.Mode == "rk" {
		err, dbImgInfo = globalOptions.DB.SelectImageByImageId(imageId)
	} else if ipi.Mode == "mb" {

	}

	// 只有检测到真实存在的文件，num才会+1
	if err == nil && len(dbImgInfo.SavePath) != 0 {
		dirPath, err := os.ReadDir(dbImgInfo.SavePath)
		if err != nil {
			return err
		} else {
			if len(dirPath) == pgcount {
				inDb = true
			}
		}
	}

	if inDb == true {
		fmt.Printf("Already downloaded in DB: %s\n", imageId)
		l.Send(slog.LevelInfo, fmt.Sprintf("Already downloaded in DB: %s\n", imageId), log.LogFiles|log.LogStdouts)
		return nil
	}

	if fileName, err = browser.MakeFilename(imgInfo, rootpath, ipi.User, ipi.Mode, ipi.Date, ipi.Content); err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("image文件夹创建失败，err=%s", err), log.LogFiles|log.LogStdouts)
		return err
	} else {
		fileName += string(os.PathSeparator)
		imgInfo.SavePath = fileName
	}
	// 暂不设置下载图像大小，默认默认
	browser.PrintInfo(imgInfo)
	globalOptions.Br.DownloadImage(imgInfo, fileName)
	if ok := globalOptions.DB.SaveImageId(imgInfo); ok != nil {
		l.Send(slog.LevelError, fmt.Sprintf("图片信息存储失败，err=%s", err), log.LogFiles|log.LogStdouts)
	}
	return nil
}
