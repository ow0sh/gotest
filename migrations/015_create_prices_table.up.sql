CREATE TABLE IF NOT EXISTS prices(
    id serial PRIMARY KEY,
    base VARCHAR(128),
    quote VARCHAR(128),
    rate FLOAT
);