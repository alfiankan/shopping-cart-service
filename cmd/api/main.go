package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	cartDeliveryHttp "github.com/alfiankan/haioo-shoping-cart/application/cart/delivery/http"
	cartRepositories "github.com/alfiankan/haioo-shoping-cart/application/cart/repositories"
	cartUseCases "github.com/alfiankan/haioo-shoping-cart/application/cart/usecases"

	middlewares "github.com/alfiankan/haioo-shoping-cart/common/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/alfiankan/haioo-shoping-cart/common"
	"github.com/alfiankan/haioo-shoping-cart/config"
	"github.com/alfiankan/haioo-shoping-cart/infrastructure"
	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
)

// initInfrastructure init all infrastructure needs to run this application
func initInfrastructure(cfg config.ApplicationConfig) (pgConn *sql.DB, redisConn redis.UniversalClient) {

	pgConn, err := infrastructure.NewPgConnection(cfg)
	common.LogExit(err, common.LOG_LEVEL_ERROR)

	redisConn, err = infrastructure.NewRedisConnection(cfg)
	common.LogExit(err, common.LOG_LEVEL_ERROR)

	return
}

// initArticleApplication init app by injecting deps
func initArticleApplication(httpServer *echo.Echo, cfg config.ApplicationConfig) {

	pgConn, redisConn := initInfrastructure(cfg)

	// repositories
	cartPersistRepo := cartRepositories.NewCartRepositoryPostgree(pgConn)
	cartCacheRepo := cartRepositories.NewCartRepositoryCacheRedis(redisConn, 5*time.Minute)

	// usecases
	cartUseCases := cartUseCases.NewCartApplication(cartPersistRepo, cartCacheRepo)

	// handle http request response
	cartDeliveryHttp.NewCartHttpApi(cartUseCases).HandleRoute(httpServer)
}

// @title haioo-shopping-cart-api
// @version 1.0
// @description Go implemented api.
// @contact.name alfiankan
// @contact.url https://github.com/alfiankan
// @contact.email alfiankan19@gmail.com
// @license.name Apache 2.0
// @BasePath /
func main() {

	cfg := config.Load()
	e := echo.New()
	e.Use(middlewares.MiddlewaresRegistry...)

	initArticleApplication(e, cfg)

	// swagger api docs
	url := echoSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger.yaml", cfg.HTTPApiPort))
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	e.Static("/docs", "docs")

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPApiPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server", err.Error())
		}
	}()

	// graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 60 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
