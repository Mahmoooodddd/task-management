package platform

import "go.uber.org/zap"

const (
	LogFilePath     = "./var/logs/task-management.log"
	TestLogFilePath = "./var/logs/task-management.log"
)

type Logger interface {
	Error(msg string, err error, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}
type logger struct {
	logger *zap.Logger
}

func (l *logger) Error(msg string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Error(err)}, fields...)
	l.logger.Error(msg, allFields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func NewLogger(configs Configs) Logger {
	cfg := zap.NewProductionConfig()
	env := configs.GetEnv()
	logPath := LogFilePath

	if env == EnvTest {
		logPath = TestLogFilePath
	}

	cfg.OutputPaths = []string{
		logPath,
	}

	zapLogger, err := cfg.Build()

	if err != nil {
		panic("can not set logger")
	}

	return &logger{
		logger: zapLogger,
	}
}
