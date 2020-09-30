package kv

import (
	"context"
	"strconv"

	"github.com/alicebob/miniredis/v2"
	"github.com/mirzakhany/pm/pkg/config"
)

var redisServer *miniredis.Miniredis

func InitMock(ctx context.Context) (*Client, error) {
	var err error
	redisServer, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	serverPort, _ := strconv.Atoi(redisServer.Port())
	host = config.RegisterStringMock("redis.host", redisServer.Host())
	port = config.RegisterIntMock("redis.port", serverPort)
	return Init(ctx)
}
