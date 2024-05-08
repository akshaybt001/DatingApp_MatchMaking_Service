package service

import (
	"context"
	"fmt"

	"github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/adapters"
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/helper"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
)

type MatchService struct {
	UserClient pb.UserServiceClient
	adapters   adapters.AdapterInterface
	pb.UnimplementedMatchServiceServer
}

func NewMatchService(adapters adapters.AdapterInterface, useraddr string) *MatchService {
	userConn, _ := helper.DialGrpc(useraddr)
	return &MatchService{
		adapters:   adapters,
		UserClient: pb.NewUserServiceClient(userConn),
	}
}

func (m *MatchService) Like(ctx context.Context, req *pb.LikeRequest) (*pb.NoPara, error) {
	if req.LikedId == req.UserId {
		return &pb.NoPara{}, fmt.Errorf("cannot like the user itself")
	}
	isExist, err := m.adapters.IsLikeExist(req.UserId, req.LikedId)
	if err != nil {
		return &pb.NoPara{}, err
	}
	if isExist {
		return &pb.NoPara{}, fmt.Errorf("liked already")
	}
	UserData, err := m.UserClient.GetUserData(context.Background(), &pb.GetUserById{Id: req.UserId})
	if err != nil {
		return &pb.NoPara{}, err
	}
	if !UserData.IsSubscribed {
		likes := UserData.LikeCount
		if likes == 0 {
			return &pb.NoPara{}, fmt.Errorf("your like limit exceeded")
		}
		fmt.Println(23456)
		err := m.adapters.Like(req.UserId, req.LikedId)
		if err != nil {
			return nil, err
		}
		fmt.Println(23456)
		_, err = m.UserClient.DecrementLikeCount(context.Background(), &pb.GetUserById{Id: req.UserId})
		if err != nil {
			return nil, err
		}

	} else {
		err := m.adapters.Like(req.UserId, req.LikedId)
		if err != nil {
			return nil, err
		}
	}
	exist, err := m.adapters.IsLikeExist(req.LikedId, req.UserId)
	if err != nil {
		return nil, err
	}
	if exist {
		err := m.adapters.Match(req.UserId, req.LikedId)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (m *MatchService) GetMatch(req *pb.GetByUserId, srv pb.MatchService_GetMatchServer) error {

	matches, err := m.adapters.GetMatch(req.Id)
	if err != nil {
		return err
	}
	for _, match := range matches {
		res := &pb.MatchResposne{
			Id:      match.Id,
			MatchId: match.MatchId,
			UserId:  match.UserId,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (m *MatchService) UnMatch(ctx context.Context, req *pb.GetByUserId) (*pb.NoPara, error) {
	match, err := m.adapters.IsMatchExist(req.Id)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf("cannot unmatch as it is not matched user")
	}
	err = m.adapters.UnMatch(req.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
