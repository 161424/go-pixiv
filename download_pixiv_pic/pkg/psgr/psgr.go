package psgr

import (
	"fmt"
	sql_ "github.com/chen/download_pixiv_pic/common/model/sql"
	"github.com/chen/download_pixiv_pic/database/sql"
	"gorm.io/gorm"
)

var Authdb *gorm.DB

func GetAuthTable() *gorm.DB {
	db := sql_.GetClient()
	neu := sql.InitAuth()
	err := db.DB.AutoMigrate(&sql.Auth{})
	Authdb := db.DB.Create(neu)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(dv)
	return Authdb
}
