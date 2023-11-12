package user

import (
	"fmt"
	"github.com/chen/download_pixiv_pic/common/addr"
	"github.com/chen/download_pixiv_pic/conf"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type OtherA struct {
	Uname    string
	Password string
	Cookie   string
}

type StateA struct {
	*OtherA
	Status bool //  用来确用户更换的状态，但是好像没什么用了
}

var Client *http.Client

func GetNewStatusA() *StateA {
	r1, r2 := Singin()
	return &StateA{
		OtherA: r1,
		Status: r2,
	}
}

func (s *StateA) HandleCk(str string) {
	if s.Status || len(str) < 20 {
		l.Send(slog.LevelError, "CK Error", log.LogStdouts)
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
	header := addr.Header

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
