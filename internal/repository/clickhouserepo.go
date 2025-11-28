package repository

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Stats struct {
	Id           string
	Url_short_id string
	User_ip      string
	User_agent   string
	Referer      string
	Country      string
	City         string
	Device       string
	OS           string
	Browser      string
	Timestamp    time.Time
}

type ClkHouseRepo struct {
	ClkhouseConn clickhouse.Conn
}

func NewClkHouseRepo(db clickhouse.Conn) *ClkHouseRepo {
	return &ClkHouseRepo{ClkhouseConn: db}
}

func (r *ClkHouseRepo) Insert(ctx context.Context, stats Stats) error {
	batch, err := r.ClkhouseConn.PrepareBatch(ctx, "INSERT INTO stats (id, url_short_id, user_ip, user_agent, referer, country, city, device, os, browser, timestamp)")

	if err != nil {
		return err
	}
	defer batch.Abort()

	err = batch.Append(
		stats.Id,
		stats.Url_short_id,
		stats.User_ip,
		stats.User_agent,
		stats.Referer,
		stats.Country,
		stats.City,
		stats.Device,
		stats.OS,
		stats.Browser,
		stats.Timestamp,
	)

	if err != nil {
		return err
	}

	return batch.Send()
}

func (r *ClkHouseRepo) GetBySHortID(ctx context.Context, shortURLID string) ([]Stats, error) {
	rows, err := r.ClkhouseConn.Query(ctx, "SELECT id, url_short_id, user_ip, user_agent, referer, country, city, device, os, browser, timestamp FROM stats WHERE url_short_id=?", shortURLID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Stats

	for rows.Next() {
		var s Stats

		if err := rows.Scan(
			&s.Id,
			&s.Url_short_id,
			&s.User_ip,
			&s.User_agent,
			&s.Referer,
			&s.Country,
			&s.City,
			&s.Device,
			&s.OS,
			&s.Browser,
			&s.Timestamp,
		); err != nil {
			return nil, err
		}

		results = append(results, s)
	}

	return results, nil
}

func (r *ClkHouseRepo) GetStatsByDateRange(ctx context.Context, shortURLID string, start, end time.Time) ([]Stats, error) {
	rows, err := r.ClkhouseConn.Query(ctx, "SELECT id, url_short_id, user_ip, user_agent, referer, country, city, device, os, browser, timestamp WHERE short_url_id=? AND timestamp BETWEEN ? AND ?", shortURLID, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Stats

	for rows.Next() {
		var s Stats
		if err := rows.Scan(
			&s.Id,
			&s.Url_short_id,
			&s.User_ip,
			&s.User_agent,
			&s.Referer,
			&s.Country,
			&s.City,
			&s.Device,
			&s.OS,
			&s.Timestamp,
		); err != nil {
			return nil, err
		}

		results = append(results, s)
	}

	return results, nil
}

func (r *ClkHouseRepo) GetByContry(ctx context.Context, shortURLID, country string) ([]Stats, error) {
	rows, err := r.ClkhouseConn.Query(ctx, "SELECT id, url_short_id, user_ip, user_agent, referer, country, city, device, os, browser, timestamp FROM stats WHERE short_url_id=? AND country=?", shortURLID, country)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Stats

	for rows.Next() {
		var s Stats

		if err := rows.Scan(
			&s.Id,
			&s.Url_short_id,
			&s.User_ip,
			&s.User_agent,
			&s.Referer,
			&s.Country,
			&s.City,
			&s.Device,
			&s.OS,
			&s.Timestamp,
		); err != nil {
			return nil, err
		}

		results = append(results, s)
	}

	return results, nil
}
