package redis

type Store interface {
	Ping() error
}