package repository

import (
	"biturl/internal/domain"
	"time"

	"gorm.io/gorm"
)

type PGrepo interface {
	CreateShortURL(input *domain.URL) error
	LoadURL(shortCode string) (domain.URL, error)
	LoadURLByAccessKey(accessKey string) (domain.URL, error)
	DeleteURL(shortCode string) error
	ShortCodeExists(statsAccessKey string) (bool, *string, *string)
	IncreaseExpiryDate(statsAccessKey string, newExpiry time.Time) error
	FindExpiredURLs() ([]domain.URL, error)
}

type pgRepo struct {
	DB *gorm.DB
}

// CreateURL implements [URLrepo].
func (u *pgRepo) CreateShortURL(url *domain.URL) error {
	return u.DB.Create(&url).Error
}

// LoadURL implements [URLrepo].
func (u *pgRepo) LoadURL(shortCode string) (domain.URL, error) {
	var url domain.URL
	err := u.DB.Where("short_code=?", shortCode).First(&url).Error
	return url, err
}

func (u *pgRepo) LoadURLByAccessKey(accessKey string) (domain.URL, error) {
	var url domain.URL
	err := u.DB.Where("stats_access_key", accessKey).First(&url).Error

	return url, err
}

// DeleteURL implements [PGrepo].
func (u *pgRepo) DeleteURL(shortCode string) error {
	url := domain.URL{}
	return u.DB.Where("short_code=?", shortCode).Delete(&url).Error
}

func (u *pgRepo) ShortCodeExists(statsAccessKey string) (bool, *string, *string) {
	var url domain.URL
	err := u.DB.Where("stats_access_key=?", statsAccessKey).First(&url).Error
	if err != nil {
		return false, nil, nil
	}

	return true, &url.ShortCode, &url.OriginalURL
}

func (u *pgRepo) IncreaseExpiryDate(accessKey string, newExpiry time.Time) error {
	return u.DB.Model(&domain.URL{}).
		Where("stats_access_key=?", accessKey).
		Update("expires_at", newExpiry).
		Error
}

func (u *pgRepo) FindExpiredURLs() ([]domain.URL, error) {
	var expiredURLs []domain.URL
	err := u.DB.Where("expired_at < NOW()").Find(&expiredURLs).Error
	return expiredURLs, err
}

func NewPostgresRepo(db *gorm.DB) PGrepo {
	return &pgRepo{DB: db}
}
