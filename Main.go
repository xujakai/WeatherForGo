package main

import (
	"./log"
	"./push"
	"./weather"
	"github.com/pmylund/go-bloom"
	"github.com/robfig/cron"
	"strings"
	"time"
)

var f = bloom.New(10000, 0.001)

var sts = []string{"雨", "雪"}

func main() {
	c := cron.New()
	c.AddFunc("0 0 9,18 * * ?", func() {
		wws := []weather.Inform{
			{Pro: "10118",
				District: "06",
				City:     "08",
				Info:     "固始县"},
			{Pro: "10118",
				District: "01",
				City:     "01",
				Info:     "郑州市"},
		}
		for _, w := range wws {
			ws := weather.GetWeather(w)

			hour := time.Now().Hour()
			if hour >= 17 && hour <= 19 {

				info := weather.GetRemindInfo(ws)
				if info != nil {
					if strings.Compare(info.CoolingInfo, "") != 0 {
						push.SendMsg(info.CoolingInfo)
					}
					if strings.Compare(info.WillRainInfo, "") != 0 {
						push.SendMsg(info.WillRainInfo)
					}
				} else {
					log.Log("明天是晴天！")
				}
			} else {
				push.SendMsg(weather.GetToString(ws, w))
			}
		}
	})
	c.AddFunc("0 0,15,30,45 * * * ? ", func() {
		hour := time.Now().Hour()
		if hour < 6 || hour > 22 {
			return
		}
		ws := []weather.Inform{
			{Pro: "10118",
				District: "06",
				City:     "08",
				Info:     "固始县"},
			{Pro: "10102",
				District: "",
				City:     "",
				Info:     "上海市"},
			{Pro: "10118",
				District: "01",
				City:     "01",
				Info:     "郑州市"},
		}
		weather.WarningInforms(ws, f)
	})

	c.AddFunc("@daily", func() {
		f.Reset()
	})
	c.Start()
	log.Log("监控程序启动！")
	select {}
}
