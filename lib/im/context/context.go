package context

import (
	"sync"

	"github.com/CanftIn/gothafoss/lib/cache"
	"github.com/CanftIn/gothafoss/lib/common"
	"github.com/CanftIn/gothafoss/lib/im/config"
	"github.com/CanftIn/gothafoss/lib/im/event"
	"github.com/CanftIn/gothafoss/lib/im/http"
	"github.com/CanftIn/gothafoss/lib/im/pool"
	"github.com/CanftIn/gothafoss/lib/log"
	"github.com/RussellLuo/timingwheel"
	"github.com/gocraft/dbr"
)

type Context struct {
	cfg          *config.Config
	mysqlSession *dbr.Session
	redisCache   *common.RedisCache
	memoryCache  cache.Cache
	EventPool    pool.Collector
	PushPool     pool.Collector // 离线push
	Event        event.Event
	timingWheel  *timingwheel.TimingWheel // Time wheel delay task
	httpRouter   *http.IMHttp

	log.Log
	valueMap sync.Map
}

// NewContext NewContext
func NewContext(cfg *config.Config) *Context {
	c := &Context{
		cfg:         cfg,
		Log:         log.NewTLog("Context"),
		EventPool:   pool.StartDispatcher(cfg.EventPoolSize),
		PushPool:    pool.StartDispatcher(cfg.Push.PushPoolSize),
		timingWheel: timingwheel.NewTimingWheel(cfg.TimingWheelTick.Duration, cfg.TimingWheelSize),
		valueMap:    sync.Map{},
	}
	c.timingWheel.Start()
	return c
}

// GetConfig 获取配置信息
func (c *Context) GetConfig() *config.Config {
	return c.cfg
}
