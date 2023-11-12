package handle

import (
	"github.com/chen/download_pixiv_pic/dao/sql"
	"github.com/chen/download_pixiv_pic/pkg/browser"
	"os"
	"strings"
)

func FRTUDB(P string) {
	// 遍历rootpath + 。。。
	// 在本地查询是否存在
	var img []sql.ImageInfo
	DB.DB.Find(img)
	var d = []string{}
	for _, value := range img {
		for _, j := range value.ImageUrls {
			pl := strings.Split(j, "/")
			pd := strings.Split(pl[len(pl)-1], ".")
			pj := browser.DelSpeChar(value.ImageTitle)
			path := value.Path + pd[0] + " - " + pj + "." + pd[len(pd)-1]
			_, err := os.Stat(path)
			if os.IsExist(err) {
				d = append(d, j)
				continue
			}

		}
		if len(d) == len(value.ImageUrls) {
			continue
		}
		value.ImageUrls = d
		DB.DB.Save(&value)

	}

}
