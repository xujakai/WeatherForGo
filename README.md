# WeatherForGo
基于go的天气系统，包含天气预报、天气预警、降温提醒、带伞提醒



具体城市代码表请访问[中国天气网](https://www.weather.com)获取

### 推送方式：

本项目目前采用钉钉机器人方式，通过申请的钉钉token发送消息。[接口文档](https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq)

### 项目配置：

项目默认读取config.yaml文件，请参照文件说明进行配置。

`

```
log:
  path: ./logs
  fileName: log

push:
  - label: dd #钉钉机器人
    value: xxx

noti:
  - pro: 10118 #省
    district: 06 #市
    city: 08 #区
    info: 固始县 #展示标识
    alarm: true #是否预警
    remind: true #前一天有雨提醒
    report: true #当天预报
  - pro: 10102
    district: 01
    city:
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

`

### 未来计划：

完善项目配置，以便更加灵活的使用。引入更多的通知方式


### 此项目依赖：

​    [patrickmn/go-bloom](https://github.com/patrickmn/go-bloom)

​    [json-iterator/go](https://github.com/json-iterator/go)

​    [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)

​    [natefinch/lumberjack](https://github.com/natefinch/lumberjack)

​    [uber-go/zap](https://github.com/uber-go/zap)

​    [pmylund/go-bloom](https://github.com/pmylund/go-bloom)

​    [robfig/cron](https://github.com/robfig/cron)

​    [spf13/viper](https://github.com/spf13/viper)

​    [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify)