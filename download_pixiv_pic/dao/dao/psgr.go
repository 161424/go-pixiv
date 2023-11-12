package dao

import (
	"bytes"
	"fmt"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"github.com/chen/download_pixiv_pic/log"
	"github.com/chen/download_pixiv_pic/pkg/artist"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

// 状态输出类的包内容

type PsgrDB struct {
	DB *gorm.DB
}

// Client 客户端

const (
	path          = "./database/sql/"
	authsqlname   = "auth.sql"
	cookiesqlname = "cookie.sql"
)

var l *log.Logs

type name interface {
	*sql.Auth | *sql.ImageInfo | *artist.PixivArtist
}

func init() {
	l = log.NewSlogGroup("dao.psgr")
}

// GetClient 获取一个数据库客户端
func GetClient() *PsgrDB {
	return &PsgrDB{}
}

// MySQLConfig 数据库配置
func (Pg *PsgrDB) Open(pg map[string]interface{}) (*PsgrDB, string) {
	buf := bytes.Buffer{}
one:
	for i, j := range pg {
		buf.WriteString(i)
		buf.WriteString("=")

		var v string
		switch j.(type) {
		case int:
			v = strconv.Itoa(j.(int))
		case string:
			v = j.(string)
		default:
			break one
		}

		buf.WriteString(v)
		buf.WriteString(" ")
	}

	var err error
	dns := buf.String()
	Pg.DB, err = gorm.Open(postgres.Open(dns))

	if err != nil {
		return nil, err.Error()
	}
	return Pg, ""

}

func (Pg *PsgrDB) GetDB() *gorm.DB {
	return Pg.DB
}

func Create[T name](Pg *PsgrDB, tb T) error {
	err := Pg.DB.AutoMigrate(&tb)
	if err != nil {
		return err
	}
	return nil
}

func (Pg *PsgrDB) CreateDb() error {

	if err := Create(Pg, &artist.PixivArtist{}); err != nil {
		fmt.Println(err)
		return err
	}

	if err := Create(Pg, &sql.ImageInfo{}); err != nil {
		fmt.Println(err)
		return err
	}

	if err := Create(Pg, &sql.Auth{}); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (Pg *PsgrDB) SelectImageByImageId(imageid string) (error, sql.ImageInfo) {
	var dbImgInfo = &sql.ImageInfo{}
	f := Pg.DB.Where("ImageId = ?", imageid).Model(&sql.ImageInfo{}).First(dbImgInfo)
	if f.Error != nil {
		l.Send(slog.LevelError, fmt.Sprintf("SelectImageByImageIdError:%s", f.Error), log.LogStdouts|log.LogFiles)
		return f.Error, *dbImgInfo
	}
	return nil, *dbImgInfo
}

func (Pg *PsgrDB) SaveImageId(imgInfo *sql.ImageInfo) error {

	f := Pg.DB.Where("ImageId=?", imgInfo.ImageId).Save(imgInfo)
	if f.Error != nil {
		l.Send(slog.LevelError, fmt.Sprintf("UpdateImageIdError:%s", f.Error), log.LogStdouts|log.LogFiles)
		return f.Error
	}
	return nil
}

func (Pg *PsgrDB) SelectArtistImages(art *artist.PixivArtist) error {

	//f := Pg.DB.Where("ArtistId = ?", art.ArtistId).Save(art)
	if f := Pg.DB.Where("ArtistId = ?", art.ArtistId).Save(art); f.Error != nil {
		l.Send(slog.LevelError, fmt.Sprintf("SelectArtistImagesErr=%s", f.Error), log.LogStdouts|log.LogFiles)
	} else {
		l.Send(slog.LevelInfo, fmt.Sprintf("Successfully Update: %s", art.ArtistId), log.LogStdouts|log.LogFiles)

	}
	return nil
}

func (Pg *PsgrDB) UpdateArtist(art *artist.PixivArtist) error {

	//f := Pg.DB.Where("ArtistId = ?", art.ArtistId).Save(art)
	if f := Pg.DB.Where("ArtistId = ?", art.ArtistId).Save(art); f.Error != nil {
		l.Send(slog.LevelError, fmt.Sprintf("SelectArtistImagesErr=%s", f.Error), log.LogStdouts|log.LogFiles)
	} else {
		l.Send(slog.LevelInfo, fmt.Sprintf("Successfully Update: %s", art.ArtistId), log.LogStdouts|log.LogFiles)

	}
	return nil
}

func (Pg *PsgrDB) SaveAuthDownloadRecord(auth *sql.Auth) {
	if f := Pg.DB.Save(auth); f.Error != nil {
		l.Send(slog.LevelError, fmt.Sprintf("SaveAuthDownloadRecordErr=%s", f.Error), log.LogStdouts|log.LogFiles)
	} else {
		l.Send(slog.LevelInfo, fmt.Sprintf("Successfully Save: %s", auth.Uname), log.LogStdouts|log.LogFiles, "err", f.Error)

	}
	return
}
