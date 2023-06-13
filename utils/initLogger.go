package utils

import (
	filerotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

func InitLogger() *zap.SugaredLogger {
	var core zapcore.Core // 定义核心
	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel // 定义 info 级别的日志写入器
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel // 定义 warn 级别的日志写入器
	})

	_logpath := viper.GetStringMapString("LOG_PATH") // 获取日志文件路径名
	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(_logpath["info"])  // 获取信息日志记录器的写入器（文件）
	warnWriter := getWriter(_logpath["error"]) // 获取错误日志记录器的写入器（文件）

	encoder := getEncoder()
	// 根据环境变量选择输出方式
	if viper.GetString("ENV") == "dev" {
		// 如果是开发环境，则使用多重输出流和文件输出
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(infoWriter), zapcore.AddSync(os.Stdout)), infoLevel),
			zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(warnWriter), zapcore.AddSync(os.Stdout)), warnLevel),
		)
	} else {
		// 如果是其他环境，则只使用文件输出
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		)
	}
	// 创建 Logger，用于记录日志
	return zap.New(core, zap.AddCaller()).Sugar() // AddCaller 表示在日志中添加调用者的信息
}

func getEncoder() zapcore.Encoder {

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder, // 设置关键词大写
		TimeKey:     "time",
		CallerKey:   "file",
		//EncodeTime:  zapcore.ISO8601TimeEncoder, //格式化输出时间格式
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder, // 只输出文件名和行号
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000) // 输出执行时间，单位为毫秒
		},
	}

	var encoder zapcore.Encoder
	// 根据环境变量 选择日志输出格式
	if viper.GetString("ENV") == "dev" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 控制台 格式输出
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // JSON 格式输出
	}

	return encoder
}

func getWriter(filename string) io.Writer {
	hook, err := filerotatelogs.New(
		filename+".%Y%m%d",
		filerotatelogs.WithLinkName(filename),
		filerotatelogs.WithMaxAge(time.Hour*24*7),     // 日志保留7天
		filerotatelogs.WithRotationTime(time.Hour*24), // 每隔1天生产一个日志文件
	)

	if err != nil {
		panic(err)
	}
	return hook
}
