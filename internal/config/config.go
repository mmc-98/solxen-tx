package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type Sol struct {
	Url       string
	Key       string
	Num       int
	Fee       uint64
	ToAddr    string
	Time      int
	ProgramID string
	HdPath    string
}
type Config struct {
	service.ServiceConf
	LogConf logx.LogConf
	Sol     Sol
}
