package balance

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

var ConfigFile = flag.String("f4", "solxen-tx.yaml", "the config file")

func Balance(_ *cobra.Command, _ []string) error {

	flag.Parse()

	var c config.Config
	conf.MustLoad(*ConfigFile, &c)
	logx.MustSetup(c.LogConf)
	logx.DisableStat()

	ser := svc.NewServiceContext(c)

	l := logic.NewProducerLogic(context.Background(), ser)

	l.Balance()

	return nil

}
