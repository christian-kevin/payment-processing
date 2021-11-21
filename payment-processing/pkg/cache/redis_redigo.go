package cache

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

type redigoImpl struct {
	Pool *redis.Pool
}

const (
	redisSet              = "SET"
	redisGet              = "GET"
	redisHGet             = "HGET"
	redisHMGet            = "HMGET"
	redisHSet             = "HSET"
	redisHSetNX           = "HSETNX"
	redisHMSet            = "HMSET"
	redisHKeys            = "HKEYS"
	redisHVals            = "HVALS"
	redisHGetAll          = "HGETALL"
	redisHLen             = "HLEN"
	redisScan             = "SCAN"
	redisMatch            = "MATCH"
	redisIncrBy           = "INCRBY"
	redisDel              = "DEL"
	redisHDel             = "HDEL"
	redisEx               = "EX"
	redisNX               = "NX"
	redisXX               = "XX"
	redisLimit            = "LIMIT"
	redisPing             = "PING"
	redisExpire           = "EXPIRE"
	redisINCR             = "INCR"
	redisZAdd             = "ZADD"
	redisZRange           = "ZRANGE"
	redisZRevRange        = "ZREVRANGE"
	redisZRangeByScore    = "ZRANGEBYSCORE"
	redisZRevRangeByScore = "ZREVRANGEBYSCORE"
	redisZRank            = "ZRANK"
	redisZRevRank         = "ZREVRANK"
	redisZScore           = "ZSCORE"
	redisZCount           = "ZCOUNT"
	redisZRemRangeByScore = "ZREMRANGEBYSCORE"
	redisZREM             = "ZREM"
	redisHExists          = "HEXISTS"
	redisHIncrBy          = "HINCRBY"
	redisExists           = "EXISTS"
)

var (
	ErrNX  = errors.New("key already exist")
	ErrXX  = errors.New("key is not exists")
	ErrHNX = errors.New("(key,field) combination already exists")
)

func newRedigo(cfg *Config) (*redigoImpl, error) {
	var dial []redis.DialOption

	if cfg.ReadTimeout != 0 {
		dial = append(dial, redis.DialReadTimeout(cfg.ReadTimeout))
	}

	if cfg.ReadTimeout != 0 {
		dial = append(dial, redis.DialWriteTimeout(cfg.ReadTimeout))
	}

	if cfg.DialTimeout != 0 {
		dial = append(dial, redis.DialConnectTimeout(cfg.ReadTimeout))
	}

	r := &redigoImpl{
		Pool: &redis.Pool{
			IdleTimeout: cfg.IdleTimeout,
			MaxIdle:     cfg.MaxIdle,
			MaxActive:   cfg.PoolSize,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", cfg.ServerAddr, dial...)
				if err != nil {
					return nil, err
				}

				return c, nil
			},
		},
	}
	_, err := r.Pool.Get().Do(redisPing)

	return r, err
}

func (r *redigoImpl) ErrorOnHashCacheMiss() error {
	return ErrNil
}

func (r *redigoImpl) ErrorOnCacheMiss() error {
	return ErrNil
}

func (r *redigoImpl) GetConn() Conn {
	return r.Pool.Get()
}

func (r *redigoImpl) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := r.GetConn()
	defer c.Close()
	return c.Do(commandName, args...)
}

func (r *redigoImpl) Set(key, value string, ttl time.Duration) error {
	var err error

	if 0 >= ttl {
		_, err = r.Do(redisSet, key, value)
	} else {
		_, err = r.Do(redisSet, key, value, redisEx, int64(ttl.Seconds()))
	}

	return err
}

func (r *redigoImpl) SetNX(key, value string, ttl time.Duration) error {
	var (
		err   error
		reply interface{}
	)

	if 0 >= ttl {
		reply, err = r.Do(redisSet, key, value, redisNX)
	} else {
		reply, err = r.Do(redisSet, key, value, redisEx, int64(ttl.Seconds()), redisNX)
	}

	if nil != err {
		return err
	}
	if nil == reply {
		return ErrNX
	}

	return err
}

func (r *redigoImpl) HSet(key, field, value string, ttl time.Duration) error {
	var err error

	_, err = r.Do(redisHSet, key, field, value)
	if err != nil {
		return err
	}
	if ttl > 0 {
		_, err = r.Do(redisExpire, key, int64(ttl.Seconds()))
	}

	return err
}

func (r *redigoImpl) HSetNX(key, field, value string, ttl time.Duration) error {
	var err error

	reply, err := redis.Int64(r.Do(redisHSetNX, key, field, value))
	if err != nil {
		return err
	}
	if reply == 0 {
		return ErrHNX
	}

	if ttl > 0 {
		_, err = r.Do(redisExpire, key, int64(ttl.Seconds()))
	}

	return err
}

