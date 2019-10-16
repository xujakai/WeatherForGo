package weather

import (
	"../push"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type WaringCode struct {
	KindMap  map[string]string
	GradeMap map[string]string
}

var waringCode WaringCode

func init() {
	waringCode.KindMap = map[string]string{
		"01": "台风预警",
		"02": "暴雨预警",
		"03": "暴雪预警",
		"04": "寒潮预警",
		"05": "大风预警",
		"06": "沙尘暴预警",
		"07": "高温预警",
		"08": "干旱预警",
		"09": "雷电预警",
		"10": "冰雹预警",
		"11": "霜冻预警",
		"12": "大雾预警",
		"13": "霾预警",
		"14": "道路结冰预警",
		"51": "海上大雾预警",
		"52": "雷暴大风预警",
		"53": "持续低温预警",
		"54": "浓浮尘预警",
		"55": "龙卷风预警",
		"56": "低温冻害预警",
		"57": "海上大风预警",
		"58": "低温雨雪冰冻预警",
		"59": "强对流预警",
		"60": "臭氧预警",
		"61": "大雪预警",
		"62": "强降雨预警",
		"63": "强降温预警",
		"64": "雪灾预警",
		"65": "森林(草原)火险预警",
		"66": "雷暴预警",
		"67": "严寒预警",
		"68": "沙尘预警",
		"69": "海上雷雨大风预警",
		"70": "海上雷电预警",
		"71": "海上台风预警",
		"72": "低温预警",
		"91": "寒冷预警",
		"92": "灰霾预警",
		"93": "雷雨大风预警",
		"94": "森林火险预警",
		"95": "降温预警",
		"96": "道路冰雪预警",
		"97": "干热风预警",
		"98": "空气重污染预警",
		"99": "冰冻预警",
	}

	waringCode.GradeMap = map[string]string{
		"01": "蓝色",
		"02": "黄色",
		"03": "橙色",
		"04": "红色",
		"05": "白色",
	}

}

//RemindInfo 提示信息
type RemindInfo struct {
	Msg *[]push.Msg
}

var sts = []string{"雨", "雪"}

func GetRemindInfo(ws []Weather) *RemindInfo {
	var today, tomorrow *Weather

	for i := 0; i < len(ws); i++ {
		if ws[i].Date == Today {
			today = &ws[i]
		}
		if ws[i].Date == Tomorrow {
			tomorrow = &ws[i]
		}
	}
	if today == nil || tomorrow == nil {
		log.Info("无法生成提示信息！")
		return nil
	}
	var msg []push.Msg
	a := tomorrow.GetCladRank() - today.GetCladRank()
	if a > 0 {
		msg = append(msg, push.Msg{Title: "降温提醒", Content: "明天" + today.District + "有明显降温，降温幅度：" + float2string(today.MaxTemperature-tomorrow.MaxTemperature) + "℃！"})
	}
	if a < 0 {
	}
	var flag string
	for _, v := range sts {
		if strings.Contains(tomorrow.WeatherRecording, v) {
			flag = v
			break
		}
	}
	if strings.Compare("", flag) != 0 {
		msg = append(msg, push.Msg{Title: "有" + flag + "提醒", Content: "明天" + today.District + "有" + flag + ",注意带伞！"})
	}
	if len(msg) == 0 {
		return nil
	}
	var r RemindInfo
	r.Msg = &msg
	return &r
}

func GetToMsg(ws []Weather, inform Inform) push.Msg {
	s := getToString(ws, inform)
	return push.Msg{Title: "【今日天气】", Content: s}
}

func getToString(ws []Weather, inform Inform) string {
	var msg []string
	for _, v := range ws {
		msg = append(msg, v.ToString())
	}
	msg = msg[0:3]
	return "【今日天气】" + time.Now().Format("2006年01月02日") + "天气#" + inform.Info + "，" + strings.Join(msg, "；")
}

func getWaringStr(kind string) string {
	tmp := waringCode.KindMap[kind[0:2]]
	tmp = tmp[0 : len(tmp)-len("预警")]
	return tmp + waringCode.GradeMap[kind[2:4]] + "预警"
}
