package session

import (
	"context"
	"github.com/mirzakhany/pm/pkg/kv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	_, err := kv.InitMock(ctx)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestSession(t *testing.T) {

	type test struct {
		Username string
		Password string
	}

	test1 := test{
		Username: "test1",
		Password: "test1",
	}

	err := Set("test1", test1, time.Minute*1)
	assert.Nil(t, err)

	var test1Val test
	err = Get("test1", &test1Val)
	assert.Nil(t, err)
	assert.Equal(t, test1, test1Val)
}
