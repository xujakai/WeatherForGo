package main

import (
	"./config"
	"./push"
	"./util"
	"./weather"
	"flag"
	"fmt"
	"github.com/pmylund/go-bloom"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Task struct {
	Log  *config.LogInfo   `mapstructure:"log"`
	Push *[]push.Push      `mapstructure:"push"`
	Info *[]weather.Inform `mapstructure:"noti"`
}

var f = bloom.New(10000, 0.001)

var sts = []string{"雨", "雪"}

func (task Task) alarm() {
	hour := time.Now().Hour()
	if hour < 6 || hour > 22 {
		return
	}
	weather.WarningInforms(*task.Info, *task.Push, f)
}

func (task Task) weatherInfo() {
	for _, w := range *task.Info {
		ws := weather.GetWeather(w)
		if w.Report {
			for _, v := range *task.Push {
				v.Push(weather.GetToString(ws, w))
			}
		}
	}
}

func (task Task) remind() {
	for _, w := range *task.Info {
		ws := weather.GetWeather(w)
		info := weather.GetRemindInfo(ws)
		if info != nil && w.Remind {
			if strings.Compare(info.CoolingInfo, "") != 0 {
				for _, v := range *task.Push {
					v.Push(info.CoolingInfo)
				}
			}
			if strings.Compare(info.WillRainInfo, "") != 0 {
				for _, v := range *task.Push {
					v.Push(info.WillRainInfo)
				}
			}
		} else {
			if w.Remind {
				log.Info("明天是晴天！")
			} else {
				log.Info("不做提醒！")
			}
		}
	}
}

var (
	help       = flag.Bool("h", false, "this help！")
	test       = flag.Bool("t", false, "test run this project！")
	configName = flag.String("c", "config", "config name")
)

func main() {
	flag.Parse()
	if help != nil && *help {
		flag.Usage()
		return
	}

	if test != nil && *test {
		fmt.Println("run test")
		info := &config.LogInfo{"./", "test.log"}
		p := []push.Push{{Label: "console"}}
		w := []weather.Inform{{Pro: "10102", District: "01", City: "00", Info: "上海市", Alarm: true, Remind: true, Report: true}}
		var task Task
		task.Log = info
		task.Push = &p
		task.Info = &w
		task.Log.LoggerToFile()

		task.weatherInfo()
		task.remind()
		task.alarm()
		return
	}

	var task Task
	config := config.NewConfigByName(*configName)
	config.GetViperUnmarshal(&task)
	task.Log.LoggerToFile()
	for e := range *task.Info {
		(*task.Info)[e].District = util.Add(2, (*task.Info)[e].District)
		(*task.Info)[e].City = util.Add(2, (*task.Info)[e].City)
	}
	c := cron.New()

	c.AddFunc("0 0 9 * * ?", task.weatherInfo)
	c.AddFunc("0 0 18 * * ?", task.remind)
	c.AddFunc("0 0,15,30,45 * * * ? ", task.alarm)

	c.AddFunc("@daily", func() {
		f.Reset()
	})
	c.Start()
	log.Info("监控程序启动！")
	select {}
}
