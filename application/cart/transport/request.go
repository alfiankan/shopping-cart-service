package transport

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ItemRequest struct {
	ProductCode string `json:"kodeProduk"`
	ProductName string `json:"namaProduk"`
	Quantity    int    `json:"kuantitas"`
}

func (request ItemRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ProductCode, validation.Required, is.Alphanumeric),
		validation.Field(&request.ProductName, validation.Required, is.ASCII),
		validation.Field(&request.Quantity, validation.Required, is.Int),
	)
}
