package cart

import (
	"time"

	"github.com/google/uuid"
)

type ItemFilter struct {
	Name string
	Qty  int
}

type Cart struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Items     []CartItem `json:"items,omitempty"`
}

type CartItem struct {
	ItemID      uuid.UUID
	ProductCode string
	ProductName string
	Quantity    int
}
