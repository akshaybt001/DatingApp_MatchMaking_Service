package db

import (
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(connetTo string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connetTo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&entities.Likes{})
	db.AutoMigrate(&entities.Match{})
	return db, nil
}
