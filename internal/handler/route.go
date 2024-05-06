package handler

import (
	"context"
	"solxen-tx/internal/logic"
	"solxen-tx/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
)

func RegisterJob(serverCtx *svc.ServiceContext, group *service.ServiceGroup) {

	group.Add(logic.NewProducerLogic(context.Background(), serverCtx))

	group.Start()

}
