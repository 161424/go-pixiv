package psgr

import (
	"fmt"
	"github.com/chen/download_pixiv_pic/dao/dao"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"gorm.io/gorm"
)

var Authdb *gorm.DB

func GetAuthTable() *gorm.DB {
	db := dao.GetClient()
	neu := sql.InitAuth()
	err := db.DB.AutoMigrate(&sql.Auth{})
	Authdb := db.DB.Create(neu)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(dv)
	return Authdb
}
