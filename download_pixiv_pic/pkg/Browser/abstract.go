package browser

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chen/download_pixiv_pic/cmd/utils"
	"github.com/chen/download_pixiv_pic/common/addr"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"github.com/chen/download_pixiv_pic/pkg/artist"
	"path/filepath"

	//"github.com/chen/download_pixiv_pic/pkg/Browser"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var NewWork = addr.NewNetWork()

func IsNotLoggedIn(page *goquery.Document) bool {
	check := page.Find(".signup_button")
	if len(check.Text()) > 0 {
		return true
	}
	check = page.Find(".ui-button _signup")
	if len(check.Text()) > 0 {
		return true
	}
	return false
}

func IsNeedPermission(page *goquery.Document) bool {
	errorMessages := []string{"この作品は.+さんのマイピクにのみ公開されています|この作品は、.+さんのマイピクにのみ公開されています",
		"This work is viewable only for users who are in .+\\'s My pixiv list",
		"Only .+\\'s My pixiv list can view this.",
		"<section class=\"restricted-content\">"}
	return HaveString(page, errorMessages)
}

func HaveString(page *goquery.Document, s []string) bool {
	for _, str := range s {
		reg := regexp.MustCompile(str)
		result := reg.FindAllStringSubmatch(page.Text(), -1)
		for _, r := range result {
			if len(r[1]) > 0 {
				return true
			}
		}
	}
	return false
}

func IsNeedAppropriateLevel(page *goquery.Document) bool {
	errorMessages := []string{"該当作品の公開レベルにより閲覧できません。"}
	return HaveString(page, errorMessages)
}

func IsDeleted(page *goquery.Document) bool {
	errorMessages := []string{"該当イラストは削除されたか、存在しないイラストIDです。|該当作品は削除されたか、存在しない作品IDです。",
		"この作品は削除されました。",
		"The following work is either deleted, or the ID does not exist.",
		"This work was deleted.",
		"Work has been deleted or the ID does not exist."}
	return HaveString(page, errorMessages)
}

func IsGuroDisabled(page *goquery.Document) bool {
	errorMessages := []string{"表示されるページには、18歳未満の方には不適切な表現内容が含まれています。",
		"The page you are trying to access contains content that may be unsuitable for minors"}
	return HaveString(page, errorMessages)
}

func IsErrorExist(page *goquery.Document) {
	//check := page.Find("span.error")
}

func PrintInfo(img *sql.ImageInfo) {
	fmt.Printf("  User Name: %s\n", img.UserName)
	fmt.Printf("  Image Title: %s\n", img.ImageTitle)
	fmt.Printf("  Tags Title: %s\n", img.ImageTags)
	fmt.Printf("  Translated Tags: %s\n", img.TranslationTag)
	fmt.Printf("  Date: %s\n", img.WorksDateDateTime)
	fmt.Printf("  Mode: %s\n", img.ImageMode)
	if img.ImageMode == "manga" {
		fmt.Printf("  Pages: %d\n", img.ImageCount)
	}
	fmt.Printf("  Bookmarks %d\n", img.BookmarkCount)

}

func MakeFilename(img *sql.ImageInfo, path string, user, mode, data, content string) (string, error) {
	fn := utils.DelSpeChar(img.UserName)
	var r string
	r = fmt.Sprintf("%s(%s)", fn, img.UserID)
	if user == "rank" {
		r = filepath.Join(user, mode, fmt.Sprintf("%s-%s", data, content), r)
	} else if user == "author" {
		// 作者的收藏，作品等等
		r = filepath.Join(user, fmt.Sprintf("%s-%s", data, content), r)
	} else if user == "tag" {
		r = filepath.Join(user, fmt.Sprintf("%s-%s", data, content), r)
	}

	path = filepath.Join(path, r)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModeDir|os.ModePerm)
			if err != nil {
				return "", err
			}
		}

	}
	fmt.Printf("  Creating directory: %s. Successful!\n", path)

	return path, nil
}

