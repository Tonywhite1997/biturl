package repository

import (
	"biturl/internal/dto"
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

func (r *ClkHouseRepo) GetTotalClicks(ctx context.Context, shortURLID string) (uint64, error) {
	var total uint64
	err := r.ClkhouseConn.QueryRow(ctx, "SELECT COUNT(*) FROM stats WHERE url_short_id=?", shortURLID).Scan(&total)

	return total, err
}

func (r *ClkHouseRepo) GetUniqueVisitors(ctx context.Context, shortURLID string) (uint64, error) {
	var total uint64
	err := r.ClkhouseConn.QueryRow(ctx, "SELECT COUNT(DISTINCT user_ip) FROM stats WHERE url_short_id=?", shortURLID).Scan(&total)
	return total, err
}

func (r *ClkHouseRepo) GetBrowserClicks(ctx context.Context, shortURLID string) ([]dto.BrowserStats, error) {

	rows, err := r.ClkhouseConn.Query(ctx,
		`SELECT browser, COUNT(*) 
	FROM stats 
	WHERE url_short_id=? 
	GROUP BY browser 
	ORDER BY COUNT(*)`, shortURLID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalBrowsers []dto.BrowserStats

	for rows.Next() {
		var c dto.BrowserStats
		if err := rows.Scan(&c.Browser, &c.Count); err != nil {
			return nil, err
		}
		totalBrowsers = append(totalBrowsers, c)
	}

	return totalBrowsers, nil
}

func (r *ClkHouseRepo) GetCountryClicks(ctx context.Context, shortURLID string) ([]dto.CountryStats, error) {

	rows, err := r.ClkhouseConn.Query(ctx,
		`SELECT COUNTRY, COUNT(*) 
	FROM stats 
	WHERE url_short_id=? 
	GROUP BY country 
	ORDER BY COUNT(*) 
	DESC`, shortURLID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalCountries []dto.CountryStats
	for rows.Next() {
		var c dto.CountryStats

		if err := rows.Scan(&c.Country, &c.Count); err != nil {
			return nil, err
		}
		totalCountries = append(totalCountries, c)
	}

	return totalCountries, err
}

func (r *ClkHouseRepo) GetDailyClicks(ctx context.Context, shortURLID string) ([]dto.DailyStats, error) {
	rows, err := r.ClkhouseConn.Query(ctx,
		`SELECT 
		toDate(timestamp) AS day,
		COUNT(*) AS clicks
	FROM stats
	WHERE
		url_short_id=?
		AND timestamp >= now() - INTERVAL 30 DAY
	GROUP BY day
	ORDER BY day ASC
	 `, shortURLID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dailyClicks []dto.DailyStats
	for rows.Next() {
		var d dto.DailyStats
		if err := rows.Scan(&d.Date, &d.Count); err != nil {
			return nil, err
		}
		dailyClicks = append(dailyClicks, d)
	}
	return dailyClicks, nil
}

func (r *ClkHouseRepo) GetDeviceClicks(ctx context.Context, shortURLID string) ([]dto.DeviceStats, error) {
	rows, err := r.ClkhouseConn.Query(ctx,
		`SELECT device, COUNT(*)
	FROM stats
	WHERE url_short_id=?
	GROUP BY device
	ORDER BY COUNT(*)
	`, shortURLID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalDevices []dto.DeviceStats
	for rows.Next() {
		var d dto.DeviceStats
		if err := rows.Scan(&d.Device, &d.Count); err != nil {
			return nil, err
		}
		totalDevices = append(totalDevices, d)
	}
	return totalDevices, nil
}

func (r *ClkHouseRepo) DeleteStatsRecord(ctx context.Context, shortURLID string) error {
	err := r.ClkhouseConn.Exec(ctx, "DELETE FROM stats WHERE url_short_id=? ", shortURLID)
	if err != nil {
		return err
	}
	return nil
}
