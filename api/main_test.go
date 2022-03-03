package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	db "phantom/db/mongo"
	"phantom/util"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode) // 把gin的运行模式改为测试模式，默认是调式模式，调试模式会给出详细的日志，不利于阅读

	os.Exit(m.Run())
}
