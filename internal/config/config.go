package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type Sol struct {
	Url string
	// TxnUrl    string
	// GrpcUrl   string
	Mnemonic   string
	WalletPath string
	Num        int
	Fee        uint64
	ToAddr     string
	Time       int
	// ProgramIdMiner string
	ProgramId string
	HdPath    string
}
type Config struct {
	service.ServiceConf
	// Redis redis.RedisConf
	//
	// DB struct {
	// 	DataSource string
	// }
	// Cache cache.CacheConf
	//
	// // KqPusherConf struct {
	// // 	Brokers []string
	// // 	Topic   string
	// // }
	// DqConf dq.DqConf
	LogConf logx.LogConf
	Sol     Sol

	// Vault struct {
	// 	Address *vault.Config
	// 	Token   string
	// }
}
