package browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chen/download_pixiv_pic/cmd/utils"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"github.com/chen/download_pixiv_pic/log"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

type Br struct {
	//Url []string
}

var ImageCount = 0
var l *log.Logs
var refere = "https://www.pixiv.net"

type RpJs struct {
	Mode      string    `json:"mode"`
	CurrDate  string    `json:"date"`
	NextDate  any       `json:"next_date"`
	PrevDate  any       `json:"prev_date"`
	CurrPage  any       `json:"page"`
	NextPage  any       `json:"next"`
	PrevPage  any       `json:"prev"`
	RankTotal int       `json:"rank_total"`
	Contents  []content `json:"contents"`
	Content   string    `json:"content"`
	Url       string    `json:"url"`
}

// 对于跨包，子特征名大小写不重要，但是变量名一定要大写才能被解析到

type content struct {
	Title             string         `json:"title"`
	Tags              []string       `json:"tags"`
	IllustType        string         `json:"illust_type"`
	IllustBookStyle   string         `json:"illust_book_style"`
	IllustPageCount   string         `json:"illust_page_count"`
	UserName          string         `json:"user_name"`
	IllustContentType map[string]any `json:"illust_content_type"`
	IllustId          float64        `json:"illust_id"`
	UserId            float64        `json:"user_id"`
}

type Cont struct {
	Timestamp string
	Illust    map[string]map[string]interface{}
	User      map[string]map[string]interface{}
}

func init() {
	l = log.NewSlogGroup("Browser")
}

func (br *Br) GetPixivPage(url, referer string, enable_cache bool) (RJ *RpJs, s string) {
	if referer == "" {
		referer = refere
	}
	url = br.fixUrl(url, true)
	rb := GetPixivPage(url, referer)
	l.Send(slog.LevelDebug, url, log.LogFiles|log.LogStdouts)
	s = ""
	if strings.Contains(string(rb), "Just a moment...") {
		s = "受到了CloudFlare 5s盾的管控， 请检查cookie和梯子是否有效"
		return
	}
	//fmt.Printf()
	if strings.Contains(string(rb), "不在排行榜统计范围内") {
		s = "不在排行榜统计范围内"
		return
	}

	if err := json.Unmarshal(rb, &RJ); err != nil {
		s = "解析错误，请检查参数是否匹配"
		return
	}

	return
}

func (br *Br) fixUrl(url string, useHttps bool) string {
	if !strings.HasPrefix(url, "http") {
		if !strings.HasPrefix(url, "/") {
			url = "/" + url
		}
		if useHttps {
			return "https://www.pixiv.net" + url
		} else {
			return "http://www.pixiv.net" + url
		}
	}
	return url
}

func (br *Br) GetPixivRanking(mode, date, content string, filter string, page int) (*RpJs, string) {
	url := fmt.Sprintf("https://www.pixiv.net/ranking.php?mode=%s", mode)
	if len(date) > 0 {
		url = fmt.Sprintf("%s&date=%s", url, date)
	}
	if len(content) > 0 && content != "all" {
		url = fmt.Sprintf("%s&content=%s", url, content)
	}
	url = fmt.Sprintf("%s&p=%d&format=json", url, page)

	df, err := br.GetPixivPage(url, "", false)

	if err != "" {
		return nil, err
	}

	return df, ""
}

func (br *Br) GetImagePage(imageId string) (*sql.ImageInfo, error) {
	l.Send(slog.LevelDebug, fmt.Sprintf("Getting image page: %s", imageId), log.LogStdouts)
	url := fmt.Sprintf("https://www.pixiv.net%s/artworks/%s/", "", imageId)
	resp := GetPixivPage(url, "")
	df := bytes.NewReader(resp)

	doc, err := goquery.NewDocumentFromReader(df)
	if err != nil {
		return nil, err
	}
	imgInfo, err := br.PixivImage(imageId, doc)
	if err != nil {
		return nil, err
	}
	if imgInfo.ImageMode == "ugoira_view" {
		ugoiraMetaUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/ugoira_meta", imageId)
		// ugoira 代表是动图，暂不进行过多考虑
		fmt.Println(ugoiraMetaUrl)
	}
	return imgInfo, nil
}

