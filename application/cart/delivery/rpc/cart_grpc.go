package cart_grpc

import (
	"context"
	"log"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	cart_grpc_generated "github.com/alfiankan/haioo-shoping-cart/application/cart/delivery/rpc/codegen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartRpc struct {
	cartUseCase cart.ICartUseCase
	cart_grpc_generated.UnimplementedCartServiceServer
}

func NewCartRpc(cartUseCase cart.ICartUseCase) cart_grpc_generated.CartServiceServer {
	return &CartRpc{
		cartUseCase: cartUseCase,
	}
}

func (h *CartRpc) CreateCart(ctx context.Context, _ *cart_grpc_generated.EmptyRequest) (response *cart_grpc_generated.EmptyDataResponse, err error) {

	log.Println("Invoked")
	if err = h.cartUseCase.CreateCart(ctx); err != nil {
		err = status.Error(codes.Internal, err.Error())
		return
	}
	response = &cart_grpc_generated.EmptyDataResponse{
		Success: true,
		Code:    0,
		Msg:     "Ok",
	}
	return
}
