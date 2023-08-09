package storage

import (
	"context"
	"encoding"
	"strconv"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"
	"github.com/redis/go-redis/v9"
)

type redisStorage struct {
	client *redis.Client
	ctx    context.Context
}

var _ Storage = (*redisStorage)(nil)

func (n *nullInt64) ScanRedis(s string) (err error) {
	return n.Scan(s)
}

var _ redis.Scanner = (*nullInt64)(nil)

func newRedisStorage(dbName string, compressor compress.Compressor, cache cache.Cache) (*redisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    "/run/redis/redis.sock",
	})

	return &redisStorage{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (s *redisStorage) Get(nodeID int64) (node *Node, err error) {
	if err := s.client.Get(
		s.ctx,
		"count",
	).Err(); err != nil {
		if nodeID == 1 {
			return nil, ErrIsEmpty
		}
	}

	node = new(Node)
	err = s.client.HGetAll(
		s.ctx,
		strconv.Itoa(int(nodeID)),
	).Scan(node)

	return
}

func (s *redisStorage) SetObject(nodeID int64, position Position, marshaler encoding.BinaryMarshaler) error {
	return s.client.HSet(
		s.ctx,
		strconv.Itoa(int(nodeID)),
		position.Object(), marshaler,
	).Err()
}

func (s *redisStorage) SetChild(nodeID int64, position Position, child int64) error {
	return s.client.HSet(
		s.ctx,
		strconv.Itoa(int(nodeID)),
		position.Child(), child,
	).Err()
}

func (s *redisStorage) NewNode(marshaler encoding.BinaryMarshaler) (nodeID int64, err error) {
	count, err := s.client.Incr(s.ctx, "count").Uint64()
	if err != nil {
		return 0, err
	}
	err = s.client.HSet(
		s.ctx,
		strconv.Itoa(int(count)),
		Left.Object(), marshaler,
	).Err()

	return int64(count), err
}

func (s *redisStorage) Close() error {
	return s.client.Close()
}
