package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var SugarLogger *zap.SugaredLogger

func init() {
	InitLogger()

}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	SugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logs/log.txt",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}

	//file, _ := os.Create("logs/log.txt")
	fileWriteSyncer := zapcore.AddSync(lumberJackLogger)
	consoleWriteSyncer := zapcore.Lock(os.Stdout)
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer)
	return multiWriteSyncer
	//return zapcore.AddSync(lumberJackLogger)
}
