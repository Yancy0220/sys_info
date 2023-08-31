package logs

import (
	"fiber/pkg/setting"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"time"
)

var Logger *zap.Logger

// 初始化日志 logger
func init() {
	logPath := "info"
	errPath := "error"
	level := "debug"
	//创建夹
	//date := time.Now().Format("20060102150405")
	infoPath := "logs/info/"
	errorPath := "logs/error/"
	// 判断日志路径是否存在，如果不存在就创建
	if exist := isExist(infoPath); !exist {
		if err := os.MkdirAll(infoPath, os.ModePerm); err != nil {
		}
	}
	if exist := isExist(errorPath); !exist {
		if err := os.MkdirAll(errorPath, os.ModePerm); err != nil {
		}
	}
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder, //将级别转换成大写
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	encoder := zapcore.NewConsoleEncoder(config)
	// 设置级别
	logLevel := zap.DebugLevel
	switch level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}
	// 实现两个判断日志等级的interface  可以自定义级别展示
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(logPath, infoPath)
	warnWriter := getWriter(errPath, errorPath)

	// 最后创建具体的Logger
	var core zapcore.Core
	if setting.RunMode == "dev" {
		core = zapcore.NewTee(
			// 将info及以下写入logPath,  warn及以上写入errPath
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
			//日志都会在console中展示
			zapcore.NewCore(zapcore.NewConsoleEncoder(config),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel),
		)
	} else {
		core = zapcore.NewTee(
			// 将info及以下写入logPath,  warn及以上写入errPath
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
			//日志都会在console中展示
			//zapcore.NewCore(zapcore.NewConsoleEncoder(config),
			//zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel),
		)
	}
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel)) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
}

func getWriter(filename, path string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	if path == "" {
		path = "logs/"
	}
	rotateLogsFile := path + filename + ".%Y%m%d%H.log"
	linkName := path + filename + ".log"
	hook, err := rotatelogs.New(
		rotateLogsFile,
		rotatelogs.WithLinkName(linkName),
		rotatelogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotatelogs.WithRotationTime(time.Hour*24), //切割频率 6小时
	)
	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// logs.Debug(...)
func Debug(format string, v ...interface{}) {
	Logger.Sugar().Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	Logger.Sugar().Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	Logger.Sugar().Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	Logger.Sugar().Errorf(format, v...)
}

func Panic(format string, v ...interface{}) {
	Logger.Sugar().Panicf(format, v...)
}
