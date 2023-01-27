package transport

import "github.com/google/uuid"

type ItemRequest struct {
	ProductCode uuid.UUID `json:"kodeProduk"`
	ProductName string    `json:"namaProduk"`
	Quantity    int       `json:"kuantitas"`
}
