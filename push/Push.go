package push

import (
	"../util"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

var funcMap = make(map[string]func(msg string))

type Push struct {
	Label string `mapstructure:"label"`
	Value string `mapstructure:"value"`
}

var ddUrl = "https://oapi.dingtalk.com/robot/send?access_token="

func (token Push) Dd(msg string) bool {
	content := `{"msgtype": "text",
		"text": {"content": "` + msg + `"}
	}`
	if token.Value == "" {
		log.Error("dd token is empty!")
		return false
	}
	res, err := http.Post(ddUrl+token.Value, "application/json", strings.NewReader(content))
	defer res.Body.Close()
	if err != nil {
		log.Error(err)
		return false
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("发送失败！", err)
		return false
	}
	log.Info(msg, string(body))
	return true
}

func (token Push) Console(msg string) bool {
	fmt.Println(msg)
	return true
}

func callReflect(any interface{}, name string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	if v := reflect.ValueOf(any).MethodByName(name); v.String() == "<invalid Value>" {
		return nil
	} else {
		return v.Call(inputs)
	}
}

func (token Push) Push(msg string) bool {
	value := callReflect(&token, util.Capitalize(token.Label), msg)
	if value != nil {
		return value[0].Bool()
	}
	return false
}
