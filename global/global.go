package global

import "go.uber.org/zap"

var (
	// 定义全局日志变量
	Logger    *zap.SugaredLogger
	StartTime string
)
