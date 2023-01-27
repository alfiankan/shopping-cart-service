package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alfiankan/haioo-shoping-cart/common"
	"github.com/alfiankan/haioo-shoping-cart/config"
	"github.com/alfiankan/haioo-shoping-cart/infrastructure"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {

	// set override config
	envs := map[string]string{
		"PG_DATABASE_HOST":     "127.0.0.1",
		"PG_DATABASE_USERNAME": "postgres",
		"PG_DATABASE_PASSWORD": "postgres",
		"PG_DATABASE_NAME":     "postgres",
		"PG_DATABASE_PORT":     "2345",
		"PG_DATABASE_SSL_MODE": "disable",
		"LOG_LEVEL":            "debug",
		"REDIS_HOST":           "127.0.0.1:9376",
		"REDIS_PASSWORD":       "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		"HTTP_API_PORT":        "5000",
	}

	for key, val := range envs {
		os.Setenv(key, val)
	}
	cfg := config.Load("void")

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// POSTGREESQL SETUP
	postgreePorts := []docker.PortBinding{{HostPort: strconv.Itoa(cfg.PostgreePort)}}
	pool.RemoveContainerByName("haioo-cart-postgree-test")

	if _, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "haioo-cart-postgree-test",
		Repository:   "postgres",
		Tag:          "14.1-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{"5432/tcp": postgreePorts},
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", cfg.PostgreeUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.PostgreePass),
			"listen_addresses = '*'",
		}}); err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// REDIS SETUP
	parsedEnv := strings.Split(cfg.RedisHost[0], ":")
	redisPorts := []docker.PortBinding{{HostPort: parsedEnv[len(parsedEnv)-1]}}
	pool.RemoveContainerByName("haioo-cart-redis-test")

	if _, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "haioo-cart-redis-test",
		Repository:   "redis",
		Tag:          "6.2-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{"6379/tcp": redisPorts},
		Cmd:          []string{"--requirepass", cfg.RedisPass},
	}); err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// PING ALL CHECK CONNECTION OK
	for {
		log.Println("try to ping redis ⏳")
		redisConn, _ := infrastructure.NewRedisConnection(cfg)
		if err := redisConn.Ping(context.Background()).Err(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	log.Println("redis up and running ✅")

	for {
		log.Println("try to ping postgree ⏳")
		pgConn, _ := infrastructure.NewPgConnection(cfg)
		if err := pgConn.Ping(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	log.Println("postgree up and running ✅")

	// SETUP MIGRATION AND SEED
	if err := common.Migration("../../../"); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	pool.RemoveContainerByName("haioo-cart-redis-test")
	pool.RemoveContainerByName("haioo-cart-elasticsearch-test")
	pool.RemoveContainerByName("haioo-cart-postgree-test")
	os.Exit(code)

}
