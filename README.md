# WeatherForGo
基于go的天气系统，包含天气预报、天气预警、降温提醒、带伞提醒



~~具体城市代码表请访问[中国天气网](https://www.weather.com)获取~~

请用 -q 查询地区代码

```
./WeatherForGo -q 河南信阳固始
省：10118 市：06 县区：08

./WeatherForGo -q 河南省信阳市固始县
省：10118 市：06 县区：08

./WeatherForGo -q 河南信阳
省：10118 市：06 县区：01

./WeatherForGo -q 上海市
省：10102 市：00 县区：01
```



### 推送方式：
serverChan方式: 申请的serverChan Key,直接推送到微信。[接口文档](https://sc.ftqq.com/)

钉钉机器人方式：申请钉钉token,直接推送到钉钉。[接口文档](https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq)

web hook方式：请配制hook地址。hook消息格式：
```
{
   "title":"【今日天气】",
   "content":"【今日天气】2019年10月16日天气#上海市，16日（今天）,阴转小雨,22/16℃,北风 3-4级；17日（明天）,阴转多云,22/18℃,北风转西北风 <3级；18日（后天）,多云转晴,23/18℃,北风 3-4级转<3级"
}
```
console控制台：加一个label console

配制需要推送的目标，push列表会逐一推送

### 项目配置：

项目默认读取config.yaml文件，请参照文件说明进行配置。本项目支持热更新，可直接修改配制文件

```
log:
  path: ./logs
  fileName: weather.log

push:
  - label: dd
    value: xxx
  - label: console
    value:
  - label: hook
    value: https://www.topme.pro/hook
  - label: serverChan
    value: key

noti:
  - pro: 10118 #省
    district: 06 #市
    city: 08 #区
    info: 固始县 #展示标识
    alarm: true #是否预警
    remind: true #前一天有雨提醒
    report: true #当天预报
  - pro: 10102
    district: 00
    city: 01
    info: 上海市
    alarm: true
    remind: true
    report: true
  - pro: 10118
    district: 01
    city: 01
    info: 郑州市
    alarm: true
    remind: true
    report: true
```



### 未来计划：

完善项目配置，以便更加灵活的使用。引入更多的通知方式


### 此项目依赖：

​    [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify)

​    [sirupsen/logrus](https://github.com/sirupsen/logrus)

​    [spf13/viper](https://github.com/spf13/viper)

​    [json-iterator/go](https://github.com/json-iterator/go)

​    [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)

​    [emirpasic/gods](https://github.com/emirpasic/gods)

​    [pmylund/go-bloom](https://github.com/pmylund/go-bloom)

​    [robfig/cron](https://github.com/robfig/cron)