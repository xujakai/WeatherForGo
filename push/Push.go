package push

import (
	"../spider"
	"../util"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"reflect"
)

var funcMap = make(map[string]func(msg string))

type Push struct {
	Label string `mapstructure:"label"`
	Value string `mapstructure:"value"`
}

type Msg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var ddUrl = "https://oapi.dingtalk.com/robot/send?access_token="

func (token Push) Dd(msg Msg) bool {
	content := `{"msgtype": "text",
		"text": {"content": "` + msg.Content + `"}
	}`
	if token.Value == "" {
		log.Error("dd token is empty!")
		return false
	}
	return spider.PostJson(ddUrl+token.Value, content)
}

var jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary

func (token Push) Hook(msg Msg) bool {
	b, err := jsonIterator.MarshalToString(msg)
	if err != nil {
		log.Error(err)
		return false
	}
	return spider.PostJson(token.Value, b)
}

func (token Push) Console(msg Msg) bool {
	fmt.Println(msg.Content)
	return true
}

func (token Push) ServerChan(msg Msg) bool {
	scUrl := "https://sc.ftqq.com/" + token.Value + ".send"
	content := `{"text": "` + msg.Title + `",
		"desp":` + msg.Content + `}
	}`
	return spider.PostJson(scUrl, content)
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

func (token Push) Push(msg Msg) bool {
	value := callReflect(&token, util.Capitalize(token.Label), msg)
	if value != nil {
		return value[0].Bool()
	}
	return false
}
