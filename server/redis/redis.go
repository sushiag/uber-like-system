package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	C *redis.Client
}

// New initializes a Redis client
func New(addr, pass string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0, // default DB
	})

	return &Client{C: rdb}
}

// SetDriverLocation stores driver's current location in Redis
func (r *Client) SetDriverLocation(ctx context.Context, driverID uint64, lat, long float64) error {
	key := fmt.Sprintf("driver:%d", driverID)
	return r.C.HSet(ctx, key, map[string]any{
		"lat":  lat,
		"long": long,
	}).Err()
}

// GetDriverLocation retrieves driver's location from Redis
func (r *Client) GetDriverLocation(ctx context.Context, driverID uint64) (float64, float64, error) {
	key := fmt.Sprintf("driver:%d", driverID)
	data, err := r.C.HGetAll(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}

	lat, okLat := data["lat"]
	long, okLong := data["long"]
	if !okLat || !okLong {
		return 0, 0, fmt.Errorf("driver location not found")
	}

	var latF, longF float64
	fmt.Sscanf(lat, "%f", &latF)
	fmt.Sscanf(long, "%f", &longF)

	return latF, longF, nil
}
