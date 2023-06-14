package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chen/download_pixiv_pic/pkg/Artist"
	"github.com/chen/download_pixiv_pic/pkg/Browser"
	"log"
	"strings"
)

func DownLoadByAuth(string2 string) {

	memberId, bookmark, tags, profile := GetInputInfo()
	sp, ep, num := GetStartAndEndNumber()
	ProcessMember(memberId, Rootpath, bookmark, tags, profile, sp, ep, num)
}

func ProcessMember(memberId, Rootpath, bookmark, tags, profile string, sp, ep, num int) {
	da := &Artist.PixivArtist{}
	br := &Browser.Br{}
	fmt.Printf("Processing Member Id: %s\n", memberId)
	var url = ""
	var offsetStop int
	limit := 48
	offset := (sp - 1) * limit
	if num > 0 {
		offsetStop = num
	} else {
		offsetStop = ep * limit
	}

	if bookmark != "" {
		// 表示 收藏
		// https://www.pixiv.net/ajax/user/6558698/illusts/bookmarks?tag=FGO&offset=0&limit=24&rest=show
		url = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illusts/bookmarks?tag=%s&offset=%d&limit=%d&rest=show", memberId, tags, offset, limit)
	} else {
		if len(tags) > 0 {
			url = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=%s&offset=%d&limit=%d", memberId, tags, offset, limit)
			//} else if r18mode != "" {
			//	url = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=R-18&offset=%s&limit=%s", memberId, offset, limit)
		} else if profile != "" {
			// 标识某人的精选
			url = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/profile/all", memberId)
		} else {
			url = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=&offset=%d&limit=%d", memberId, offset, limit)
		}
	}
	fmt.Printf("Member Url %s\n", url)
	resp := Browser.GetPixivPage(url, "")
	var r map[string]interface{}
	json.Unmarshal(resp, &r)
	da.GetPA(memberId, r, false, offset, limit)
	ArtistSomeInfo(da)
	if len(da.ImageList) == 0 {
		fmt.Printf("No images for Member Id: %s, from Bookmark: %s", memberId, bookmark)
	}
	//da.ReferenceImageId = da.ImageList[0]
	da.GetMemberInfoWhitecube(memberId, bookmark)

	fmt.Printf("Member Name %s\n", da.ArtistName)
	fmt.Printf("Member Avatar %s\n", da.ArtistAvatar)
	fmt.Printf("Member Backgrd %s\n", da.ArtistBackground)
	var printOffsetStop int
	if offsetStop < da.TotalImages {
		printOffsetStop = offsetStop
	} else {
		printOffsetStop = da.TotalImages
	}
	fmt.Printf("Processing images from %d to %d of %d\n", offset+1, printOffsetStop, da.TotalImages)

	if !da.HaveImages {
		log.Printf("No image found for: %d\n", memberId)
	}

	var imgid string
	var dnlist []string
	var id int
	var ns = 0
	for {
		rootpath := GetRootPath(Rootpath, da.ArtistName, bookmark, tags, profile)
		fmt.Println(rootpath)
		for id, imgid = range da.ImageList {
			up := fmt.Sprintf("[ %d of %d ]", id+1, printOffsetStop)
			if da.TotalImages > 0 {
				printOffsetStop -= offset
			} else {
				printOffsetStop = (sp-1)*20 + len(da.ImageList)
			}
			for retryCount := 0; retryCount < Browser.NewWork.Retry; retryCount++ {
				fmt.Printf("MemberId: %s Page: %d Post %d of %s\n", memberId, sp, up, da.TotalImages)

				result := br.ProcessImage(imgid, rootpath, "mb")
				if result == "YES" {
					ns += 1
					dnlist = append(dnlist, imgid)
					break
				}
			}
		}
		sp += 1
		if sp > ep {
			break
		}
	}

	// 更新artist数据库数据
	DB.UpdateArtist(da)
	da.LocalImagesList = dnlist
	da.LocalImages = ns
	fmt.Printf("last image_id: %s\n", imgid)
	fmt.Printf("Member_id: %s  completed: %s\n")
}

