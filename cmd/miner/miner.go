package miner

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

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var ConfigFile = flag.String("f", "solxen-tx.yaml", "the config file")

func Miner(_ *cobra.Command, _ []string) error {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*ConfigFile, &c)
	logx.MustSetup(c.LogConf)
	logx.DisableStat()

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
			return nil
		case syscall.SIGHUP:
		default:
			return nil
		}
	}

	return nil
}
