package sql

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/chen/download_pixiv_pic/pkg/Artist"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log"
	"time"
)

type ImageInfo struct {
	gorm.Model
	ImageId string `gorm:"column:imageid;NOT NULL;primaryKey"`
	Path    string `gorm:"column:Path;"`

	ImageMode        string         `gorm:"column:ImageMode"`
	ImageUrls        pq.StringArray `gorm:"column:ImageUrls;type:text[]"`
	ImageResizedUrls pq.StringArray `gorm:"column:ImageResizedUrls;type:text[]"`
	ImageTitle       string         `gorm:"column:ImageTitle"`
	//ImageCaption     string                 `gorm:"column:ImageCaption"`
	ImageCount    int     `gorm:"column:ImageCount"`
	SeriesNavData MapSS   `gorm:"column:SeriesNavData;type:text"`
	JdRtv         float64 `gorm:"column:Jd_rtv"`
	JdRtc         float64 `gorm:"column:Jd_rtc"`
	//Jd_rtt           string `gorm:"column:ImageMode"`
	ImageTags          pq.StringArray `gorm:"column:ImageTags;type:text[]"`
	TranslationTag     pq.StringArray `gorm:"column:TranslationTag;type:text[] "`
	WorksDateDateTime  time.Time      `gorm:"column:WorksDateDateTime"`
	WorksResolution    string         `gorm:"column:WorksResolution"`
	BookmarkCount      int64          `gorm:"column:BookmarkCount"`
	ImageResponseCount float64        `gorm:"column:ImageResponseCount"`
	//ParseUrlFromCaption
	AiType      int    `gorm:"column:AiType"`
	UserID      string `gorm:"column:UserID"`
	UserName    string `gorm:"column:UserName"`
	UserAccount string `gorm:"column:UserAccount"`
	*Artist.PixivArtist
}

func (c *ImageInfo) Create(db *gorm.DB) error {
	a := pq.StringArray{}
	b := []pq.StringArray{c.ImageUrls, c.ImageResizedUrls, c.ImageTags, c.TranslationTag}
	for _, _b := range b {
		if _b != nil && len(_b) > 0 {
			for _, v := range _b {
				a = append(a, v)
			}
			_b = a
		}
	}

	if result := db.Create(c); result.Error != nil {
		log.Printf("Error creating company: %s", c.ImageId)
		return result.Error
	} else {
		log.Printf("Successfully created company: %s", c.ImageId)
		return nil
	}
}

//func (c *ImageInfo) Find() error {
//	db := sql2.GetClient()
//	test := []ImageInfo{}
//	db.DB.Find(&test)
//	for i, j := range test {
//		log.Printf("%T,%T, %d, %s", i, j, i, j.Local)
//	}
//	return nil
//}

type MapSS map[string]string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *MapSS) Scan(value interface{}) (err error) {

	var skills map[string]string
	switch value.(type) {
	case string:
		err = json.Unmarshal([]byte(value.(string)), &skills)
	case []byte:
		err = json.Unmarshal(value.([]byte), &skills)
	default:
		return errors.New("Incompatible type for Skills")
	}
	if err != nil {
		return err
	}
	*j = skills
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j MapSS) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	if n := len(j); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+3*n)
		b[0] = '{'
		var _n = 0
		for key, value := range j {
			_n += 1
			b = appendArrayQuotedBytes(b, []byte(key))
			b = append(b, ':')
			b = appendArrayQuotedBytes(b, []byte(value))
			if _n == n {
				continue
			}
			b = append(b, ',')
		}
		return string(append(b, '}')), nil
	}

	return "{}", nil
}

func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '"')
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '"')
}
