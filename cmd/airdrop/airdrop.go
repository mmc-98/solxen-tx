package airdrop

import (
	"context"
	"flag"
	"solxen-tx/internal/config"
	"solxen-tx/internal/logic"
	"solxen-tx/internal/svc"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var ConfigFile = flag.String("f3", "solxen-tx.yaml", "the config file")

func Airdrop(_ *cobra.Command, _ []string) error {

	flag.Parse()

	var c config.Config
	conf.MustLoad(*ConfigFile, &c)
	c.Sol.Url = "http://69.10.34.226:8899"
	logx.MustSetup(c.LogConf)
	logx.DisableStat()

	ser := svc.NewServiceContext(c)

	l := logic.NewProducerLogic(context.Background(), ser)

	l.Airdrop()

	return nil

}
