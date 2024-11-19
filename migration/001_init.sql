
CREATE TABLE IF NOT EXISTS trx_supply (
    id SERIAL PRIMARY KEY,
    fleet_number VARCHAR(50) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    driver_id VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS map_hexagon_place (
    id SERIAL PRIMARY KEY,
    hexagon_id VARCHAR(50) NOT NULL,
    place_id INTEGER NOT NULL,
    place_type_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_hexagon_id ON map_hexagon_place(hexagon_id);
CREATE INDEX idx_fleet_number ON trx_supply(fleet_number);