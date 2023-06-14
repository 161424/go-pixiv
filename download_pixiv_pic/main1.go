package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	b "github.com/chen/download_pixiv_pic/pkg/Browser"
)

type Cont struct {
	Timestamp string
	Illust    map[string]map[string]interface{}
}

func main() {
	h := "https://www.pixiv.net/artworks/108640782/"
	resp := b.GetPixivPage(h, "")

	fmt.Println(string(resp))
	rep := bytes.NewReader(resp)
	doc, err := goquery.NewDocumentFromReader(rep)
	if err != nil {
		panic(err)
	}
	//str, _ := doc.Find("mate#meta-preload-data").Attr("content")
	//fmt.Printf(str)
	r, _ := doc.Find("meta#meta-preload-data").Attr("content")
	//fmt.Println("+++++++++", r)

	var con = &Cont{}
	json.Unmarshal([]byte(r), &con)
	fmt.Println((con.Illust["108640782"]["illustId"]).(string))

}
