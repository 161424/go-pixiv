package addr

import "github.com/chen/download_pixiv_pic/conf"

type proxy struct {
	Ip   string
	Port string
}

var Proxy proxy
var Header header
var l = conf.Conf

type NetWork struct {
	Retry         int
	RetryWait     int
	DownloadDelay int
}

func init() {
	//fmt.Println("123", ConfigData["Header"]["User_agent"])
	Proxy.Ip = l.GetString("NetWork.ProxyIp")
	Proxy.Port = l.GetString("NetWork.ProxyPort")
	Header.UserAgent = l.GetString("Authentication.UserAgent")

}

type header struct {
	UserAgent string
}

//func GetNewHeader() *header {
//	return &header{
//		User_agent: ConfigData["Header"]["User_agent"],
//	}
//}

func NewNetWork() *NetWork {
	return &NetWork{
		Retry:         l.GetInt("NetWork.Retry"),
		RetryWait:     l.GetInt("NetWork.RetryWait"),
		DownloadDelay: l.GetInt("NetWork.DownloadDelay"),
	}
}
