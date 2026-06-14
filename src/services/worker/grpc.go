package worker

import (
	"app/src/services/grpc"
	"app/src/services/worker/orders"
)

func InitGrpc(server *grpc.Server, dependencies *Dependencies) {
	orders.InitGrpc(server, dependencies.Orders)
}
