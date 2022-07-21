package user

import (
	"github.com/mongmx/echo-template/model"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByID(id string) (*model.User, error) {
	user := model.User{ID: id}
	if err := s.db.Where(&user).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByEmail(email string) (*model.User, error) {
	user := model.User{Email: email}
	if err := s.db.Where(&user).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *Store) UpdateUser(user *model.User) error {
	return s.db.Save(user).Error
}

func (s *Store) DeleteUser(user *model.User) error {
	return s.db.Delete(user).Error
}

func (s * Store) GetHistory(id string) ([]model.History, error) {
	var histories []model.History
	if err := s.db.Where("user_id = ?", id).Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

func (s *Store) CreateHistory(history *model.History) error {
	return s.db.Create(history).Error
}

func (s *Store) GetProfile(username string) (*model.User, error) {
	u := model.User{Name: username}
	if err := s.db.Where(&u).Preload("Followers").First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) AddFollower(u *model.User, followerID string) error {
	f := model.Follow{FollowerID: followerID, FollowingID: u.ID}
	return s.db.Model(&u).Association("Followers").Append(&f)
}

func (s *Store) RemoveFollower(u *model.User, followerID string) error {
	f := model.Follow{FollowerID: followerID, FollowingID: u.ID}
	if err := s.db.Model(u).Association("Followers").Find(&f); err != nil {
		return err
	}
	return s.db.Delete(&f).Error
}

func (s *Store) IsFollower(userID, followerID string) (bool, error) {
	f := model.Follow{FollowerID: followerID, FollowingID: userID}
	if err := s.db.Where(&f).First(&f).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
