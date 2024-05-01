package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(connetTo string) (*gorm.DB,error){
	db,err:=gorm.Open(postgres.Open(connetTo),&gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err!=nil{
		return nil,err
	}

	return db,nil
}