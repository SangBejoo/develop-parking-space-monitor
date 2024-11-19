BEGIN;

-- Drop tables if they exist to ensure clean state
DROP TABLE IF EXISTS monitoring_place;
DROP TABLE IF EXISTS map_hexagon_place;
DROP TABLE IF EXISTS trx_supply;

-- Create tables
CREATE TABLE map_hexagon_place (
    id SERIAL PRIMARY KEY,
    hexagon_id VARCHAR(255) NOT NULL,
    place_id INTEGER,
    place_type_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trx_supply (
    id SERIAL PRIMARY KEY,
    fleet_number VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    driver_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;