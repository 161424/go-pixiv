package Artist

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

//type DA struct {
//	*sql.PixivArtist
//}

type PixivArtist struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	ArtistId         string         `gorm:"column:ArtistId;type:text;primaryKey"`
	ArtistName       string         `gorm:"column:ArtistName;type:text"`
	ArtistAvatar     string         `gorm:"column:ArtistAvatar;type:text"`
	ArtistBackground string         `gorm:"column:ArtistBackground;type:text"`
	TotalImages      int            `gorm:"column:TotalImages"`
	LocalImages      int            `gorm:"column:LocalImages"`
	ImageList        pq.StringArray `gorm:"column:ImageList;type:text[]"`
	LocalImagesList  pq.StringArray `gorm:"column:LocalImagesList;type:text[]"`
	IsLastPage       bool           `gorm:"column:IsLastPage"`
	HaveImages       bool           `gorm:"column:HaveImages"`
	Offset           int            `gorm:"column:Offset"`
	Limit            int            `gorm:"column:Limit"`
	//ReferenceImageId string         `gorm:"column:ReferenceImageId;type:text"`
	MangaSeries pq.StringArray `gorm:"column:MangaSeries;type:text[]"`
	NovelSeries pq.StringArray `gorm:"column:NovelSeries;type:text[]"`
}

//var DA *PixivArtist

func (da *PixivArtist) GetPA(mid string, r map[string]interface{}, fromImage bool, offset, limit int) {

	da.Offset = offset
	da.Limit = limit
	da.ArtistId = mid

	da.ArtistAvatar = "no_profile"
	da.ArtistName = "self"
	da.ArtistBackground = "no_background"

	if !fromImage {
		if (r["error"]).(bool) {
			log.Panicln("err for get page")
		}
		if r["body"] == nil {
			log.Panicln("Missing body content, possible artist id doesn't exists.")
		}
		body := (r["body"]).(map[string]interface{})
		da.ParseImages(body)
		da.ParseMangaList(body)
		da.ParseNovelList(body)
		//da.ArtistSomeInfo(mid)

	} else {
		da.IsLastPage = true
		da.HaveImages = true
		da.ArtistId = (r["userId"]).(string)
		da.ArtistAvatar = strings.Replace((r["image"]).(string), "_50.", ".", -1)
		da.ArtistAvatar = strings.Replace((r["image"]).(string), "_170.", ".", -1)
		da.ArtistName = (r["name"]).(string)
		if r["background"] != nil {
			bk := (r["background"]).(map[string]interface{})
			da.ArtistBackground = (bk["url"]).(string)
		}

	}

	//da.ParseInfo(r, fromImage, false)
	//return da
}

func (da *PixivArtist) ParseImages(body map[string]interface{}) {
	if body["works"] != nil {
		//fmt.Println(body["works"])
		work := (body["works"]).([]interface{})

		for _, img := range work {
			id := (img.(map[string]interface{}))["id"].(string)
			da.ImageList = append(da.ImageList, id)
		}
		da.TotalImages = int((body["total"]).(float64))
		if len(da.ImageList) > 0 {
			da.HaveImages = true
		}
		if len(da.ImageList)+da.Offset == da.TotalImages {
			da.IsLastPage = true
		} else {
			da.IsLastPage = false
		}
	} else {
		if body["illusts"] != nil {
			illusts := (body["illusts"]).(map[string]interface{})
			for img, _ := range illusts {
				da.ImageList = append(da.ImageList, img)
			}
			da.ArtistId = (illusts["illust_user_id"]).(string)
			da.ArtistName = (illusts["user_name"]).(string)
		}

		if body["manga"] != nil {
			illusts := (body["manga"]).(map[string]interface{})
			for _, img := range illusts {
				da.ImageList = append(da.ImageList, (img).(string))
			}
		}
		if body["illust"] != nil {
			illust := (body["illust"]).(map[string]interface{})
			da.ArtistId = (illust["illust_user_id"]).(string)
			da.ArtistName = (illust["user_name"]).(string)
		} else if body["novel"] != nil {
			novel := (body["novel"]).(map[string]interface{})
			da.ArtistId = (novel["user_id"]).(string)
			da.ArtistName = (novel["user_name"]).(string)
		}

		da.TotalImages = len(da.ImageList)
		if da.Offset+da.Limit >= da.TotalImages {
			da.IsLastPage = true
		} else {
			da.IsLastPage = false
		}
		if len(da.ImageList) > 0 {
			da.HaveImages = true
		}
	}
}

func (da *PixivArtist) ParseMangaList(body map[string]interface{}) {
	if len(body) != 0 && body["mangaSeries"] != nil {
		ms := (body["mangaSeries"]).([]map[string]interface{})
		for _, i := range ms {
			da.MangaSeries = append(da.MangaSeries, (i["id"]).(string))
		}
	}
}

func (da *PixivArtist) ParseNovelList(body map[string]interface{}) {
	if len(body) != 0 && body["novelSeries"] != nil {
		ms := (body["novelSeries"]).([]map[string]interface{})
		for _, i := range ms {
			da.MangaSeries = append(da.MangaSeries, (i["id"]).(string))
		}
	}
}

func (da *PixivArtist) GetMemberInfoWhitecube(memid, bookmark string) {
	//
}
