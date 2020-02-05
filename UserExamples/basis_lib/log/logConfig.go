package log

import (
	"context"
	"fmt"
	"go-micro-learning/UserExamples/basis_lib/configuration"
	oti "go-micro-learning/UserExamples/basis_lib/tracer/opentracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"regexp"
	"sync"
)


/*
export ETCDCTL_API=3
etcdctl --endpoints='127.0.0.1:2380' put  /gateway/config/log_config '{"log_file_name":"info.log","ac_file_name":"access.log","max_size":128,"max_backups":30,"max_age":7,"log_level":"debug"}'
*/


var (
	logger *zap.SugaredLogger
)


// 动态log配置
type DyLogConfig struct {
	LogFilename   string `json:"log_file_name"`
	AcFilename   string `json:"ac_file_name"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	LogLevel   string `json:"log_level"`
}

var (
	dyLogConf *DyLogConfig = &DyLogConfig{}
	initLogConfOnce sync.Once
	m      sync.RWMutex
)

// 初始化
//func init(){
//	basis_lib.Register(InitLogConf)
//	basis_lib.Register(InitAcLogConf)
//}

//func initLogConf() {
//	yamlFileCfg := file_config.Cfg()
//	fmt.Println("读取Etcd配置", yamlFileCfg.Etcd.Addr, yamlFileCfg.Etcd.PathPrefix, yamlFileCfg.Etcd.LogPathKey)
//	logConfEtcd := etcd.InitEtcdConfig(
//		etcd.WithAddress(yamlFileCfg.Etcd.Addr),
//		etcd.WithPathPrefix(yamlFileCfg.Etcd.PathPrefix),
//		etcd.WithPathKey(yamlFileCfg.Etcd.LogPathKey),
//		etcd.WithValStruct(dyLogConf),
//		etcd.WithHookFunc([]func(){initLog,initAcLog}), // 动态配置修改完毕后需要执行的函数
//	)
//	logConfEtcd.InitEtcdConfig()
//	logConfEtcd.GetKey()
//}


func InitLog() {
	cfg := configuration.C()
	fmt.Printf("log cfg %p \n",cfg)

	err := cfg.GetEtcdCfg("log_config",dyLogConf)
	if err != nil {
		fmt.Println("InitLog GetEtcdCfg",err,dyLogConf)
		os.Exit(1)
	}

	cfg.EtcdAutoUpdateCfg("log_config",dyLogConf,[]func(){initLog,initAcLog}) // 动态更新
	initLog()
	initAcLog()
}

//func init(){
//	basis_lib.Register(InitLog)
//}

func initLog() {
	m.Lock()
	defer m.Unlock()


	//initLogConfOnce.Do(func() {
	//	initLogConf()
	//})

	r1 := regexp.MustCompile("^debug$|^info$|^warn$|^error$")
	b1 := r1.MatchString(dyLogConf.LogLevel)
	if !b1 {
		fmt.Println("日志级别不对！")
		return
	}

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

	hook := lumberjack.Logger{
		//Filename:   "./logs/micro-srv.log", // 日志文件路径
		Filename:   dyLogConf.LogFilename,   // 日志文件路径
		MaxSize:    dyLogConf.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: dyLogConf.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     dyLogConf.MaxAge,     // 文件最多保存多少天
		Compress:   true,                 // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 不是全路径		// zapcore.FullCallerEncoder 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	//atomicLevel.SetLevel(zap.DebugLevel)
	logLevelText := dyLogConf.LogLevel
	logLevelByte := []byte(logLevelText)
	unmarshalTextErr := atomicLevel.UnmarshalText(logLevelByte)

	if unmarshalTextErr != nil {
		fmt.Println("unmarshalTextErr logs  err ", unmarshalTextErr)
		os.Exit(1)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),               // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),zapcore.AddSync(&hook)), // zapcore.AddSync(os.Stdout),打印到控制台 zapcore.AddSync(&hook)文件
		atomicLevel,                                         // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()             // 当前栈
	addCallerSkip := zap.AddCallerSkip(1) // 向上跳一层
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	//filed := zap.Fields(zap.String("serviceName", "serviceName"))

	// 构造日志
	//log := zap.New(core, caller, addCallerSkip, development, filed)
	log := zap.New(core, caller, addCallerSkip, development)
	logger = log.Sugar()
	logger.With()
	Warn("initLog 初始化成功")

}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)

}



func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func ErrorWith(ctx context.Context,args ...interface{}){
	keyId := withFieldTraceIdKeyContext(ctx)
	logger.With(zap.String("Dy-Trace-Id",keyId)).Error(args...)
}

func ErrorFWith(ctx context.Context,template string, args ...interface{}){
	keyId := withFieldTraceIdKeyContext(ctx)
	logger.With(zap.String("Dy-Trace-Id",keyId)).Errorf(template, args...)
}


func WarnWith(ctx context.Context,args ...interface{}){
	keyId := withFieldTraceIdKeyContext(ctx)
	logger.With(zap.String("Dy-Trace-Id",keyId)).Warn(args...)
}

func WarnFWith(ctx context.Context,template string, args ...interface{}){
	keyId := withFieldTraceIdKeyContext(ctx)
	logger.With(zap.String("Dy-Trace-Id",keyId)).Warnf(template, args...)
}


func InfoWith(ctx context.Context,args ...interface{}){
	keyId := withFieldTraceIdKeyContext(ctx)
	logger.With(zap.String("Dy-Trace-Id",keyId)).Info(args...)
}

func InfoFWith(ctx context.Context,template string, args ...interface{}){
	keyId := withFieldTraceIdKeyContext(ctx)
	logger.With(zap.String("Dy-Trace-Id",keyId)).Infof(template, args...)
}

func DebugWith(ctx context.Context,args ...interface{}){
	keyId := withFieldTraceIdKeyContext(ctx)
	logger.With(zap.String("Dy-Trace-Id",keyId)).Debug(args...)
}

func DebugFWith(ctx context.Context,template string, args ...interface{}){
	keyId := withFieldTraceIdKeyContext(ctx)
	logger.With(zap.String("Dy-Trace-Id",keyId)).Debugf(template, args...)
}




func withFieldTraceIdKeyContext(ctx context.Context) (string) {
	keyId := ctx.Value(oti.TraceIdKey{})
	if keyId == nil {
		return ""
	}
	keyIdStr := keyId.(string)
	return keyIdStr
}