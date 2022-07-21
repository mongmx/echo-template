package admin

import (
	"time"

	"github.com/mongmx/echo-template/model"
)

type profileResponse struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	Image     string    `json:"avatar"`
}

func marshalProfileResponse(u *model.User) *profileResponse {
	return &profileResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		Name:      u.Name,
		Email:     u.Email,
		Bio:       u.Bio,
		Image:     u.Avatar,
	}
}

// type singleProfileResponse struct {
// 	Profile *profileResponse `json:"profile"`
// }

// func newSingleProfileResponse(u *model.User) *singleProfileResponse {
// 	return &singleProfileResponse{
// 		Profile: marshalProfileResponse(u),
// 	}
// }

type listProfileResponse struct {
	Profiles []*profileResponse `json:"profiles"`
}

func newListProfileResponse(users []*model.User) *listProfileResponse {
	profiles := make([]*profileResponse, len(users))
	for i, u := range users {
		profiles[i] = marshalProfileResponse(u)
	}
	return &listProfileResponse{
		Profiles: profiles,
	}
}
