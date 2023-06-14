package Browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	sql2 "github.com/chen/download_pixiv_pic/common/model/sql"
	"github.com/chen/download_pixiv_pic/database/sql"
	"log"
	"os"
	"strings"
	"time"
)

type Br struct {
}

var ImageCount = 0
var RJ *RpJs

type RpJs struct {
	Mode      interface{}              `json:"mode"`
	CurrDate  interface{}              `json:"date"`
	NextDate  interface{}              `json:"next_date"`
	PrevDate  interface{}              `json:"prev_date"`
	CurrPage  interface{}              `json:"page"`
	NextPage  interface{}              `json:"next"`
	PrevPage  interface{}              `json:"prev"`
	RankTotal interface{}              `json:"rank_total"`
	Contents  []map[string]interface{} `json:"contents"`
}

type Cont struct {
	Timestamp string
	Illust    map[string]map[string]interface{}
	User      map[string]map[string]interface{}
}

func (br *Br) GetPixivPage(url, referer string, enable_cache bool) (*RpJs, string) {
	if referer == "" {
		referer = "https://www.pixiv.net"
	}
	url = br.fixUrl(url, true)
	rb := GetPixivPage(url, referer)

	if strings.Contains(string(rb), "Just a moment...") {
		return nil, "受到了CloudFlare 5s盾的管控， 可能是Cookies已经失效"
	}
	json.Unmarshal(rb, &RJ)
	return RJ, ""
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

func (br *Br) GetPixivRanking(mode, date, content string, filter string, page int) *RpJs {
	url := fmt.Sprintf("https://www.pixiv.net/ranking.php?mode=%s", mode)
	if len(date) > 0 {
		url = fmt.Sprintf("%s&date=%s", url, date)
	}
	if len(content) > 0 {
		url = fmt.Sprintf("%s&content=%s", url, content)
	}
	url = fmt.Sprintf("%s&p=%d&format=json", url, page)
	fmt.Println(url)
	df, err := br.GetPixivPage(url, "", false)
	if err != "" {
		return nil
	}
	return df
}

func (br *Br) ProcessImage(imageId, rootpath, tp string) string {
	//todo 调用数据库
	var inDb = false
	var num int64
	//var fileName string
	downloadImageFlag := true
	//imageId := fmt.Sprintf(strconv.FormatFloat(iid, 'f', 0, 64))

	db := sql2.GetLn()

	referer := "https://www.pixiv.net/artworks/" + imageId
	filename := "no-filename-" + strings.Split(imageId, "\n")[0] + ".tmp"

	fmt.Printf("referer: %s\n", referer)
	fmt.Printf("filename: %s\n", filename)
	fmt.Printf("Processing Image Id: %s\n", imageId)
	imageInfo, _ := br.GetImagePage(imageId)
	imageInfo.ImageId = imageId

	// 创建DB

	if tp == "rk" {
		num = db.SelectImageByImageId(imageId)
	} else if tp == "mb" {

	}
	sql2.CreateDb(db)

	fileName := MakeFilename(imageInfo, rootpath, tp)
	fileName += string(os.PathSeparator)
	imageInfo.Path = filename
	total := len(imageInfo.ImageUrls)
	ImageCount = total

	if num != 0 {
		dirPath, err := os.ReadDir(filename)
		if err != nil {
			log.Println(err)
		} else {
			if len(dirPath) == total {
				inDb = true
			}
		}
	}

	if inDb == true {
		fmt.Printf("Already downloaded in DB: %s\n", imageId)
		downloadImageFlag = false
	}

	if downloadImageFlag {
		PrintInfo(imageInfo)
		for currentImg, imgurl := range imageInfo.ImageUrls {
			fmt.Printf("[%d/%d]Image URL : %s\n", currentImg+1, total, imgurl)
			rp := strings.Split(imgurl, ".")
			_rp := rp[len(rp)-1]
			It := DelSpeChar(imageInfo.ImageTitle)
			_p := fmt.Sprintf("%s_p%d - %s.%s", imageId, currentImg, It, _rp)
			_fileName := fileName + _p
			result, _ := br.DownLoadImage(imgurl, _fileName, referer, imageInfo)
			fmt.Printf(result, imageId)
		}
		db.DB.Create(imageInfo)
	}

	return "YES"

}

func (br *Br) GetImagePage(imageid string) (*sql.ImageInfo, *bytes.Reader) {
	fmt.Printf("Getting image page: %s\n", imageid)
	url := fmt.Sprintf("https://www.pixiv.net%s/artworks/%s/", "", imageid)
	resp := GetPixivPage(url, "")
	rep := bytes.NewReader(resp)
	doc, err := goquery.NewDocumentFromReader(rep)
	if err != nil {
		panic(err)
	}
	imgInfo := br.PixivImage(imageid, doc)
	if imgInfo.ImageMode == "ugoira_view" {
		ugoiraMetaUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/ugoira_meta", imageid)
		// ugoira 代表是动图，暂不进行过多考虑
		fmt.Println(ugoiraMetaUrl)
	}
	return imgInfo, rep
}

func (br *Br) PixivImage(imgid string, doc *goquery.Document) *sql.ImageInfo {
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

func (br *Br) DownLoadImage(imgurl, fileName, referer string, image *sql.ImageInfo) (string, bool) {
	fileNameSave := fileName
	//max_retry := conf.ConfigData["Network"]["Retry"]

	for retryCount := 0; retryCount < NewWork.Retry; retryCount++ {
		if CheckImageIsExiste() {
			return "文件%s已被下载过，现跳过\n", true
		}
		b := br.PerformDownload(imgurl, fileNameSave, referer)
		if b {
			return "文件%s下载成功\n", true
		}
		fmt.Printf("Retry, 第%d次\n", retryCount)
		time.Sleep(time.Duration(NewWork.RetryWait) * time.Second)

	}
	return "重试过后文件%s仍下载失败，现跳过\n", true
}

func (br *Br) PerformDownload(imgurl, fileNameSave, referer string) bool {
	fmt.Printf("PerformDownload. fileNameSave: %s; referer: %s\n", fileNameSave, referer)
	resp := GetPixivPage(imgurl, referer)
	buf := new(bytes.Buffer)
	f, err := os.OpenFile(fileNameSave, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	buf.Write(resp)
	f.Write(resp)
	defer f.Close()

	st, _ := os.Stat(fileNameSave)
	if st.Size() == int64(len(resp)) {
		return true
	}

	return false

}
