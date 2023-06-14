package conf

type proxy struct {
	Ip   string
	Port string
}

var Proxy proxy
var Header header

type NetWork struct {
	Retry         int
	RetryWait     int
	DownloadDelay int
}

func init() {
	//fmt.Println("123", ConfigData["Header"]["User_agent"])
	Proxy.Ip = (ConfigData["NetWork"]["ProxyIp"]).(string)
	Proxy.Port = (ConfigData["NetWork"]["ProxyPort"]).(string)
	Header.UserAgent = (ConfigData["Header"]["UserAgent"]).(string)

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
		Retry:         (ConfigData["NetWork"]["Retry"]).(int),
		RetryWait:     (ConfigData["NetWork"]["RetryWait"]).(int),
		DownloadDelay: (ConfigData["NetWork"]["DownloadDelay"]).(int),
	}
}
