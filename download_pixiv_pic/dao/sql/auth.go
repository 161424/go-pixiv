package sql

import (
	"encoding/base64"
	"fmt"
	"github.com/chen/download_pixiv_pic/cmd/utils"
	"github.com/chen/download_pixiv_pic/common/addr"
	"github.com/chen/download_pixiv_pic/conf"
	"gorm.io/gorm"
)

const (
	Success = "success"
	Fail    = "fail"
)

type Auth struct {
	gorm.Model
	Uname    string `gorm:"column:uname;NOT NULL;type:varchar(32);comment:登录账户"json:"uname"`
	Password string `gorm:"column:password;NOT NULL;type:varchar(32);comment:登录密码"json:"password"`
	Cookies  string `gorm:"column:cookies;comment:cookie"json:"cookies"`
	InterIP  string `gorm:"column:InterIP;default:null;type:cidr;comment:InterIP"json:"ip_1"`
	EnterIP  string `gorm:"column:EnterIP;default:null;type:inet;comment:EnterIP"`
	Mac      string `gorm:"column:macaddr;default:null;type:macaddr;comment:macaddr"`
	Output   string `gorm:"column:output;default:null"` // 执行结果
	//RunTimer time.Time `gorm:"column:run_timer;default:null"` // 执行时间
	CostTime float64 `gorm:"column:cost_time"` // 执行耗时
	Status   string  `gorm:"column:status;NOT NULL"`
	//Issue string 	`gorm:"column:issue"`
}

var l = conf.Conf
var A *Auth

func DefaultAuth() (h *Auth) {
	//fmt.Println((conf2.ConfigData["Authentication"]["username"]).(string))
	A = &Auth{
		Uname:    l.GetString("Authentication.username"),
		Password: l.GetString("Authentication.password"),
		Cookies:  l.GetString("Authentication.cookie"),
		InterIP:  addr.GetInterIP(),
		EnterIP:  addr.GetExternalIP(),
		Mac:      addr.GetMacAddr(),
		Output:   "",
		//RunTimer: time.Now(),
		CostTime: 0.0,
		Status:   "",
	}
	return A

}

func DeleAuthAll_id(db *gorm.DB, id int) {
	db.Delete(&Auth{}, id)
}

func DeleAuthAll_hard(db *gorm.DB) {
	db.Unscoped().Delete(&Auth{})
	//db.Exec("DELETE FROM auths")
}

func (auth *Auth) BeforeSave(db *gorm.DB) error {
	if auth.Password != "" {
		keyName := fmt.Sprintf("username_%v", auth.Uname)

		rs, err := utils.AesEcrypt([]byte(auth.Password), []byte(utils.MD5(keyName)))
		if err != nil {
			return err
		}
		auth.Password = base64.StdEncoding.EncodeToString(rs)
	}
	return nil
}

// 返回前解密数据
func (auth *Auth) AfterFind(tx *gorm.DB) error {
	if auth.Password != "" {
		b, err := base64.StdEncoding.DecodeString(auth.Password)
		if err == nil {
			keyName := fmt.Sprintf("username_%v", auth.Uname)

			rs, err := utils.AesDecrypt(b, utils.MD5(keyName))
			if err == nil {
				auth.Password = rs
			} else {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
