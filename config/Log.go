package config

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

type LogInfo struct {
	LogPath  string `mapstructure:"path"`
	FileName string `mapstructure:"fileName"`
}
type logFileWriter struct {
	logPath      string
	fileName     string
	file         *os.File
	size         int64
	thisFileName string
}

func (c LogInfo) LoggerToFile() {
	if !isExist(c.LogPath) {
		os.MkdirAll(c.LogPath, os.ModePerm)
	}
	name := getFileName(c.LogPath, c.FileName)
	file, i := getFile(name)

	fileWriter := logFileWriter{logPath: c.LogPath, fileName: c.FileName, file: file, size: i, thisFileName: name}
	logger.SetOutput(&fileWriter)

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})
}

func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	name := getFileName(p.logPath, p.fileName)

	if strings.Compare(p.thisFileName, name) == 0 {
		n, e := p.file.Write(data)
		p.size += int64(n)
		return n, e
	} else {
		p.file.Close()
		file, i := getFile(name)
		p.file = file
		p.size = i
		n, e := p.file.Write(data)
		return n, e
	}
	/*//文件最大 64K byte
	if p.size > 1024*64 {
		p.file.Close()
		fmt.Println("log file full")
		currentTime := time.Now()
		p.file, _ = os.OpenFile(path.Join(p.logPath, p.fileName)+"."+currentTime.Format("2006-01-02")+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		p.size = 0
	}*/
}

func getFile(fileName string) (*os.File, int64) {
	src, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil && os.IsNotExist(err) {
		file, er := os.Create(fileName)
		if er != nil {
			fmt.Println("create log file failed!")
		}
		return file, 0

	} else {
		info, err := src.Stat()
		if err != nil {
			fmt.Println(err)
		}
		return src, info.Size()
	}
}

func getFileName(logPath, logName string) string {
	currentTime := time.Now()
	return path.Join(logPath, logName) + "." + currentTime.Format("2006-01-02") + ".log"
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
