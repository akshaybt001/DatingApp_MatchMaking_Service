package adapters

import (
	"time"

	"github.com/akshaybt001/DatingApp_MatchMaking_Service/entities"
	helperstruct "github.com/akshaybt001/DatingApp_MatchMaking_Service/entities/helperStruct"
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

func (m *MatchAdapter) UnMatch(id string) error {
	if err := m.DB.Exec(`DELETE FROM matches WHERE id=?`, id).Error; err != nil {
		return err
	}
	return nil
}

func (m *MatchAdapter) GetMatch(userId string) ([]helperstruct.Match, error) {
	var res []helperstruct.Match
	selectMatchQuery := `SELECT id , user_id , match_id FROM matches WHERE user_id=? OR match_id=?`
	if err := m.DB.Raw(selectMatchQuery, userId, userId).Scan(&res).Error; err != nil {
		return []helperstruct.Match{}, err
	}
	return res, nil
}

func (m *MatchAdapter) IsMatchExist(id string) (bool, error) {
	var count int
	if err := m.DB.Raw(`SELECT COUNT(*) FROM matches WHERE id=?`, id).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *MatchAdapter) FindWhoLikesUser(id string)([]entities.Likes,error){
	var res []entities.Likes
	selectUserlikes:=`SELECT * FROM likes WHERE liked_id=?`
	if err:=m.DB.Raw(selectUserlikes,id).Scan(&res).Error;err!=nil{
		return []entities.Likes{},nil
	}
	return res,nil
}