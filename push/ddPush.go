package push

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

var ddUrl = "https://oapi.dingtalk.com/robot/send?access_token="

func sendDdMsg(token, msg string) {
	content := `{"msgtype": "text",
		"text": {"content": "` + msg + `"}
	}`
	res, err := http.Post(ddUrl+token, "application/json", strings.NewReader(content))
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("发送失败！")
		return
	}
	log.Info(msg, string(body))
}
