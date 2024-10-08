package ports

import (
	"github.com/jumaniyozov/ordc/order/internal/application/core/domain"
)

type APIPort interface {
	PlaceOrder(order domain.Order) (domain.Order, error)
}
