package adapters

import (
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/entities"
	helperstruct "github.com/akshaybt001/DatingApp_MatchMaking_Service/entities/helperStruct"
)

type AdapterInterface interface {
	IsLikeExist(userId, likedId string) (bool, error)
	Like(userId, likedId string) error
	Match(userId, matchId string) error
	UnMatch(id string) error
	GetMatch(userId string) ([]helperstruct.Match, error)
	IsMatchExist(id string) (bool, error)
	FindWhoLikesUser(id string)([]entities.Likes,error)
}
