package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/alfiankan/haioo-shoping-cart/application/cart"
	cartGrpcDelivery "github.com/alfiankan/haioo-shoping-cart/application/cart/delivery/rpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	cart_grpc_generated "github.com/alfiankan/haioo-shoping-cart/application/cart/delivery/rpc/codegen"
	cartRepositories "github.com/alfiankan/haioo-shoping-cart/application/cart/repositories"
	cartUseCases "github.com/alfiankan/haioo-shoping-cart/application/cart/usecases"

	"google.golang.org/grpc"

	"github.com/alfiankan/haioo-shoping-cart/common"
	"github.com/alfiankan/haioo-shoping-cart/config"
	"github.com/alfiankan/haioo-shoping-cart/infrastructure"
	"github.com/go-redis/redis/v9"
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
func initArticleApplication(cfg config.ApplicationConfig) cart.ICartUseCase {
	log.Println("cart application started")

	pgConn, redisConn := initInfrastructure(cfg)

	// repositories
	cartPersistRepo := cartRepositories.NewCartRepositoryPostgree(pgConn)
	cartCacheRepo := cartRepositories.NewCartRepositoryCacheRedis(redisConn, 5*time.Minute)

	// usecases
	return cartUseCases.NewCartApplication(cartPersistRepo, cartCacheRepo)

}

func InitgRPCServer(servers ...any) {

	fmt.Println(servers...)
}

type GrpcServerHandler struct {
	Registrar any
	Srv       any
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

	cartApplication := initArticleApplication(cfg)

	grpcServer := grpc.NewServer()

	cart_grpc_generated.RegisterCartServiceServer(grpcServer, cartGrpcDelivery.NewCartRpc(cartApplication))

	listener, err := net.Listen("tcp", ":5300")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		mux := runtime.NewServeMux()
		err := cart_grpc_generated.RegisterCartServiceHandlerFromEndpoint(context.Background(), mux, "localhost:5300", []grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			log.Fatal(err)
		}
		server := http.Server{
			Handler: mux,
		}
		l, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatal(err)
		}
		// start server
		err = server.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("gRPC server starting")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