func (r *redigoImpl) HMSet(key string, fieldsMap map[string]string, ttl time.Duration) error {
	_, err := r.Do(redisHMSet, redis.Args{}.Add(key).AddFlat(fieldsMap)...)
	if err != nil {
		return err
	}

	if ttl > 0 {
		_, err = r.Do(redisExpire, key, int64(ttl.Seconds()))
	}

	return err
}

func (r *redigoImpl) HGet(key, field string) (string, error) {
	reply, err := redis.String(r.Do(redisHGet, key, field))
	if err == redis.ErrNil {
		return reply, ErrNil
	}
	return reply, err
}

func (r *redigoImpl) HMGet(key string, fields ...string) ([]string, error) {
	reply, err := r.Do(redisHMGet, redis.Args{}.Add(key).AddFlat(fields)...)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) HDel(key string, fields ...string) (int64, error) {
	return redis.Int64(r.Do(redisHDel, redis.Args{}.Add(key).AddFlat(fields)...))
}

func (r *redigoImpl) HKeys(key string) ([]string, error) {
	reply, err := r.Do(redisHKeys, key)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) HVals(key string) ([]string, error) {
	reply, err := r.Do(redisHVals, key)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) HGetAll(key string) (map[string]string, error) {
	reply, err := r.Do(redisHGetAll, key)
	return redis.StringMap(reply, err)
}

func (r *redigoImpl) HLen(key, field string) (int64, error) {
	var (
		reply interface{}
		err   error
	)

	if "" == field {
		reply, err = r.Do(redisHLen, key)
	} else {
		reply, err = r.Do(redisHLen, key, field)
	}

	return redis.Int64(reply, err)
}

func (r *redigoImpl) Get(key string) (string, error) {
	reply, err := redis.String(r.Do(redisGet, key))
	if err == redis.ErrNil {
		return reply, ErrNil
	}
	return reply, err
}

func (r *redigoImpl) IncrBy(key string, incr int64) (int64, error) {
	reply, err := r.Do(redisIncrBy, key, incr)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) Del(key ...string) error {
	if len(key) == 0 {
		return ErrInsufficientArgument
	}
	_, err := r.Do(redisDel, redis.Args{}.AddFlat(key)...)
	return err
}

func (r *redigoImpl) ZAddXX(key, member string, score int) error {
	_, err := r.Do(redisZAdd, key, redisXX, score, member)
	return err
}

func (r *redigoImpl) ZAddNX(key, member string, score int64) (int64, error) {
	return redis.Int64(r.Do(redisZAdd, key, redisNX, score, member))
}

func (r *redigoImpl) ZAddINCR(key, member string, score int) error {
	_, err := r.Do(redisZAdd, key, redisINCR, score, member)
	return err
}

func (r *redigoImpl) ZRange(key string, start, stop int) ([]string, error) {
	reply, err := r.Do(redisZRange, key, start, stop)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) ZRevRange(key string, start, stop int) ([]string, error) {
	reply, err := r.Do(redisZRevRange, key, start, stop)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) ZRangeByScore(key string, min, max, offset, count int) ([]string, error) {
	reply, err := r.Do(redisZRangeByScore, key, min, max, redisLimit, offset, count)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) ZRevRangeByScore(key string, max, min, offset, count int) ([]string, error) {
	reply, err := r.Do(redisZRevRangeByScore, key, max, min, redisLimit, offset, count)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) ZRank(key, member string) (int64, error) {
	reply, err := r.Do(redisZRank, key, member)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) ZRevRank(key, member string) (int64, error) {
	reply, err := r.Do(redisZRevRank, key, member)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) ZScore(key, member string) (int64, error) {
	reply, err := r.Do(redisZScore, key, member)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) ZCount(key string, min, max int) (int64, error) {
	reply, err := r.Do(redisZCount, key, min, max)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) ZRemRangeByScore(key string, start, stop int) (int64, error) {
	return redis.Int64(r.Do(redisZRemRangeByScore, key, start, stop))
}

func (r *redigoImpl) ZRem(key string, members ...string) (int64, error) {
	return redis.Int64(r.Do(redisZREM, redis.Args{}.Add(key).AddFlat(members)...))
}

func (r *redigoImpl) ZAddXXIncrBy(key, member string, incrValue int64) (int64, error) {
	reply, err := redis.Int64(r.Do(redisZAdd, key, redisXX, redisINCR, incrValue, member))
	if err == redis.ErrNil {
		return 0, ErrXX
	}
	return reply, err
}

func (r *redigoImpl) HExists(key, field string) (bool, error) {
	return redis.Bool(r.Do(redisHExists, key, field))
}

func (r *redigoImpl) HIncrBy(key, field string, incrValue int64) (int64, error) {
	return redis.Int64(r.Do(redisHIncrBy, key, field, incrValue))
}

func (r *redigoImpl) Expire(key string, ttl time.Duration) (int64, error) {
	return redis.Int64(r.Do(redisExpire, key, int64(ttl.Seconds())))
}

func (r *redigoImpl) Exists(key string) (bool, error) {
	return redis.Bool(r.Do(redisExists, key))
}