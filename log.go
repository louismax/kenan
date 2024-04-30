package kenan

import (
	"fmt"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"time"
)

var log *logrus.Logger

func LogInit() {
	log = logrus.New()
	//设置日志级别
	log.SetLevel(logrus.TraceLevel)
	// 设置将日志输出到标准输出（默认的输出为stderr,标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)
	// 设置日志格式为json格式
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//日志钩子
	log.AddHook(newLfsLogHook())
	fmt.Printf("***** 日志钩子初始化完成,输出路径:[%s]\n", settings.LogFilePath)
}

func DBG(format string, a ...interface{}) {
	myText := fmt.Sprintf(format, a...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	Line := fmt.Sprintf("%s:%d", file, line)
	Method := f.Name()

	log.WithFields(logrus.Fields{
		"Line":   Line,
		"Method": Method,
	}).Debug(myText)
}

func LogInfo(text string, args ...interface{}) {
	myText := fmt.Sprintf(text, args...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	Line := fmt.Sprintf("%s:%d", file, line)
	Method := f.Name()
	log.WithFields(logrus.Fields{
		"Line":   Line,
		"Method": Method,
	}).Info(myText)
}

// LogWarn 警告日志
func LogWarn(text string, args ...interface{}) {
	myText := fmt.Sprintf(text, args...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	Line := fmt.Sprintf("%s:%d", file, line)
	Method := f.Name()
	log.WithFields(logrus.Fields{
		"Line":   Line,
		"Method": Method,
	}).Warn(myText)
}

// LogError 错误日志
func LogError(text string, args ...interface{}) {
	myText := fmt.Sprintf(text, args...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	Line := fmt.Sprintf("%s:%d", file, line)
	Method := f.Name()
	log.WithFields(logrus.Fields{
		"Line":   Line,
		"Method": Method,
	}).Error(myText)
}

// newLfsLogHook 日志钩子
func newLfsLogHook() logrus.Hook {
	rotationTime := 24 * time.Hour

	maxAge := 72 * time.Hour

	writer, err := rotateLogs.New(
		settings.LogFilePath+"_%Y%m%d"+".log",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotateLogs.WithLinkName(settings.LogFilePath),
		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotateLogs.WithRotationTime(rotationTime), // 日志切割时间间隔
		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		rotateLogs.WithMaxAge(maxAge),
		//rotateLogs.WithRotationCount(maxRemainCnt),
		//rotateLogs.WithMaxAge(time.Minute), // 文件最大保存时间
	)
	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return lfsHook
}
