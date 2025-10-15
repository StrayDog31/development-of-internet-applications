package redis

import (
	"context"
	"fmt"
	"time"

	"development-of-internet-applications/internal/app/jwt"

	"github.com/go-redis/redis/v8"
)

const servicePrefix = "dia-backend:"
const jwtPrefix = "jwt:"

type Client struct {
	cfg    config.RedisConfig
	client *redis.Client
}

func New(cfg config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		Username:     cfg.User,
		DB:           0,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.ReadTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{
		cfg:    cfg,
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func getJWTKey(token string) string {
	return servicePrefix + jwtPrefix + token
}

func (c *Client) WriteJWTToBlacklist(ctx context.Context, jwtStr string, jwtTTL time.Duration) error {
	return c.client.Set(ctx, getJWTKey(jwtStr), true, jwtTTL).Err()
}

func (c *Client) CheckJWTInBlacklist(ctx context.Context, jwtStr string) error {
	return c.client.Get(ctx, getJWTKey(jwtStr)).Err()
}