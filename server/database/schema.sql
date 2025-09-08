-- riders
CREATE TABLES IF NOT EXISTS riders (
    id BIGSERIAL PRIMARY KEY,
    username TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() 
);

-- drivers
CREATE TABLE IF NOT EXIST drivers {
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    status SMALLINT DEFAULT 0, -- 0= available, 1=assigned, 2=enroute, 3=completed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
};

-- ride_requests
CREATE TABLE IF NOT EXISTS ride_requests (
    id BIGSERIAL PRIMARY KEY,
    rider_id BIGINT NOT NULL REFERENCES rider(id),
    driver_id BIGINT REFERENCES driver(id),
    status SMALLINT DEFAULT 0, -- 0=requested, 1=assigned, 2=accepted, 3=completed
    pickup_lat DOUBLE PRECISION NOT NULL,  
    pickup_long DOUBLE PRECISION NOT NULL,
    dropoff_lat DOUBLE PRECISION NOT NULL,
    dropoff_long DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
);

-- driver_location_path
CREATE TABLE IF NOT EXISTS driver_location_path (
    id BIGSERIAL PRIMARY KEY,
    driver_id BIGINT NOT NULL REFERENCES driver(id),
    lat DOUBLE PRECISION,
    long DOUBLE PRECISION,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
);