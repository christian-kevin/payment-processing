package cache

//go:generate mockgen -destination mock_cache/mock_cache.go . Cache,Cacher,Conn,HashCacher,MultiCacher,Scripter

import (
	"errors"
	"time"
)

// Cacher is an interface for basic caching operations.
type Cacher interface {
	GetConn() Conn
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Del(key ...string) error
	ErrorOnCacheMiss() error
}

type Conn interface {
	Close() error

	Err() error

	Do(commandName string, args ...interface{}) (reply interface{}, err error)

	Send(commandName string, args ...interface{}) error

	Flush() error

	Receive() (reply interface{}, err error)
}

// HashCacher is an interface for redis hash operations
type HashCacher interface {
	// HSet Sets field in the hash stored at key to value. If key does not exist, a new key holding a hash is created. If field already exists in the hash, it is overwritten.
	HSet(key, field, value string, ttl time.Duration) error
	// HSetNX Sets field in the hash stored at key to value, only if field does not yet exist.
	HSetNX(key, field, value string, ttl time.Duration) error
	// HMSet Sets multiple pair of field & value
	HMSet(key string, fieldsMap map[string]string, ttl time.Duration) error
	// HGet Returns the value associated with field in the hash stored at key.
	HGet(key, field string) (string, error)
	// HMGet Returns values with multiple fields
	HMGet(key string, fields ...string) ([]string, error)
	// HDel Removes the specified fields from the hash stored at key. Specified fields that do not exist within this hash are ignored. If key does not exist, it is treated as an empty hash and this command returns 0.
	HDel(key string, fields ...string) (int64, error)
	// HKeys Returns all field names in the hash stored at key.
	HKeys(key string) ([]string, error)
	// HVals Returns all values in the hash stored at key.
	HVals(key string) ([]string, error)
	// HGetAll Returns all fields and values of the hash stored at key. In the returned value, every field name is followed by its value, so the length of the reply is twice the size of the hash.
	HGetAll(key string) (map[string]string, error)
	HExists(key, field string) (bool, error)
	HIncrBy(key, field string, incrValue int64) (int64, error)
	ErrorOnHashCacheMiss() error
}

type Cache interface {
	Cacher
	HashCacher
	Scripter
	// SetNX et key to hold string value if key does not exist
	SetNX(key, value string, ttl time.Duration) error
	// ZAddXX add member to sorted set. Score is updated with new value. Never add score.
	ZAddXX(key, member string, score int) error
	// ZAddNX add member to sorted set. Score is added with new value. Never update score.
	ZAddNX(key, member string, score int64) (int64, error)
	// ZAddINCR add member to sorted set. Score is incremented from previous value if exist
	ZAddINCR(key, member string, score int) error
	// ZRange return members with its score from sorted set. Sorted ascending
	ZRange(key string, start, stop int) ([]string, error)
	// ZRevRange return members with its score from sorted set. Sorted descending
	ZRevRange(key string, start, stop int) ([]string, error)
	// ZRangeByScore return members with its score from sorted set, using score as range. Sorted ascending.
	ZRangeByScore(key string, min, max, offset, count int) ([]string, error)
	// ZRevRangeByScore return members with its score from sorted set, using score as range. Sorted ascending.
	ZRevRangeByScore(key string, max, min, offset, count int) ([]string, error)
	// ZRank return member rank in sorted set. Sorted ascending.
	ZRank(key, member string) (int64, error)
	// ZRevRank return member rank in sorted set. Sorted descending.
	ZRevRank(key, member string) (int64, error)
	// ZScore return members score
	ZScore(key, member string) (int64, error)
	// ZCount return count of members which score is between min and max.
	ZCount(key string, min, max int) (int64, error)
	// ZRemRangeByScore Delete member by score range
	ZRemRangeByScore(key string, start, stop int) (int64, error)
	ZRem(key string, members ...string) (int64, error)
	Expire(key string, ttl time.Duration) (int64, error)
	Exists(key string) (bool, error)
}

// Scripter is interface contract for redis scripting
type Scripter interface {
	DecrX(key string, limit int64, value int64, ttl time.Duration) (int64, error)
}

type HMulti struct {
	Field    string
	Value    int64
	MaxValue int64
}

type Implementation int

const (
	Redis = Implementation(iota)
	RedisCluster
	RedisSentinel
)

type Config struct {
	ServerAddr   string
	MaxIdle      int
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	DialTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
}

// https://godoc.org/github.com/go-redis/redis#ClusterOptions
type ConfigCacheCluster struct {
	Addrs          []string
	MaxRedirects   int
	ReadOnly       bool
	RouteByLatency bool
	RouteRandomly  bool
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	PoolSize       int
	DialTimeout    time.Duration
	IdleTimeout    time.Duration
}

var (
	ErrNil                  = errors.New("redis status: nil")
	ErrInsufficientArgument = errors.New("wrong number of arguments")
)

// New return ready to use Cache instance
func New(impl Implementation, cfg interface{}) (Cache, error) {
	switch impl {
	case Redis:
		return newRedigo(cfg.(*Config))
	case RedisCluster:
		return nil, nil
	case RedisSentinel:
		return nil,nil
	}

	return nil, errors.New("no cache implementations found")
}
