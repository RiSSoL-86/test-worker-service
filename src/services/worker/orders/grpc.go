package orders

import (
	orderspb "app/src/core/proto/orders"
	"app/src/services/grpc"
)

func InitGrpc(server *grpc.Server, dependencies *Dependencies) {
	orderspb.RegisterOrdersServiceServer(server.GRPC(), dependencies.Handler)
}
