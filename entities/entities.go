package entities

import (
	"time"

	"github.com/google/uuid"
)

type Likes struct {
	Id      int       `json:"id" gorm:"primaryKey"`
	UserId  string    `json:"user_id"`
	LikedId string    `json:"liked_id"`
	Time    time.Time `json:"time"`
}
type Match struct {
	Id      uuid.UUID `json:"id" `
	UserId  string    `json:"user_id"`
	MatchId string    `json:"match_id"`
	Time    time.Time `json:"time"`
}
