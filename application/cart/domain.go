package cart

import "github.com/google/uuid"

type ItemFilter struct {
	Name string
	Qty  string
}

type Cart struct {
	ID    uuid.UUID
	Items []CartItem
}

type CartItem struct {
	ProductCode uuid.UUID
	ProductName string
	Quantity    int
}
