package adapters

type AdapterInterface interface {
	IsLikeExist(userId, likedId string) (bool, error)
	Like(userId, likedId string) error
	Match(userId, matchId string) error
	UnMatch(userId, matchId uint) error
}
