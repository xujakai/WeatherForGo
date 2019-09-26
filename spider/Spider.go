package spider

import (
	"fmt"
	. "github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var timeout = time.Duration(5 * time.Second) //超时时间5s
var client = &http.Client{
	Timeout: timeout,
}

func GetResponseBodyAddHeader(url string, header map[string]string) string {
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("初始化request失败")
		return ""
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)

	if header != nil {
		for e := range header {
			request.Header.Add(e, header[e])
		}
	}

	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + url + "失败")
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("发送失败！")
		return ""
	}
	return string(body)
}

func GetResponseBody(url string) string {
	return GetResponseBodyAddHeader(url, nil)
}

func GetHtml(url string) *Document {
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("初始化request失败")
		return nil
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + url + "失败")
		return nil
	}
	defer res.Body.Close()
	document, err := NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + url + "失败")
		return nil
	}
	return document
}
