package worker

import (
	"app/src/services/brokers"
	"app/src/services/worker/orders"
)

func InitConsumers(runner *brokers.Runner, dependencies *Dependencies) {
	orders.InitConsumer(runner, dependencies.Orders)
}
