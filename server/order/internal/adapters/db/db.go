package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jumaniyozov/ordc/order/internal/application/core/domain"
	"time"
)

type Order struct {
	ID         uint
	CustomerID int64
	Status     string
	OrderItems []OrderItem
	CreatedAt  time.Time
}

type OrderItem struct {
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type Adapter struct {
	db *sql.DB
}

func (a *Adapter) Get(id string) (domain.Order, error) {
	var orderEntity Order

	qry := `SELECT id, customer_id, status, order_items, created_at FROM orders WHERE id = $1`
	res := a.db.QueryRow(qry, id)
	err := res.Scan(&orderEntity.ID, &orderEntity.CustomerID, &orderEntity.Status, &orderEntity.OrderItems)
	if err != nil {
		return domain.Order{}, err
	}

	var orderItems []domain.OrderItem
	for _, orderItem := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.UnixNano(),
	}

	return order, nil
}

func (a *Adapter) Save(order *domain.Order) error {
	var orderItems []OrderItem
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	orderModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItems,
		CreatedAt:  time.Unix(0, order.CreatedAt),
	}

	qry := `INSERT INTO orders(customer_id, status, order_items, created_at) VALUES($1, $2, $3, $4)`
	_, err := a.db.Exec(qry, orderModel.CustomerID, orderModel.Status, orderModel.OrderItems, orderModel.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, err := sql.Open("pgx", dataSourceUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Adapter{db: db}, nil
}
