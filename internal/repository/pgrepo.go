package repository

import (
	"biturl/internal/domain"

	"gorm.io/gorm"
)

type PGrepo interface {
	CreateShortURL(input *domain.URL) error
	LoadURL(shortCode string) (domain.URL, error)
	DeleteURL(shortCode string) error
	ShortCodeExists(shortCode string) bool
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

// DeleteURL implements [PGrepo].
func (u *pgRepo) DeleteURL(shortCode string) error {
	url := domain.URL{}
	return u.DB.Where("short_code=?", shortCode).Delete(&url).Error
}

func (u *pgRepo) ShortCodeExists(shortCode string) bool {
	var url domain.URL
	err := u.DB.Where("short_code=?", shortCode).First(&url).Error
	if err != nil {
		return false
	}

	return true
}

func NewPostgresRepo(db *gorm.DB) PGrepo {
	return &pgRepo{DB: db}
}
