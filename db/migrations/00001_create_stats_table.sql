-- +goose Up
CREATE TABLE IF NOT EXISTS stats (
    id String,
    url_short_id String,
    user_ip String,
    user_agent String,
    referer String,
    country String,
    city String,
    device String,
    os String,
    browser String,
    timestamp DateTime
) 
ENGINE = MergeTree()
ORDER BY (url_short_id, timestamp);

-- +goose Down
DROP TABLE IF EXISTS stats;