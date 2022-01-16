package redis_lock

import (
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	network = "tcp"
	port    = "82.157.50.197:6380"
)

func ReturnLock(bookid string) *redsync.Mutex {
	pool := newPool()
	r := redsync.New([]redsync.Pool{pool})
	mutex := r.NewMutex(bookid, redsync.SetExpiry(time.Duration(2)*time.Second),
		redsync.SetRetryDelay(time.Duration(5)*time.Millisecond))
	return mutex
}

func newPool() *redis.Pool {
	return &redis.Pool{

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, port)
			if err != nil {
				panic(err.Error())
				return nil, err
			}
			return c, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return nil
		},
		// MaxIdle :最大的空闲池
		MaxIdle: 3,
		// IdleTimeout : 如果连接上了空闲了这么就没有动静就会关闭
		IdleTimeout: time.Duration(3) * time.Second,
		// MaxConnLifetime : 最大的连接时间
		MaxConnLifetime: time.Duration(5) * time.Second,
	}
}
