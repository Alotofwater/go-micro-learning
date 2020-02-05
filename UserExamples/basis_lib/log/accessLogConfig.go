package log

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
)

var (
	// the local logger
	acLogger *zap.SugaredLogger
)


func initAcLog(){
	//fileName := "micro-srv.log"
	//syncWriter := zapcore.AddSync(&lumberjack.Logger{
	//	Filename:  fileName,
	//	MaxSize:   128, //MB
	//	LocalTime: true,
	//	Compress:  true,
	//})
	//encoder := zap.NewProductionEncoderConfig()
	//encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//core := zapcore.NewCore(
	//	zapcore.NewJSONEncoder(encoder),
	//	syncWriter,
	//	zap.NewAtomicLevelAt(zap.DebugLevel),
	//	)
	//log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	//logger = log.Sugar()

	acHook := lumberjack.Logger{
		Filename:   dyLogConf.AcFilename,   // 日志文件路径
		MaxSize:    dyLogConf.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: dyLogConf.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     dyLogConf.MaxAge,     // 文件最多保存多少天
		Compress:   true,                     // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		//LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		//StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 不是全路径		// zapcore.FullCallerEncoder 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	acAtomicLevel := zap.NewAtomicLevel()
	//atomicLevel.SetLevel(zap.DebugLevel)
	acLogLevelText := "debug"
	acLogLevelByte := []byte(acLogLevelText)
	acUnmarshalTextErr := acAtomicLevel.UnmarshalText(acLogLevelByte)

	if acUnmarshalTextErr != nil {
		fmt.Println("acUnmarshalTextErr logs  err ",acUnmarshalTextErr)
		os.Exit(1)
	}

	acCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&acHook)), // 打印到控制台和文件
		acAtomicLevel,                                                                     // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	//acCaller := zap.AddCaller()
	// 开启文件及行号
	//acDevelopment := zap.Development()
	// 设置初始化字段
	//acFiled := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	//acLog := zap.New(acCore, acCaller, acDevelopment, acFiled)
	acLog := zap.New(acCore)
	acLogger = acLog.Sugar()
}



func AcWarnf(template string, args ...interface{}) {
	acLogger.Warnf(template,args...)
}


func AcHttpDebug(req  *http.Request,processingTime string,statusCode int) {

	traceId := withFieldTraceIdKeyContext(req.Context())
	acLogger.Debugw("access",
		zap.String("Dy-Trace-Id",traceId),
		zap.String("method",req.Method),
		zap.String("proto",req.Proto),
		zap.String("url",req.RequestURI),
		zap.String("user_agent",req.Header.Get("User-Agent")),
		zap.String("content_length",fmt.Sprintf("%d",req.ContentLength)),
		zap.String("processing_time",processingTime),
		zap.String("user_host",req.Host),
		zap.Int("status_code",statusCode),
		)
}

func AcGinDebug(c  *gin.Context,processingTime string) {
	req := c.Request
	traceId := withFieldTraceIdKeyContext(req.Context())
	acLogger.Debugw("access",
		zap.String("Dy-Trace-Id",traceId),
		zap.String("method",req.Method),
		zap.String("proto",req.Proto),
		zap.String("url",req.RequestURI),
		zap.String("user_agent",req.Header.Get("User-Agent")),
		zap.String("content_length",fmt.Sprintf("%d",c.Writer.Size())),
		zap.String("processing_time",processingTime),
		zap.String("user_host",req.Host),
		zap.Int("status_code",c.Writer.Status()),
	)
}

func AcGrpcDebug(ctx  context.Context,processingTime string) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	traceId := withFieldTraceIdKeyContext(ctx)
	acLogger.Debugw("access",
		zap.String("Dy-Trace-Id",traceId),
		zap.String("Micro-Method",md["Micro-Method"]),
		zap.String("Remote",md["Remote"]),
		zap.String("User-Agent",md["User-Agent"]),
		zap.String("Micro-Endpoint",md["Micro-Endpoint"]),
		zap.String("processing_time",processingTime),
	)
}

func AcDebugw(req  *http.Request) {
	acLogger.Debugw("access",
		zap.String("method",req.Method),
		zap.String("url",req.URL.Path),
		zap.String("host",req.URL.Opaque),
		zap.String("Hostname",req.URL.RawQuery),
	)
}

//func AcWith(){
//
//	acLogger.With(
//		zap.String("hello", "world"),
//		zap.String("failure", "oh no"),
//		zap.Int("count", 42),
//	)
//   //"hello", "world",
//   //"failure", "ccccc",
//   //"count", 42,)
//
//}

