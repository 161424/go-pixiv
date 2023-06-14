package sql

import (
	"bytes"
	conf2 "github.com/chen/download_pixiv_pic/common/conf"
	"github.com/chen/download_pixiv_pic/database/sql"
	"github.com/chen/download_pixiv_pic/pkg/Artist"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
)

// 状态输出类的包内容

type PsgrDB struct {
	DB *gorm.DB
}

var Ps = &PsgrDB{}

// Client 客户端

func GetLn() *PsgrDB {
	return Ps
}

func NewClient(db *gorm.DB) *PsgrDB {
	return &PsgrDB{
		DB: db,
	}

}

// GetClient 获取一个数据库客户端
func GetClient() *PsgrDB {

	//fmt.Println(conf2.ConfigData["postgres"])
	Ps.Open(conf2.ConfigData["postgres"])
	return Ps
}

// MySQLConfig 数据库配置
func (Pg *PsgrDB) Open(pg map[string]interface{}) {
	buf := bytes.Buffer{}
	for i, j := range pg {
		buf.WriteString(i)
		buf.WriteString("=")

		var v string
		if _, ok := j.(string); ok != true {
			v = strconv.Itoa(j.(int))
		} else {
			v = j.(string)
		}
		buf.WriteString(v)
		buf.WriteString(" ")
	}

	var err error
	dns := buf.String()
	Pg.DB, err = gorm.Open(postgres.Open(dns))
	//fmt.Println(Pg)

	if err != nil {
		log.Panic(err)
	}
	//fmt.Printf("PostGre connect success. %s", dns)

}

func (Pg *PsgrDB) GetDB() *gorm.DB {
	return Pg.DB
}

func (Pg *PsgrDB) SelectImageByImageId(imageid string) int64 {
	//Pg.DB.AutoMigrate(&sql.ImageInfo{})
	Create(Pg, &sql.ImageInfo{})
	f := Pg.DB.Where("ImageId = ?", imageid).Find(&sql.ImageInfo{})
	if f.Error != nil {
		log.Panic("寻找image_id error")
	}
	return f.RowsAffected
}

func (Pg *PsgrDB) UpdateArtist(art *Artist.PixivArtist) {
	Create(Pg, &Artist.PixivArtist{})
	//f := Pg.DB.Where("ArtistId = ?", art.ArtistId).Save(art)
	if result := Pg.DB.Where("ArtistId = ?", art.ArtistId).Save(art); result.Error != nil {
		log.Printf("Error Update: %s", art.ArtistId)
	} else {
		log.Printf("Successfully Update: %s", art.ArtistId)
	}
}

type name interface {
	*sql.Auth | *sql.ImageInfo | *Artist.PixivArtist
}

func Create[T name](Pg *PsgrDB, tb T) error {
	err := Pg.DB.AutoMigrate(&tb)
	if err != nil {
		log.Panic("create table err")
	}
	return nil
}

func CreateDb(Pg *PsgrDB) {
	Create(Pg, &Artist.PixivArtist{})
	Create(Pg, &sql.ImageInfo{})
}
