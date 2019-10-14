package weather

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	wws := []Inform{
		{Pro: "10118",
			District: "06",
			City:     "08",
			Info:     "固始县"},
		{Pro: "10118",
			District: "01",
			City:     "01",
			Info:     "郑州市"}}
	for _, w := range wws {
		ws := GetWeather(w)
		info := GetRemindInfo(ws)
		fmt.Println(GetToString(ws, w))
		if info != nil {
			if strings.Compare(info.CoolingInfo, "") != 0 {
				fmt.Println(info.CoolingInfo)
			}
			if strings.Compare(info.WillRainInfo, "") != 0 {
				fmt.Println(info.WillRainInfo)
			}
		} else {
			fmt.Println("明天是晴天！")
		}
	}
}
