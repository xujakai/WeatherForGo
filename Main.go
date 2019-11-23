package main

import (
	"./config"
	"./push"
	"./spider"
	"./util"
	"./weather"
	"encoding/json"
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Log  *config.LogInfo   `mapstructure:"log"`
	Push *[]push.Push      `mapstructure:"push"`
	Info *[]weather.Inform `mapstructure:"noti"`
}

var m = util.NewFilter()

var sts = []string{"雨", "雪"}

func (task Task) alarm() {
	hour := time.Now().Hour()
	if hour < 6 || hour > 22 {
		return
	}
	weather.WarningInforms(*task.Info, *task.Push, m)
}

func (task Task) weatherInfo() {
	for _, w := range *task.Info {
		if w.Report {
			ws := weather.GetWeather(w)
			for _, v := range *task.Push {
				v.Push(weather.GetToMsg(ws, w))
			}
		}
	}
}

func (task Task) remind() {
	for _, w := range *task.Info {
		if w.Remind {
			ws := weather.GetWeather(w)
			info := weather.GetRemindInfo(ws)
			if info != nil {
				for e := range *info.Msg {
					for _, v := range *task.Push {
						v.Push((*info.Msg)[e])
					}
				}
			} else {
				log.Info(w.Info, "明天是晴天！")
			}
		} else {
			log.Info(w.Info, "不做提醒！")
		}

	}
}

var (
	help       = flag.Bool("h", false, "this help！")
	test       = flag.Bool("t", false, "test this config.")
	configName = flag.String("c", "config.yaml", "set config name. default config name is config.yaml")
	query      = flag.String("q", "", "query area code")
	//areaCode   = flag.String("area-code", "", "area code, PS. 101020001")
)

var jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary

func getMap(url string) *map[string]string {
	codeMap := make(map[string]string)
	body := spider.GetResponseBody(url)
	if err := json.Unmarshal([]byte(body), &codeMap); err != nil {
		fmt.Println(url, "解析数据出错:", body)
		return nil
	}
	return &codeMap
}

func codeCompare(codeMap map[string]string, query string) (code, value *string) {
	for k, v := range codeMap {
		if strings.Contains(query, v) {
			return &k, &v
		}
	}
	return nil, nil
}

func stringCompare(all, this string, codeMap map[string]string) (code, value, query *string) {
	if codeMap == nil {
		return nil, nil, nil
	}
	if this != "" {
		p := all[strings.Index(all, this):]
		code, value = codeCompare(codeMap, p)
		if code != nil {
			return code, value, &p
		}
	}
	code, value = codeCompare(codeMap, all)
	return code, value, &all
}

type MyCron struct {
	c *cron.Cron
}

func (c *MyCron) reload(task Task) bool {
	if c.c != nil {
		c.c.Stop()
	}
	c.c = cron.New()
	c.c.AddFunc("0 0 9 * * ?", task.weatherInfo)
	c.c.AddFunc("0 0 18 * * ?", task.remind)
	c.c.AddFunc("0 0,15,30,45 * * * ? ", task.alarm)

	c.c.AddFunc("0 0 0 L * ? ", func() {
		m.Reset()
	})
	c.c.Start()
	return true
}

func readTask(config *config.Config, task *Task) {
	config.GetViperUnmarshal(task)
	for e := range *task.Info {
		(*task.Info)[e].District = util.Add(2, (*task.Info)[e].District)
		(*task.Info)[e].City = util.Add(2, (*task.Info)[e].City)
	}
}

func queryCode(query *string) {
	m1 := getMap("http://www.weather.com.cn/data/city3jdata/china.html?_=" + strconv.FormatInt(time.Now().Unix(), 10) + "667")
	if m1 == nil {
		return
	}
	c1, v1, q1 := stringCompare(*query, "", *m1)
	if c1 == nil {
		return
	}
	m2 := getMap("http://www.weather.com.cn/data/city3jdata/provshi/" + *c1 + ".html?_=" + strconv.FormatInt(time.Now().Unix(), 10) + "667")
	c2, v2, q2 := stringCompare(*q1, *v1, *m2)
	if c2 == nil {
		return
	}
	m3 := getMap("http://www.weather.com.cn/data/city3jdata/station/" + *c1 + *c2 + ".html?_=" + strconv.FormatInt(time.Now().Unix(), 10) + "667")
	c3, _, _ := stringCompare(*q2, *v2, *m3)
	fmt.Printf("省：%s 市：%s 县区：%s", *c1, *c2, *c3)
}

func main() {
	flag.Parse()
	if help != nil && *help {
		flag.Usage()
		return
	}
	if query != nil && *query != "" {
		queryCode(query)
		return
	}

	var myCron MyCron
	var task Task
	config := config.NewConfigByName(*configName)
	readTask(config, &task)
	task.Log.LoggerToFile()

	if test != nil && *test {
		task.Log.LoggerToFile()
		task.weatherInfo()
		task.remind()
		task.alarm()
		fmt.Println("test succeed! ")
		return
	}

	config.WatchConfig(func() {
		var tmpTask Task
		readTask(config, &tmpTask)
		tmpTask.Log.LoggerToFile()
		myCron.reload(tmpTask)
		log.Info("配置文件更新成功！")
	})
	myCron.reload(task)
	log.Info("监控程序启动！")
	select {}
}
