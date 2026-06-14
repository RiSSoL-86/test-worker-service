package orders

import "app/src/services/brokers"

type OrderNotificationType string

const (
	OrderNotificationCreate OrderNotificationType = "order.create"
	OrderNotificationUpdate OrderNotificationType = "order.update"
	OrderNotificationDelete OrderNotificationType = "order.delete"
)

const (
	topic         = "orders"
	consumerGroup = "orders-worker"
)

func InitConsumer(runner *brokers.Runner, dependencies *Dependencies) {
	runner.Subscribe(topic, consumerGroup).
		On(string(OrderNotificationCreate), dependencies.Handler.Create).
		On(string(OrderNotificationUpdate), dependencies.Handler.Update).
		On(string(OrderNotificationDelete), dependencies.Handler.Delete)
}
