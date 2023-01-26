package cart

import "github.com/google/uuid"

type Cart struct {
	ID    uuid.UUID
	Items []CartItem
}

type CartItem struct {
	ProductCode uuid.UUID
	ProductName string
	Quantity    int
}
