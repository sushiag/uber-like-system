-- riders
CREATE TABLE IF NOT EXISTS riders (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- drivers
CREATE TABLE IF NOT EXISTS drivers (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    status SMALLINT DEFAULT 0, -- 0= available, 1=assigned, 2=enroute, 3=completed
    lat DOUBLE PRECISION,      -- current latitude
    long DOUBLE PRECISION,     -- current longitude
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- rides
CREATE TABLE IF NOT EXISTS rides (
    id BIGSERIAL PRIMARY KEY,
    rider_id BIGINT NOT NULL REFERENCES riders(id),
    driver_id BIGINT REFERENCES drivers(id),
    status SMALLINT DEFAULT 0, -- 0=requested, 1=assigned, 2=accepted, 3=completed
    pickup_lat DOUBLE PRECISION NOT NULL,  
    pickup_long DOUBLE PRECISION NOT NULL,
    dropoff_lat DOUBLE PRECISION NOT NULL,
    dropoff_long DOUBLE PRECISION NOT NULL,
    requested_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    accepted_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- driver_location_path
CREATE TABLE IF NOT EXISTS driver_location_path (
    id BIGSERIAL PRIMARY KEY,
    driver_id BIGINT NOT NULL REFERENCES drivers(id),
    lat DOUBLE PRECISION,
    long DOUBLE PRECISION,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);