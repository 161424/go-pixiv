package browser

import (
	"bytes"
	"fmt"
	"github.com/chen/download_pixiv_pic/common/addr"
	"github.com/chen/download_pixiv_pic/dao/sql"
	"github.com/chen/download_pixiv_pic/log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//var ck struct {
//	reps sync.Map
//}

var Client *http.Client

func init() {
	//fmt.Println("Get")
	sockets, _ := url.Parse(addr.Proxy.Ip + ":" + addr.Proxy.Port)
	l.Send(slog.LevelInfo, fmt.Sprintf("Proxy:%s", sockets), log.LogFiles|log.LogStdouts, )
	Client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(sockets),
		},
		Timeout: 10 * time.Second,
	}
}

func hashString(s string) uint64 {
	var h uint64 = 14695981039346656037 // offset
	for i := 0; i < len(s); i++ {
		h = h ^ uint64(s[i])
		h = h * 1099511628211 // prime
	}
	return h
}

func GetPixivPage(urls string, ref string) []byte {
	client := Client
	var req *http.Request
	var err error

	req, err = http.NewRequest("GET", urls, nil)
	if err != nil {
		return nil
	}
	//fmt.Printf("conf.Header.UserAgent, %T", conf.Header.UserAgent)
	req.Header.Set("User-Agent", addr.Header.UserAgent)
	req.Header.Set("Referer", ref)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for _, i := range strings.Split(sql.A.Cookies, ";") {
		a := strings.Split(i, "=")
		req.AddCookie(&http.Cookie{Name: a[0], Value: a[1]})
	}
	//ck.reps.Store(hashString(sql.DefaultAuth().Cookies),req)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	bts := bytes.Buffer{}

	bts.ReadFrom(resp.Body)
	//fmt.Println(a, b, bts.Bytes(), resp.Body)
	defer resp.Body.Close()
	return bts.Bytes()
}
