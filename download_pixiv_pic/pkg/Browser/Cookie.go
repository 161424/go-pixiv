package Browser

import (
	"fmt"
	"github.com/chen/download_pixiv_pic/common/conf"
	"github.com/chen/download_pixiv_pic/database/sql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type StateA struct {
	*OtherA
	Status bool //  用来确用户更换的状态，但是好像没什么用了
}

type OtherA struct {
	Uname    string
	Password string
	Cookie   string
}

var Client *http.Client

func init() {
	//fmt.Println("Get")
	url, _ := url.Parse(conf.Proxy.Ip + ":" + conf.Proxy.Port)
	fmt.Println("Proxy", url)
	Client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(url),
		},
		Timeout: 10 * time.Second,
	}
}

func GetNewStatusA() *StateA {
	r1, r2 := Singin()
	return &StateA{
		OtherA: r1,
		Status: r2,
	}
}

func (s *StateA) HandleCk() {
	if s.Status || len((conf.ConfigData["Authentication"]["cookie"]).(string)) < 20 {
		log.Panicln("Cookie err")
	}

}

func Singin() (user *OtherA, status bool) {
	user = &OtherA{
		Uname:    (conf.ConfigData["Authentication"]["username"]).(string),
		Password: (conf.ConfigData["Authentication"]["password"]).(string),
		Cookie:   (conf.ConfigData["Authentication"]["cookie"]).(string),
	}

	var s string
	fmt.Printf("Using Username: %s? ", (conf.ConfigData["Authentication"]["username"]).(string))
	fmt.Println("是否更换账户？(Y?N)")
	fmt.Scanln(&s)
	if s == "Y" {
		fmt.Println("请输入新的账户密码以空格作为分割，例如：账户 密码")
		fmt.Scanln(&user.Uname, &user.Password)

		return user, true
	}

	return user, false
}

func CheckCk(user *OtherA) {
	header := conf.Header

	req, _ := http.NewRequest("GET", "https://www.pixiv.net", nil)
	req.Header.Set("User-Agent", header.UserAgent)

	for _, i := range strings.Split(user.Cookie, ";") {
		a := strings.Split(i, "=")
		req.AddCookie(&http.Cookie{Name: a[0], Value: a[1]})
	}

	resp, err := Client.Do(req)
	if err != nil {
		panic(err)
	}

	rb, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var result = false
	parsed_str := string(rb)
	if strings.Contains(parsed_str, "logout.php") {
		result = true
	} else if strings.Contains(parsed_str, "pixiv.user.loggedIn = true") {
		result = true
	} else if strings.Contains(parsed_str, "_gaq.push(['_setCustomVar', 1, 'login', 'yes'") {
		result = true
	} else if strings.Contains(parsed_str, "var dataLayer = [{ login: 'yes',") {
		result = true
	}

	if result {
		fmt.Println("Logged in using cookie")
		re := regexp.MustCompile("user_id: \\\"(\\d+)\\\",")
		found := re.FindStringSubmatch(parsed_str)
		fmt.Printf("My User Id: %s", found[1])
	}

	log.Panicln("Cookie already expired/invalid.")

}

func GetCk() string {
	return ""
}

func UpdateCk() {

}

func getMyId() {

}

func GetPixivPage(url string, ref string) []byte {
	client := Client
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("conf.Header.UserAgent, %T", conf.Header.UserAgent)
	req.Header.Set("User-Agent", conf.Header.UserAgent)
	req.Header.Set("Referer", ref)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//fmt.Printf(sql.DefaultAuth.Cookies)
	for _, i := range strings.Split(sql.DefaultAuth.Cookies, ";") {
		a := strings.Split(i, "=")
		req.AddCookie(&http.Cookie{Name: a[0], Value: a[1]})
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	r_, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return r_
}
