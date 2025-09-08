package redis

import (
	"context"
	"fmt"
	"uber-like-system/redis"
)

type Client struct {
	C *redis.Client
}

func new(addr, pass string) *Client {
	c := redis.NewClient(&redis.Options{
		Addr: addr,
		Pass: pass,
		DB:   0,
	})
	return &Client{C: c}
}

func (r *Client) SetDriverLocation(ctx context.Context, driveID uint64, lat, long float64) error {
	key := fmt.Sprintf(driverID)
	return r.C.HSet(ctx, key, map[string]interface{}{"lat": lat, "long": long}).Err()
}
