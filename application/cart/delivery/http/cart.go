package http_delivery

import (
	"net/http"
	"strconv"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	"github.com/alfiankan/haioo-shoping-cart/application/cart/transport"
	"github.com/alfiankan/haioo-shoping-cart/common"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CartHttpApi struct {
	cartUseCase cart.ICartUseCase
}

func NewCartHttpApi(cartUseCase cart.ICartUseCase) *CartHttpApi {
	return &CartHttpApi{cartUseCase}
}

// @Description create new cart/bucket
// @Tags cart
// @Accept json
// @Produce json
// @Success 200
// @Router /carts [post]
func (handler *CartHttpApi) NewCart(c echo.Context) error {

	err := handler.cartUseCase.CreateCart(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &transport.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    common.EmptyResponseData,
		})
	}
	return c.JSON(http.StatusOK, &transport.BaseResponse{
		Message: common.HttpSuccessCreated,
		Data:    common.EmptyResponseData,
	})

}

// @Description get carts bucket
// @Tags cart
// @Accept json
// @Produce json
// @Success 200
// @Router /carts [get]
func (handler *CartHttpApi) GetAllCarts(c echo.Context) error {

	var carts []cart.Cart
	carts, err := handler.cartUseCase.GetCarts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &transport.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	if len(carts) == 0 {
		carts = []cart.Cart{}
	}

	return c.JSON(http.StatusOK, &transport.BaseResponse{
		Message: common.HttpSuccess,
		Data:    carts,
	})
}

// @Description Add product to cart
// @Tags cart
// @Param cart_id path string true "cart_id uuid"
// @Param article body transport.ItemRequest  true "Article detail"
// @Accept json
// @Produce json
// @Success 200
// @Router /carts/{cart_id}/items [post]
func (handler *CartHttpApi) AddProductToCart(c echo.Context) error {

	cartID, err := uuid.Parse(c.Param("cart_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &transport.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	var reqBody transport.ItemRequest
	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &transport.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	if err := handler.cartUseCase.AddToCart(c.Request().Context(), cart.Cart{
		ID: cartID,
		Items: []cart.CartItem{{
			ProductCode: reqBody.ProductCode,
			ProductName: reqBody.ProductName,
			Quantity:    reqBody.Quantity,
		}},
	}); err != nil {
		common.Log(common.LOG_LEVEL_ERROR, err.Error())

		return c.JSON(http.StatusInternalServerError, &transport.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	return c.JSON(http.StatusOK, &transport.BaseResponse{
		Message: common.HttpSuccess,
		Data:    common.EmptyResponseData,
	})

}

// @Description delete item/produk from cart
// @Tags cart
// @Param cart_id path string true "cart_id uuid"
// @Param product_code path string true "kodeProduk uuid"
// @Accept json
// @Produce json
// @Success 200
// @Router /carts [delete]
func (handler *CartHttpApi) DeleteCartItem(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

// @Description get all product/item from cart
// @Tags cart
// @Param cart_id path string true "cart_id uuid"
// @Param namaProduk query string false "filter by nama produk"
// @Param kuantitas query int false "filter by kuantitas"
// @Accept json
// @Produce json
// @Success 200
// @Router /carts/{cart_id}/items [get]
func (handler *CartHttpApi) GetAllCartitems(c echo.Context) error {

	cartID, err := uuid.Parse(c.Param("cart_id"))
	filterProdukName := c.QueryParam("namaProduk")

	kuantitasParamStr := "0"
	if c.QueryParam("kuantitas") != "" {
		kuantitasParamStr = c.QueryParam("kuantitas")
	}

	filterKuantitas, err := strconv.Atoi(kuantitasParamStr)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &transport.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	cartRes, err := handler.cartUseCase.GetCartItems(
		c.Request().Context(),
		cartID,
		cart.ItemFilter{Name: filterProdukName, Qty: filterKuantitas},
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &transport.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	if len(cartRes.Items) == 0 {
		cartRes.Items = []cart.CartItem{}
	}

	return c.JSON(http.StatusOK, &transport.BaseResponse{
		Message: common.HttpSuccess,
		Data:    cartRes.Items,
	})

}

func (handler *CartHttpApi) HandleRoute(e *echo.Echo) {

	e.POST("/carts", handler.NewCart)
	e.GET("/carts", handler.GetAllCarts)

	// kodeProduk, namaProduk, kuantitas
	e.POST("/carts/:cart_id/items", handler.AddProductToCart)

	// delete produk by kodeproduk
	e.DELETE("/carts/:cart_id/:product_code", handler.DeleteCartItem)

	//{kodeProduk}- {namaProduk} - ({kuantitas }) -> filter by namaProduk and kuantitas
	e.GET("/carts/:cart_id/items", handler.GetAllCartitems)

}
