package admin

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

func (s *Store) ListUsers() ([]*model.User, error) {
	var users []*model.User
	err := s.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
