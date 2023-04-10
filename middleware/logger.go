package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	//	软连接
	//linkName := "latest_log.log"
	//	将日志写在文件中 两步即可
	filePath := "log/log.log"
	src, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("打开文件出错", err)
	}
	logger := logrus.New()
	//	日志输出写入文件
	logger.Out = src

	//	设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//	日志分割
	logWriter, _ := retalog.New(
		//	日志文件名按照年月日来保存
		filePath+"%Y%m%d.log",
		//	最大保存时间 保存一周
		retalog.WithMaxAge(7*24*time.Hour),
		//	什么时候分割一次日志  24小时分割一次
		retalog.WithRotationTime(24*time.Hour),
		//	软连接 windows下需要以管理员权限运行
		//retalog.WithLinkName(linkName),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.AddHook(Hook)
	return func(c *gin.Context) {
		//	统计开始时间
		startTime := time.Now()
		//	洋葱模型中间件，类似于入栈的方式
		c.Next()
		stopTime := time.Since(startTime).Milliseconds()
		//	格式化请求话费时间
		spendTime := fmt.Sprintf("%d ms", stopTime)

		//	获取发送请求的域名路径
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		//	获取请求状态码
		statusCode := c.Writer.Status()
		//	获取客户端Ip
		clientIp := c.ClientIP()
		//	获取请求客户端信息，如IE浏览器版本或者chrome浏览器版本
		userAgent := c.Request.UserAgent()
		//	获取参数大小
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		//	获取请求方式
		method := c.Request.Method
		//	获取请求路径
		path := c.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"Status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		if len(c.Errors) > 0 {
			// gin系统内部错误
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			//	警告
			entry.Warn()
		} else {
			entry.Info()
		}

	}
}