func GetRootPath(path, an, bookmark, tags, profile string) string {
	path = path + Spt + "Artist" + Spt + an
	// 表示收藏的图片
	if bookmark != "" {
		_path := path + Spt + bookmark
		Browser.CheckPathIsExit(_path)
		return _path
	} else {
		// 表示 标签的图片
		if tags != "" {
			if tags == "R18" {
				_path := path + Spt + tags
				Browser.CheckPathIsExit(_path)
				return _path
			}
			_path := path + Spt + tags
			Browser.CheckPathIsExit(_path)
			return _path

		} else if profile != "" {
			_path := path + Spt + profile
			Browser.CheckPathIsExit(_path)
			return _path
		} else {
			_path := path + Spt + "all"
			Browser.CheckPathIsExit(_path)
			return _path
		}
	}

	return path + Spt + "other"

}

func GetInputInfo() (md, bm, ts, pf string) {
mb:
	fmt.Println("选择输入选项代号\n 1.用户名\n 2.用户id")
	var mb int
	//var md string
	fmt.Scanln(&mb)
	if mb == 1 {
		fmt.Print("请输入用户名: ")
		fmt.Scanln(&md)
		url := fmt.Sprintf("https://www.pixiv.net/search_user.php?s_mode=s_usr&i=0&nick=%s", md)
		rb := Browser.GetPixivPage(url, "")
		rep := bytes.NewReader(rb)
		doc, _ := goquery.NewDocumentFromReader(rep)
		n, o := doc.Find(".user-recommendation-item > a").Attr("href")
		if o {
			md = strings.Split(n, "/")[2]
		} else {
			fmt.Printf("未找到输入的 %s", md)
			goto mb
		}
	} else if mb == 2 {
		fmt.Print("请输入用户id: ")
		fmt.Scanln(&md)

	} else {
		goto mb
	}

	var t_ string
	//var bm, ts, pf string
info:
	fmt.Print("请输入要下载用户的 1.插画(illus), 2.收藏(bookmark): ")
	fmt.Scanln(&t_)
	if t_ == "1" {
		bm = ""
		fmt.Print("是否要下载用户插画的 profile(y/n): ")
		fmt.Scanln(&t_)

		if t_ == "y" {
			pf = "true"
		} else {
			pf = ""
		}
	} else if t_ == "2" {
		bm = "true"
	} else {
		goto info
	}

	fmt.Print("请输入要下载用户插画的 tags: ")
	fmt.Scanln(&t_)

	if len(t_) == 0 {
		ts = ""
		fmt.Println("tags is all")
	} else if strings.Contains(t_, "R18") {
		ts = "R-18"
		fmt.Println("tags is R-18")
	} else {
		ts = t_
		//fmt.Printf("tags is %s\n", ts)
	}

	fmt.Printf("medid: %s; bookmark: %s; tags: %s; pf: %s.\n", md, bm, ts, pf)
	return md, bm, ts, pf

}

func ArtistSomeInfo(da *Artist.PixivArtist) {
	// https://www.pixiv.net/ajax/user/83739
	// https://www.pixiv.net/ajax/user/6558698
	url := fmt.Sprintf("https://www.pixiv.net/ajax/user/%s", da.ArtistId)
	var r map[string]interface{}
	rb := Browser.GetPixivPage(url, "")
	json.Unmarshal(rb, &r)

	da.ArtistName = ((r["body"]).(map[string]interface{})["name"]).(string)
	da.ArtistAvatar = strings.Replace(((r["body"]).(map[string]interface{})["image"]).(string), "_50.", ".", -1)
	if (r["body"]).(map[string]interface{})["background"] != nil {
		bk := (r["body"]).(map[string]interface{})["background"].(map[string]interface{})
		da.ArtistBackground = (bk["url"]).(string)
	}
	//da.ArtistBackground = ((r["body"]).(map[string]interface{})["background"]).(string)
}
