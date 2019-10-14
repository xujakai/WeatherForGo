package weather

import (
	"../spider"
	. "github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
)

type MyDate int

const (
	Today    = MyDate(1)
	Tomorrow = MyDate(2)
	D3       = MyDate(3)
	D4       = MyDate(4)
	D5       = MyDate(5)
)

type Weather struct {
	Date             MyDate  `今天、明天`
	DateStr          string  `日期字符串`
	WeatherRecording string  `天气`
	District         string  `地区`
	MaxTemperature   float64 `最高气温`
	MinTemperature   float64 `最低所温`
	WindDirection    string  `风向`
	WindForce        string  `风力`
	CladRank         int     `穿衣指数`
}

func (weather Weather) GetCladRank() int {
	tmp := weather.MaxTemperature
	if tmp > 28.0 {
		return 1
	}
	if tmp > 24.0 {
		return 2
	}
	if tmp > 21.0 {
		return 3
	}
	if tmp > 18.0 {
		return 4
	}
	if tmp > 15.0 {
		return 5
	}
	if tmp > 11.0 {
		return 6
	}
	if tmp > 6.0 {
		return 7
	}
	return 8
}

func (weather Weather) ToString() string {
	if weather.MinTemperature ==weather.MaxTemperature {
		return weather.DateStr + "," + weather.WeatherRecording + ","  + float2string(weather.MinTemperature) + "℃," + weather.WindDirection + " " + weather.WindForce
	}
	return weather.DateStr + "," + weather.WeatherRecording + "," + float2string(weather.MaxTemperature) + "/" + float2string(weather.MinTemperature) + "℃," + weather.WindDirection + " " + weather.WindForce
}

func float2string(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

var r, _ = regexp.Compile("^[0-9\\.]*")

func string2float(s string) float64 {
	s = r.FindString(s)
	f, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0.0
	} else {
		return f
	}
}

func GetWeather(inform Inform) []Weather {
	doc := spider.GetHtml(getUrl(inform))
	content := doc.Find("#7d").Find("ul")
	//fmt.Println(content.Html())
	var weathers []Weather

	var index = 1

	content.Find("li").Each(func(i int, selection *Selection) {
		ti := selection.Find("h1").Text()
		if len(ti) == 0 {
			return
		}
		var weather Weather
		switch index {
		case 1:
			weather.Date = Today
			break
		case 2:
			weather.Date = Tomorrow
			break
		case 3:
			weather.Date = D3
			break
		case 4:
			weather.Date = D4
			break
		case 5:
			weather.Date = D5
			break
		}
		weather.DateStr = ti
		weather.District = inform.Info

		wea := selection.Find(".wea").Text()
		weather.WeatherRecording = wea

		find := selection.Find(".tem")



		var f = find.Find("span")
		weather.MinTemperature = string2float(find.Find("i").Text())
		if f.Nodes !=nil {
			weather.MaxTemperature = string2float(find.Find("span").Text())
		}else {
			weather.MaxTemperature = weather.MinTemperature
		}

		var winTmp []string
		selection.Find(".win>em>span").Each(func(i int, se *Selection) {
			winTmp = append(winTmp, se.AttrOr("title", ""))
		})
		join := strings.Join(ThisRemoveDuplicatesAndEmpty(winTmp), "转")
		weather.WindDirection = join

		weather.WindForce = selection.Find(".win").Find("i").Text()

		weathers = append(weathers, weather)
		index++
	})
	return weathers
}

func getUrl(inform Inform) string {
	return "http://www.weather.com.cn/weather/" + inform.Pro + inform.District + inform.City + ".shtml"
}

/**
 * 数组去重 去空
 */
func ThisRemoveDuplicatesAndEmpty(a []string) (ret []string) {
	aLen := len(a)
	for i := 0; i < aLen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
