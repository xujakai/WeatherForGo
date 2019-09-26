package weather

import (
	"../log"
	"../push"
	"../spider"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pmylund/go-bloom"
	"strconv"
	"strings"
	"time"
)

type ResData struct {
	Count     string     `json:"count"`
	DataArray [][]string `json:"data"`
}

func (data ResData) getWarning(pro string, district string, city string) *[]Warning {
	var ws []Warning
	cityCode := pro + district + city
	for _, v := range data.DataArray {
		tmpInfo := v[1][0:strings.Index(v[1], ".")]
		split := strings.Split(tmpInfo, "-")
		tmpCityCode := split[0]
		if strings.Compare(tmpCityCode, cityCode) == 0 {

			var warningStr string
			var warningTimeStr string
			var warningCity string
			var warningUil string

			warningStr = split[2]
			warningTimeStr = split[1]
			warningUil = v[1]
			warningCity = v[0]

			warningTimeStr = warningTimeStr[0:4] + "年" + warningTimeStr[4:6] + "月" + warningTimeStr[6:8] + "日 " + warningTimeStr[8:10] + "：" + warningTimeStr[10:12] + "：" + warningTimeStr[12:]
			warning := Warning{City: warningCity, Url: warningUil, Time: warningTimeStr, Info: getWaringStr(warningStr)}
			infoString := warning.getWarningInfoString()
			if infoString == nil {
				s := warning.getWarningInfoStringPro()
				warning.Content = s
			} else {
				warning.Content = infoString
			}
			ws = append(ws, warning)
		}
	}
	return &ws

}

type Warning struct {
	City    string
	Url     string
	Time    string
	Info    string
	Content *string
}

type WarningInfo struct {
	Head         string `json:"head"`
	AlertId      string `json:"ALERTID"`
	Province     string `json:"PROVINCE"`
	City         string `json:"CITY"`
	StationName  string `json:"STATIONNAME"`
	SignalType   string `json:"SIGNALTYPE"`
	SignalLevel  string `json:"SIGNALLEVEL"`
	TypeCode     string `json:"TYPECODE"`
	LevelCode    string `json:"LEVELCODE"`
	IssueTime    string `json:"ISSUETIME"`
	IssueContent string `json:"ISSUECONTENT"`
	Underwriter  string `json:"UNDERWRITER"`
	RelieveTime  string `json:"RELIEVETIME"`
	NameEn       string `json:"NAMEEN"`
	YjtypeEn     string `json:"YJTYPE_EN"`
	YjycEn       string `json:"YJYC_EN"`
	Time         string `json:"TIME"`
	Effect       string `json:"EFFECT"`
}

func getWeatherWarningResData() *ResData {
	body := spider.GetResponseBody("http://product.weather.com.cn/alarm/grepalarm_cn.php?" + strconv.FormatInt(time.Now().Unix(), 10) + "667")
	bo := body[len("var alarminfo=") : len(body)-1]
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	var data ResData
	if err := json.Unmarshal([]byte(bo), &data); err != nil {
		log.Log("解析天气数据出错")
		return nil
	}
	return &data
}

type Inform struct {
	Pro      string
	District string
	City     string
	Info     string
}

func WarningInforms(informs []Inform, f *bloom.Filter) {
	data := getWeatherWarningResData()
	for _, v := range informs {
		warnings := data.getWarning(v.Pro, v.District, v.City)
		if warnings == nil {
			log.Log(v.Info + "无预警信息！")
		} else {
			for _, warning := range *warnings {
				pushStr := "【" + warning.Info + "】" + warning.Time + " " + warning.City + "#" + *warning.Content
				bytes := []byte(pushStr)
				if !f.Test(bytes) {
					push.SendMsg(pushStr)
					f.Add(bytes)
				}
			}
		}
	}
}

func (code Warning) getWarningInfoStringPro() *string {
	//Referer: http://www.weather.com.cn/alarm/newalarmcontent.shtml?file=10102-20190906080000-0101.html

	header := make(map[string]string)
	header["Referer"] = "http://www.weather.com.cn/alarm/newalarmcontent.shtml?file=" + code.Url
	header["Accept"] = `text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01`
	//header["Accept-Encoding"]=`gzip, deflate`
	//header["Accept-Language"]=`zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7`
	header["Connection"] = `keep-alive`
	//header["Cookie"]=`userNewsPort0=1; f_city=%E4%B8%8A%E6%B5%B7%7C101020100%7C; zs=101020100%7C%7C%7Cyd-uv; defaultCty=101020100; defaultCtyName=%u4E0A%u6D77; Wa_lvt_20=1567490146; Wa_lvt_1=1567503860,1567571525,1567580597,1567757199; Wa_lpvt_1=1567757216`
	header["Host"] = `www.weather.com.cn`
	header["X-Requested-With"] = `XMLHttpRequest`

	s := strings.Split(code.Url, "-")[2]
	body := spider.GetResponseBodyAddHeader("http://www.weather.com.cn/data/alarminfo/"+s+"?_="+strconv.FormatInt(time.Now().Unix(), 10)+"667", header)
	bo := body[len("var alarminfo="):]
	var data []string
	if err := json.Unmarshal([]byte(bo), &data); err != nil {
		fmt.Println("解析天气数据出错")
		return nil
	}
	return &data[2]
}

func (code Warning) getWarningInfoString() *string {
	body := spider.GetResponseBody("http://product.weather.com.cn/alarm/webdata/" + code.Url + "?_=" + strconv.FormatInt(time.Now().Unix(), 10) + "667")
	bo := body[len("var alarminfo="):]
	var data WarningInfo
	if err := json.Unmarshal([]byte(bo), &data); err != nil {
		fmt.Println("解析天气数据出错")
		return nil
	}
	return &data.IssueContent
}
