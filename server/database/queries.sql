-- name: CreateDriver :one
INSERT INTO drivers(username, password) 
VALUES ($1, $2)
RETURNING id, username, password;

-- name: GetDriverByUsername :one
SELECT id, username, password FROM drivers
WHERE username = $1
LIMIT 1;

-- name: CreateRider :one
INSERT INTO riders(username, password) 
VALUES ($1, $2)
RETURNING id, username, password;

-- name: GetRiderByUsername :one
SELECT id, username, password FROM riders
WHERE username = $1
LIMIT 1;

-- name: CreateRide :one
INSERT INTO rides(rider_id, driver_id, pickup_lat, pickup_long, dropoff_lat, dropoff_long)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRideStatus :one
SELECT status FROM rides WHERE id = $1;