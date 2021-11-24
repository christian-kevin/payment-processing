package redis

import (
	"spenmo/payment-processing/payment-processing/pkg/cache"
)

type storeImpl struct {
	client cache.Cache
}

func NewStore(cache cache.Cache) Store {
	return &storeImpl{client: cache}
}

func (s *storeImpl) Ping() error {
	_, err := s.client.Get("")
	if err == cache.ErrNil {
		return nil
	}
	return err
}