func (br *Br) GetImageInfo(imageId string) (imageInfo *sql.ImageInfo, err error) {
	//todo 调用数据库

	//referer := "https://www.pixiv.net/artworks/" + imageId
	//filename := "no-filename-" + strings.Split(imageId, "\n")[0] + ".tmp"
	//
	//fmt.Printf("referer: %s\n", referer)
	//fmt.Printf("filename: %s\n", filename)
	//fmt.Printf("Processing Image Id: %s\n", imageId)

	if imageInfo, err = br.GetImagePage(imageId); err != nil {
		return nil, err
	}
	imageInfo.ImageId = imageId
	return
}

func (br *Br) DownloadImage(imgInfo *sql.ImageInfo, path string) {
	var wg = sync.WaitGroup{}
	var _p string
	for index, imgUrl := range imgInfo.ImageUrls {
		l.Send(slog.LevelInfo, fmt.Sprintf(" - [%d/%d]Image URL : %s", index+1, imgInfo.ImageCount, imgUrl), log.LogFiles|log.LogStdouts)
		rp := strings.Split(imgUrl, ".") // 图片类型
		It := utils.DelSpeChar(imgInfo.ImageTitle)
		if len(imgInfo.ImageUrls) > 1 {
			_p = fmt.Sprintf("[%s_p%d]%s.%s", imgInfo.ImageId, index, It, rp[len(rp)-1])
		} else {
			_p = fmt.Sprintf("[%s]%s.%s", imgInfo.ImageId, It, rp[len(rp)-1])
		}

		filePath := path + _p
		wg.Add(1)
		go func(imgurl string, filepath string, imginfo sql.ImageInfo) {
			result := br.DownLoadImage(imgurl, filepath, refere, &imginfo)
			if result {
				l.Send(slog.LevelInfo, fmt.Sprintf("Image %s DownLoad Success And Save In: %s", imginfo.ImageId, filepath), log.LogFiles|log.LogStdouts)

			} else {
				l.Send(slog.LevelWarn, fmt.Sprintf("Image %s DownLoad Fail", imginfo.ImageId), log.LogFiles|log.LogStdouts)
			}
			wg.Done()
		}(imgUrl, filePath, *imgInfo)
	}
	wg.Wait()
	imgInfo.Status = true

}

func (br *Br) PixivImage(imgid string, doc *goquery.Document) (*sql.ImageInfo, error) {
	if IsNotLoggedIn(doc) {
		panic("Not Logged In!")
	}
	if IsNeedPermission(doc) {
		panic("Not in MyPick List, Need Permission!")
	}
	if IsNeedAppropriateLevel(doc) {
		panic("Public works can not be viewed by the appropriate level!")
	}
	if IsDeleted(doc) {
		panic("Image not found/already deleted!")
	}
	if IsGuroDisabled(doc) {
		panic("Image is disabled for under 18, check your setting page (R-18/R-18G)!")
	}

	// Artist

	return ParseInfo(doc, imgid)
}

func (br *Br) DownLoadImage(imgurl, fileName, referer string, image *sql.ImageInfo) bool {

	for retryCount := 0; retryCount < NewWork.Retry; retryCount++ {
		if br.PerformDownload(imgurl, fileName, referer) {
			return true
		}
		l.Send(slog.LevelWarn, fmt.Sprintf("Retry time %d", retryCount), log.LogStdouts)
		time.Sleep(time.Duration(NewWork.RetryWait) * 100 * time.Millisecond)
	}
	return false
}

func (br *Br) PerformDownload(imgUrl, fileNameSave, referer string) bool {

	resp := GetPixivPage(imgUrl, referer)

	f, err := os.OpenFile(fileNameSave, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		l.Send(slog.LevelError, "图片文件创建失败", log.LogFiles|log.LogStdouts)
		return false
	}
	df := bytes.NewReader(resp)
	//resp.Read(buf.Bytes())
	//buf.Write(resp)
	f.Write(resp)
	defer f.Close()

	st, _ := os.Stat(fileNameSave)
	if st.Size() == int64(df.Len()) {
		return true
	}
	return false

}
