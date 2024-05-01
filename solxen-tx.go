package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"solxen-tx/internal/config"
	"solxen-tx/internal/handler"
	"solxen-tx/internal/svc"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "solxen-tx.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	// 配置
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// 注册job
	group := service.NewServiceGroup()
	handler.RegisterJob(ctx, group)

	// 捕捉信号
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-ch
		logx.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Printf("stop group")
			group.Stop()
			logx.Info("job exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
