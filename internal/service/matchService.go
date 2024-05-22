package service

import (
	"context"
	"fmt"
	"log"

	"github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/adapters"
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/helper"
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/kafka"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
)

type MatchService struct {
	UserClient         pb.UserServiceClient
	NotificationClient pb.NotificationClient
	adapters           adapters.AdapterInterface
	pb.UnimplementedMatchServiceServer
}

func NewMatchService(adapters adapters.AdapterInterface, useraddr, notificationaddr string) *MatchService {
	userConn, _ := helper.DialGrpc(useraddr)
	notifyConn, _ := helper.DialGrpc(notificationaddr)
	return &MatchService{
		adapters:           adapters,
		UserClient:         pb.NewUserServiceClient(userConn),
		NotificationClient: pb.NewNotificationClient(notifyConn),
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

	LikedData, err := m.UserClient.GetUserData(context.Background(), &pb.GetUserById{Id: req.LikedId})
	if err != nil {
		return &pb.NoPara{}, err
	}
	if !UserData.IsSubscribed {
		likes := UserData.LikeCount
		if likes == 0 {
			return &pb.NoPara{}, fmt.Errorf("your like limit exceeded")
		}

		err := m.adapters.Like(req.UserId, req.LikedId)
		if err != nil {
			return nil, err
		}

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
	_, err = m.NotificationClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
		UserId:  req.LikedId,
		Message: fmt.Sprintf(`{"message": "hey there %s likes you"}`, UserData.Name),
	})
	if err != nil {
		log.Print("error while sending notification ", err)
	}
	fmt.Println("notification sent.........")
	exist, err := m.adapters.IsLikeExist(req.LikedId, req.UserId)
	if err != nil {
		return nil, err
	}
	if exist {
		err := m.adapters.Match(req.UserId, req.LikedId)
		if err != nil {
			return nil, err
		}
		_, err = m.NotificationClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
			UserId:  req.LikedId,
			Message: fmt.Sprintf(`{"message": "hey there %s matches you"}`, UserData.Name),
		})
		if err != nil {
			log.Print("error while sending notification ", err)
		}
		_, err = m.NotificationClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
			UserId:  req.UserId,
			Message: fmt.Sprintf(`{"message": "hey there %s matches you"}`, LikedData.Name),
		})
		if err != nil {
			log.Print("error while sending notification ", err)
		}
		fmt.Println("notification sent.........")
		err = kafka.ProduceMatchUserMessage(LikedData.Name, UserData.Email)
		if err != nil {
			fmt.Println(err)
		}
		err = kafka.ProduceMatchUserMessage(UserData.Name, LikedData.Email)
		if err != nil {
			fmt.Println(err)
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

func (m *MatchService) GetWhoLikesUser(req *pb.GetByUserId, srv pb.MatchService_GetWhoLikesUserServer) error {
	userlist, err := m.adapters.FindWhoLikesUser(req.Id)
	if err != nil {
		return err
	}
	UserData, err := m.UserClient.GetUserData(context.Background(), &pb.GetUserById{Id: req.Id})
	if err != nil {
		return err
	}
	if !UserData.IsSubscribed {
		return fmt.Errorf("you are not subscribed \n kindly please do Subscription")
	}
	for _, user := range userlist {
		UserData, err := m.UserClient.GetUserData(context.Background(), &pb.GetUserById{Id: user.UserId})
		if err != nil {
			return err
		}
		res := &pb.LikedUsersResposne{
			Id:        UserData.Id,
			LikedName: UserData.Name,
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
