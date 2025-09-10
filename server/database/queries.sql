-- name: CreateDriver :one
INSERT INTO drivers(username, password) 
VALUES ($1, $2)
RETURNING id, username, password;

-- name: GetDriverByUsername :one
SELECT id, username, password FROM drivers
WHERE username = $1
LIMIT 1;

-- name: UpdateDriverLocation :exec
UPDATE drivers
SET lat = $2, long = $3, available = true
WHERE id = $1;

-- name: CreateRider :one
INSERT INTO riders(username, password) 
VALUES ($1, $2)
RETURNING id, username, password;

-- name: GetRiderByUsername :one
SELECT id, username, password FROM riders
WHERE username = $1
LIMIT 1;

-- name: CreateRide :one
INSERT INTO rides(rider_id, driver_id, pickup_lat, pickup_long, dropoff_lat, dropoff_long, status, requested_at)
VALUES ($1, $2, $3, $4, $5, $6, 0, now()) -- 0=requested
RETURNING id, rider_id, driver_id, pickup_lat, pickup_long, dropoff_lat, dropoff_long, status, requested_at;

-- name: GetRideStatus :one
SELECT status FROM rides
WHERE id = $1;

-- name: GetRideByID :one
SELECT id, rider_id, driver_id, pickup_lat, pickup_long, dropoff_lat, dropoff_long, status, requested_at, accepted_at, completed_at
FROM rides
WHERE id = $1;

-- name: AssignDriverToRide :exec
UPDATE rides
SET driver_id = $2, status = 2, accepted_at = now()
WHERE id = $1 AND status = 0;

-- name: GetNearbyDrivers :many
SELECT d.id AS driver_id, d.username, dlp.lat, dlp.long
FROM driver_location_path AS dlp
JOIN drivers d ON d.id = dlp.driver_id
WHERE d.status = 0  -- available
AND earth_distance(
    ll_to_earth($1, $2), 
    ll_to_earth(dlp.lat, dlp.long)
) < $3;

-- name: GetAnalytics :one
SELECT 
    AVG(EXTRACT(EPOCH FROM (accepted_at - requested_at))/60) AS avg_wait,
    COUNT(*) FILTER (WHERE status='completed') AS completed_count
FROM rides;