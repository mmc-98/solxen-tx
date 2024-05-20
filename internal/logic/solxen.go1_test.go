package logic

import (
	"context"
	"flag"
	"solxen-tx/internal/config"
	"solxen-tx/internal/svc"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/zeromicro/go-zero/core/conf"
)

var (
	l         = new(Producer)
	programId = solana.MustPublicKeyFromBase58("F8yUTgMN6E96QYhVVY9UVkKzKrDjpHLqrQ7bPCoqaJHz")
)

func TestNewXnRecordAddress(t *testing.T) {

}

func init() {

	var configFile = flag.String("f", "../../build/etc/solxen-tx.dev.yaml", "the config file")

	var c config.Config
	conf.MustLoad(*configFile, &c)

	c.Sol.ProgramId = "F8yUTgMN6E96QYhVVY9UVkKzKrDjpHLqrQ7bPCoqaJHz"

	s := svc.NewServiceContext(c)
	l = NewProducerLogic(context.TODO(), s)
}

func TestMint(t *testing.T) {
	l.Mint()
}

func TestFetchIDL(t *testing.T) {

	// spew.Dump(addr, err)
}
