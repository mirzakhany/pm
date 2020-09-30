package kv

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/mirzakhany/pm/pkg/config"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var err error
	redisServer, err = miniredis.Run()
	if err != nil {
		panic(err)
	}

	serverPort, _ := strconv.Atoi(redisServer.Port())
	host = config.RegisterStringMock("redis.host", redisServer.Host())
	port = config.RegisterIntMock("redis.port", serverPort)

	code := m.Run()
	redisServer.Close()
	os.Exit(code)
}

func TestInit(t *testing.T) {
	ctx := context.Background()
	client, err := Init(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestConn(t *testing.T) {
	ctx := context.Background()
	client, err := Init(ctx)
	assert.Nil(t, err)

	c := client.Conn()
	assert.IsType(t, &redis.Conn{}, c)
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	client, err := Init(ctx)
	assert.Nil(t, err)

	c := Get()
	assert.IsType(t, &Client{}, c)
	assert.Equal(t, client, c)
}

func TestWith(t *testing.T) {

	type testCtx int

	const testK testCtx = 0

	ctx := context.Background()
	client, err := Init(ctx)
	assert.Nil(t, err)

	ctx = context.WithValue(ctx, testK, "test-value")

	client1 := client.With(ctx)

	v := client1.Context().Value(testK)
	assert.NotNil(t, v)
	assert.Equal(t, v, "test-value")
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	client, err := Init(ctx)
	assert.Nil(t, err)

	err = client.Set("test", "test-value", time.Minute*1)
	assert.Nil(t, err)

	testVal, err := client.client.Get(ctx, "test").Result()
	assert.Nil(t, err)
	assert.Equal(t, testVal, "test-value")
}

func TestGetString(t *testing.T) {
	ctx := context.Background()
	client, err := Init(ctx)
	assert.Nil(t, err)

	err = client.Set("test", "test-value", time.Minute*1)
	assert.Nil(t, err)

	testVal, err := client.GetString("test")
	assert.Nil(t, err)
	assert.Equal(t, testVal, "test-value")
}
