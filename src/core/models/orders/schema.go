package orders

import "time"

type Order struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey"`
	Title         string    `gorm:"column:title;type:text;not null"`
	Description   string    `gorm:"column:description;type:text;not null"`
	CustomerName  string    `gorm:"column:customer_name;type:text;not null"`
	CustomerPhone string    `gorm:"column:customer_phone;type:text;not null"`
	Address       string    `gorm:"column:address;type:text;not null"`
	Priority      string    `gorm:"column:priority;type:text;not null"`
	Status        string    `gorm:"column:status;type:text;not null;default:new"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (Order) TableName() string {
	return "orders"
}
