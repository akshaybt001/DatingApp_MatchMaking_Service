package service

import (
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/adapters"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
)

type MatchService struct {
	adapters adapters.AdapterInterface
	pb.UnimplementedMatchServiceServer
}

func NewMatchService(adapters adapters.AdapterInterface) *MatchService{
	return &MatchService{
		adapters: adapters,
	}
}