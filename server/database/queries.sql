-- name: CreateDriver :one
INSERT INTO drivers(username) VALUES ($1)
RETURNING *;

-- name: CreateRider :one
INSERT INTO riders(username) VALUES ($1)
RETURNING *;

-- name: CreateRide :one
INSERT INTO rides(rider_id, driver_id, pickup_lat, pickup_long, dropoff_lat, dropoff_long)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRideStatus :one
SELECT status FROM rides WHERE id = $1;