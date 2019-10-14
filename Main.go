package main

import (
	"./config"
	"./log"
	"./push"
	"./util"
	"./weather"
	"fmt"
	"github.com/pmylund/go-bloom"
	"github.com/robfig/cron"
	"strings"
	"time"
)

type Task struct {
	LogPath *LogInfo          `mapstructure:"log"`
	Push    *[]push.PushToken `mapstructure:"push"`
	Info    *[]weather.Inform `mapstructure:"noti"`
}

type LogInfo struct {
	Path string `mapstructure:"path"`
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

func (task Task) run() {
	for _, w := range *task.Info {
		ws := weather.GetWeather(w)
		hour := time.Now().Hour()
		if hour >= 17 && hour <= 19 {
			info := weather.GetRemindInfo(ws)
			if info != nil && w.Remind {
				if strings.Compare(info.CoolingInfo, "") != 0 {
					for _, v := range *task.Push {
						v.Push(info.CoolingInfo)
					}
				}
				if strings.Compare(info.WillRainInfo, "") != 0 {
					for _, v := range *task.Push {
						v.Push(info.CoolingInfo)
					}
				}
			} else {
				if w.Remind {
					log.Log("明天是晴天！")
				} else {
					log.Log("不做提醒！")
				}
			}
		} else {
			if w.Report {
				for _, v := range *task.Push {
					v.Push(weather.GetToString(ws, w))
				}
			}
		}
	}
}

func main() {
	var task Task
	config.GetViperUnmarshal(&task)
	value := config.GetValue("test")

	fmt.Println(value)
	for e := range *task.Info {
		(*task.Info)[e].District = util.Add(2, (*task.Info)[e].District)
		(*task.Info)[e].City = util.Add(2, (*task.Info)[e].City)
	}
	c := cron.New()
	//task.run()
	task.alarm()
	//c.AddFunc("0 0 9,18 * * ?", task.run)
	//c.AddFunc("0 0,15,30,45 * * * ? ", task.alarm)

	c.AddFunc("@daily", func() {
		f.Reset()
	})
	c.Start()
	log.Log("监控程序启动！")
	select {}
}
