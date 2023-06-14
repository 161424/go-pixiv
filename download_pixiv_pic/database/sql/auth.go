package sql

import (
	conf2 "github.com/chen/download_pixiv_pic/common/conf"
	"github.com/chen/download_pixiv_pic/common/model/ip"
	"gorm.io/gorm"
)

var DefaultAuth *Auth

type Auth struct {
	gorm.Model
	Uname    string `gorm:"column:uname;NOT NULL;type:varchar(32);comment:登录账户"json:"uname"`
	Password string `gorm:"column:password;NOT NULL;type:varchar(32);comment:登录密码"json:"password"`
	Cookies  string `gorm:"column:cookies;comment:cookie"json:"password"`
	InterIP  string `gorm:"column:InterIP;default:null;type:cidr;comment:InterIP"json:"ip_1"`
	EnterIP  string `gorm:"column:EnterIP;default:null;type:inet;comment:EnterIP"`
	Mac      string `gorm:"column:macaddr;default:null;type:macaddr;comment:macaddr"`
	Output   string `gorm:"column:output;default:null"` // 执行结果
	//RunTimer time.Time `gorm:"column:run_timer;default:null"` // 执行时间
	CostTime float64 `gorm:"column:cost_time"` // 执行耗时
	Status   int     `gorm:"column:status;NOT NULL"`
}

//func NewAuth(){
//
//}

func init() {
	InitAuth()
}

func InitAuth() (h *Auth) {
	//fmt.Println((conf2.ConfigData["Authentication"]["username"]).(string))
	DefaultAuth = &Auth{
		Uname:    (conf2.ConfigData["Authentication"]["username"]).(string),
		Password: (conf2.ConfigData["Authentication"]["password"]).(string),
		Cookies:  (conf2.ConfigData["Authentication"]["cookie"]).(string),
		InterIP:  ip.GetInterIP(),
		EnterIP:  ip.GetExternalIP(),
		Mac:      ip.GetMacAddr(),
		Output:   "",
		//RunTimer: time.Now(),
		CostTime: 0.0,
		Status:   0,
	}
	return
}

func DeleAuthAll_id(db *gorm.DB, id int) {
	db.Delete(&Auth{}, id)
}

func DeleAuthAll_hard(db *gorm.DB) {
	db.Unscoped().Delete(&Auth{})
	//db.Exec("DELETE FROM auths")
}

func GetDAuth() *Auth {
	return DefaultAuth
}
