package initializer

import (
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/adapters"
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/service"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB)*service.MatchService{
	repo:=adapters.NewMatchAdapter(db)
	service:=service.NewMatchService(repo)
	
	return service
}