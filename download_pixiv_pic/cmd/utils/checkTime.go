package utils

import (
	"github.com/chen/download_pixiv_pic/log"
	"log/slog"
	"time"
)

var l = log.NewDefaultSlog()

func CheckTime(times string) error {

	now := time.Now()
	inputTime, err := time.Parse("20060102", times)

	if err != nil {
		l.Send(slog.LevelError, "时间格式错误", log.LogFiles|log.LogStdouts)
		return err

	}
	if now.Sub(inputTime) < time.Hour*24 {
		l.Send(slog.LevelError, "日期错误，日期还未到达", log.LogFiles|log.LogStdouts)
		return err
	}

	return nil
}
