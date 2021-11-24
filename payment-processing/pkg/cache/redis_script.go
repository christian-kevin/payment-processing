package cache

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	ErrClusterNotSupport = errors.New("scripting does not support on this cluster library")
	ErrLimitExceeded     = errors.New("limit exceeded")
)

var decrWithLimitScript = redis.NewScript(1, `
	local key = KEYS[1]
	local limit = tonumber(ARGV[1])
	local value = tonumber(ARGV[2])
	local expireSecond = tonumber(ARGV[3])
	local cnt = tonumber(redis.call('get', key)) or limit
	if (cnt < value ) then
		return -165535
	else
		cnt = cnt - value
		redis.call('set', key, cnt)
		redis.call('expire', key, expireSecond)
		return cnt
	end
`)

func (r *redigoImpl) DecrX(key string, limit int64, value int64, ttl time.Duration) (reply int64, err error) {
	if value > limit {
		return 0, ErrLimitExceeded
	}

	if r.Pool == nil {
		return 0, ErrClusterNotSupport
	}

	conn := r.GetConn()
	defer conn.Close()

	reply, err = redis.Int64(decrWithLimitScript.Do(conn, key, limit, value, ttl.Seconds()))
	if err != nil {
		return 0, err
	}

	if reply == -165535 {
		return 0, ErrLimitExceeded
	}

	return
}
