package adapters

import "gorm.io/gorm"

type MatchAdapter struct {
	DB *gorm.DB
}

func NewMatchAdapter(db *gorm.DB) *MatchAdapter{
	return &MatchAdapter{
		DB: db,
	}
}