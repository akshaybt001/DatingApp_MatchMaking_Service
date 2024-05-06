package adapters

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MatchAdapter struct {
	DB *gorm.DB
}

func NewMatchAdapter(db *gorm.DB) *MatchAdapter {
	return &MatchAdapter{
		DB: db,
	}
}

func (m *MatchAdapter) IsLikeExist(userId, likedId string) (bool, error) {
	var count int
	if err := m.DB.Raw(`SELECT COUNT(*) FROM likes WHERE user_id=? AND liked_id=?`, userId, likedId).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *MatchAdapter) Like(userId, likedId string) error {
	fmt.Println(234567)

	selectquery := `INSERT INTO likes (user_id, liked_id, time) VALUES($1,$2,$3)`
	if err := m.DB.Exec(selectquery, userId, likedId, time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (m *MatchAdapter) Match(userId, matchId string) error {
	id := uuid.New()
	if err := m.DB.Exec(`INSERT INTO matches(id, user_id,match_id,time) VALUES(?,?,?,?)`, id, userId, matchId, time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (m *MatchAdapter) UnMatch(userId, matchId uint) error {
	if err := m.DB.Exec(`DELETE FROM matches WHERE user_id=? AND match_id=? `, userId, matchId).Error; err != nil {
		return err
	}
	return nil
}