func ParseInfo(doc *goquery.Document, imgId string) (*sql.ImageInfo, error) {
	imageInfo := &sql.ImageInfo{
		PixivArtist: &artist.PixivArtist{},
	}

	r, _ := doc.Find("meta#meta-preload-data").Attr("content")
	var con = &Cont{}
	json.Unmarshal([]byte(r), &con)
	if len(con.Illust[imgId]) == 0 {
		return nil, nil
	}
	// Artist
	parent := false
	if !parent {
		for artistId, value := range con.User {
			imageInfo.GetPA(artistId, value, true, 0, 0)
		}
	}

	root := con.Illust[imgId]

	imageCount := int((root["pageCount"]).(float64))
	tempUrl := ((root["urls"]).(map[string]interface{})["original"]).(string)
	tempResizedUrl := ((root["urls"]).(map[string]interface{})["regular"]).(string)
	imageInfo.ImageCount = imageCount

	//  不够齐全
	if imageCount == 1 {
		if strings.Contains(tempUrl, "ugoira") { // 动图
			imageInfo.ImageMode = "ugoira_view"
			tempUrl = strings.Replace(tempUrl, "/img-original/", "/img-zip-ugoira/", -1)
			tempUrl = strings.Split(tempUrl, "_ugoira0")[0]
			tempUrl = tempUrl + "_ugoira1920x1080.zip"
			imageInfo.ImageUrls = append(imageInfo.ImageUrls, tempUrl)

			tempResizedUrl = strings.Replace(tempResizedUrl, "/img-original/", "/img-zip-ugoira/", -1)
			tempResizedUrl = strings.Split(tempResizedUrl, "_ugoira0")[0]
			tempResizedUrl = tempResizedUrl + "_ugoira600x600.zip"
			imageInfo.ImageResizedUrls = append(imageInfo.ImageResizedUrls, tempResizedUrl)
		} else {
			imageInfo.ImageMode = "big"
			imageInfo.ImageUrls = append(imageInfo.ImageUrls, tempUrl)
			imageInfo.ImageResizedUrls = append(imageInfo.ImageResizedUrls, tempResizedUrl)
		}
	} else {
		imageInfo.ImageMode = "manga"
		for i := 0; i < imageCount; i++ {
			_tempUrl := strings.Replace(tempUrl, "_p0", "_p"+strconv.Itoa(i), -1)
			imageInfo.ImageUrls = append(imageInfo.ImageUrls, _tempUrl)
			_tempResizedUrl := strings.Replace(tempResizedUrl, "_p0", "_p"+strconv.Itoa(i), -1)
			imageInfo.ImageResizedUrls = append(imageInfo.ImageResizedUrls, _tempResizedUrl)
		}

	}
	//fmt.Println(imageInfo.ImageUrls)
	imageInfo.ImageTitle = (root["illustTitle"]).(string)
	//imageInfo.ImageCaption = (root["illustComment"]).(string)
	if root["seriesNavData"] == nil { // 是否是系列
		imageInfo.SeriesNavData = map[string]string{"nil": "nil"}
	} else {
		mss := make(map[string]string, 6)
		mss["seriesType"] = (root["seriesNavData"].(map[string]interface{})["seriesType"]).(string)
		mss["seriesId"] = (root["seriesNavData"].(map[string]interface{})["seriesId"]).(string)
		mss["title"] = (root["seriesNavData"].(map[string]interface{})["title"]).(string)
		mss["order"] = strconv.Itoa((root["seriesNavData"].(map[string]interface{})["order"]).(int))
		if root["seriesNavData"].(map[string]interface{})["prev"] == nil {
			mss["prev"] = ""
		} else {
			mss["prev"] = ((root["seriesNavData"].(map[string]map[string]interface{})["prev"])["id"]).(string)
		}

		if root["seriesNavData"].(map[string]interface{})["prev"] == nil {
			mss["next"] = ""
		} else {
			mss["next"] = ((root["seriesNavData"].(map[string]map[string]interface{})["prev"])["id"]).(string)
		}

		imageInfo.SeriesNavData = mss
	}

	imageInfo.JdRtv = (root["viewCount"]).(float64) // 查看人数
	imageInfo.JdRtc = (root["likeCount"]).(float64) // 喜欢人数
	tags := root["tags"]
	if tg, ok := tags.(map[string]interface{}); ok == true && tg != nil {

		//if tg == nil {
		//	panic("tg err")
		//}

		gt := tg["tags"]
		p := (gt).([]interface{})
		for _, tag := range p {
			tagp := tag.(map[string]interface{})
			imageInfo.ImageTags = append(imageInfo.ImageTags, (tagp["tag"]).(string))
			if v, isok := tagp["translation"]; isok == true {
				if v != nil {
					for _, valu := range v.(map[string]interface{}) {
						imageInfo.TranslationTag = append(imageInfo.TranslationTag, (valu).(string))
					}
				}
			}
		}
	}

	//fmt.Println("createDate", root["createDate"])
	imageInfo.WorksDateDateTime, _ = time.Parse(time.RFC3339, (root["createDate"]).(string))
	imageInfo.WorksResolution = fmt.Sprintf("%.0fx%.0f", (root["width"]).(float64), (root["height"]).(float64))
	//if imageCount > 1 {
	//	imageInfo.WorksResolution = fmt.Sprintf("Multiple images: %dP", imageCount)
	//}

	imageInfo.BookmarkCount = int64((root["bookmarkCount"]).(float64))
	imageInfo.ImageResponseCount = (root["responseCount"]).(float64)

	if root["aiType"] != nil {
		imageInfo.AiType = int((root["aiType"]).(float64)) //1 == non-AI, 2 == AI-generated
		if imageInfo.AiType == 2 {
			imageInfo.ImageTags = append(imageInfo.ImageTags, "AI-generated")
		}
	}
	imageInfo.UserID = (root["userId"]).(string)
	imageInfo.UserName = (root["userName"]).(string)
	//imageInfo.UserAccount = (root["userAccount"]).(string)

	return imageInfo, nil

}

func CheckImageIsExiste() bool {
	return false
}
