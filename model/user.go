package model

import (
	// "log"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            string `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Avatar        string
	Name          string
	Email         string
	Password      string
	Bio           string
	EmailVerified bool
	Followers     []Follow  `gorm:"foreignkey:FollowingID"`
	Followings    []Follow  `gorm:"foreignkey:FollowerID"`
}

type Follow struct {
	Follower    User
	FollowerID  string `gorm:"primary_key" sql:"type:int not null"`
	Following   User
	FollowingID string `gorm:"primary_key" sql:"type:int not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	return nil
}

type Auth struct {
	AccessToken           string    `json:"token" gorm:"-"`
	AuthState             AuthState `json:"authState" gorm:"-"`
	ExpiresIn             int64     `json:"expiresIn" gorm:"-"`
	RefreshToken          string    `json:"refreshToken" gorm:"-"`
	RefreshTokenExpiresIn int64     `json:"refreshTokenExpiresIn" gorm:"-"`
	AccessTokenUUID       string    `json:"-" gorm:"-"`
	RefreshTokenUUID      string    `json:"-" gorm:"-"`
}

type AuthState struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

type History struct {
	CreatedAt time.Time `json:"loginTime"`
	UserID    string    `json:"userId"`
	Email     string    `json:"email"`
}

type EmailVerification struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Code   string `json:"code"`
}

// type Profile struct {
// 	ID        string         `gorm:"primaryKey" json:"id"`
// 	CreatedAt time.Time      `json:"createdAt"`
// 	UpdatedAt time.Time      `json:"updatedAt"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
// 	Avatar    string         `json:"avatar"`
// 	Name      string         `json:"name"`
// 	Email     string         `json:"email"`
// 	Password  string
// 	Version   int
// }

// FollowedBy Followings should be pre loaded
func (u *User) FollowedBy(id string) bool {
	if u.Followers == nil {
		return false
	}
	for _, f := range u.Followers {
		if f.FollowerID == id {
			return true
		}
	}
	return false
}
