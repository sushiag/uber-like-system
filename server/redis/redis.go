package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	C *redis.Client
}

// New creates a new Redis client
func New(addr, pass string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})
	return &Client{C: rdb}
}

func (r *Client) SetDriverLocation(ctx context.Context, driverID uint64, lat, long float64) error {
	key := "driver_locations"
	return r.C.GeoAdd(ctx, key, &redis.GeoLocation{
		Name:      fmt.Sprintf("%d", driverID),
		Latitude:  lat,
		Longitude: long,
	}).Err()
}

func (r *Client) GetDriverLocation(ctx context.Context, driverID uint64) (float64, float64, error) {
	key := "driver_locations"
	pos, err := r.C.GeoPos(ctx, key, fmt.Sprintf("%d", driverID)).Result()
	if err != nil {
		return 0, 0, err
	}
	if len(pos) == 0 || pos[0] == nil {
		return 0, 0, fmt.Errorf("driver location not found")
	}
	return pos[0].Latitude, pos[0].Longitude, nil
}

func (r *Client) GetNearbyDrivers(ctx context.Context, lat, long float64, radius float64) ([]string, error) {
	key := "driver_locations"

	results, err := r.C.GeoSearch(ctx, key, &redis.GeoSearchQuery{
		Latitude:   lat,
		Longitude:  long,
		Radius:     radius,
		RadiusUnit: "m",
		Count:      5,
		Sort:       "ASC",
	}).Result()
	if err != nil {
		return nil, err
	}

	driverIDs := make([]string, len(results))
	copy(driverIDs, results)
	return driverIDs, nil
}
