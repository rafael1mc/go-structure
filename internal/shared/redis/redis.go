package redis

import (
	"context"
	"fmt"
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/util/consts"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
	redis "github.com/redis/go-redis/v9"
)

type RedisWrapper struct {
	client *redis.Client
	ctx    context.Context
	logger *slog.Logger
}

func NewRedisWrapper(
	env *env.Env,
	logger *slog.Logger,
) *RedisWrapper {
	ctx := context.Background()
	var client *redis.Client

	addr := fmt.Sprintf("%s:%d", env.Redis.Host, env.Redis.Port)
	logger.Debug("Initializing redis.", slog.String("conn_string", addr))

	conn := func() error {
		client = redis.NewClient(
			&redis.Options{
				Addr:     addr,
				Password: env.Redis.Pass,
			},
		)

		if err := client.Ping(ctx).Err(); err != nil {
			logger.Error("could not ping redis", consts.SlogError(err))
			return err
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 2 * time.Second
	expBackoff.MaxElapsedTime = 30 * time.Second
	expBackoff.Reset()

	err := backoff.Retry(conn, expBackoff)
	if err != nil {
		panic(fmt.Errorf("Failed to connect to redis after retrying: %w", err))
	}

	return &RedisWrapper{
		client: client,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *RedisWrapper) Publish(key, payload string) error {
	return r.client.Publish(r.ctx, key, payload).Err()
}

func (r *RedisWrapper) Subscribe(
	key string,
	callback func(message *redis.Message),
) {
	subscriber := r.client.Subscribe(r.ctx, key)
	defer subscriber.Close()

	for message := range subscriber.Channel() {
		callback(message)
	}
}

func (r *RedisWrapper) Save(
	key string,
	value string,
	expiration time.Duration,
) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisWrapper) Increment(
	key string,
) int64 {
	return r.client.Incr(r.ctx, key).Val()
}

func (r *RedisWrapper) Expire(
	key string,
	d time.Duration,
) {
	r.client.Expire(r.ctx, key, d)
}

func (r *RedisWrapper) Del(
	key ...string,
) {
	r.client.Del(r.ctx, key...)
}

func (r *RedisWrapper) Get(
	key string,
) (string, error) {
	result, err := r.client.Get(r.ctx, key).Result()

	return string(result), err
}

func (r *RedisWrapper) GetWithPattern(
	pattern string,
) (map[string]string, error) {
	values := make(map[string]string)
	err := r.GetWithPatternIterable(pattern, func(key, value string) bool {
		values[key] = value
		return false
	})

	return values, err
}

func (r *RedisWrapper) GetWithPatternIterable(
	pattern string,
	iterable func(key, value string) bool,
) error {
	// Scan for matching keys
	iter := r.client.Scan(r.ctx, 0, pattern, 0).Iterator()

	for iter.Next(r.ctx) {
		key := iter.Val()
		value, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			return err
		}
		if shouldStop := iterable(key, value); shouldStop {
			return err
		}
	}
	return nil
}
